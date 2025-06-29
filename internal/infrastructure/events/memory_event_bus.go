// Package events 提供事件处理的基础设施实现
package events

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/justinwongcn/ant/internal/domain/shared/events"
)

// MemoryEventBus 基于内存的事件总线实现
type MemoryEventBus struct {
	handlers map[string][]events.EventHandler
	mu       sync.RWMutex
}

// NewMemoryEventBus 创建新的内存事件总线
func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		handlers: make(map[string][]events.EventHandler),
	}
}

// Publish 发布事件
func (bus *MemoryEventBus) Publish(ctx context.Context, domainEvents ...events.DomainEvent) error {
	for _, event := range domainEvents {
		if err := bus.publishSingle(ctx, event); err != nil {
			log.Printf("发布事件失败: %v", err)
			// 继续处理其他事件，不中断整个流程
		}
	}
	return nil
}

// publishSingle 发布单个事件
func (bus *MemoryEventBus) publishSingle(ctx context.Context, event events.DomainEvent) error {
	bus.mu.RLock()
	handlers, exists := bus.handlers[event.EventType()]
	bus.mu.RUnlock()

	if !exists || len(handlers) == 0 {
		// 没有处理器，记录日志但不报错
		log.Printf("事件 %s 没有注册的处理器", event.EventType())
		return nil
	}

	// 异步处理事件，避免阻塞主流程
	for _, handler := range handlers {
		go func(h events.EventHandler, e events.DomainEvent) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("事件处理器发生panic: %v", r)
				}
			}()

			if err := h.Handle(ctx, e); err != nil {
				log.Printf("事件处理失败: %v", err)
			}
		}(handler, event)
	}

	return nil
}

// Subscribe 订阅事件
func (bus *MemoryEventBus) Subscribe(handler events.EventHandler) error {
	if handler == nil {
		return fmt.Errorf("事件处理器不能为空")
	}

	eventTypes := handler.EventTypes()
	if len(eventTypes) == 0 {
		return fmt.Errorf("事件处理器必须指定处理的事件类型")
	}

	bus.mu.Lock()
	defer bus.mu.Unlock()

	for _, eventType := range eventTypes {
		bus.handlers[eventType] = append(bus.handlers[eventType], handler)
	}

	return nil
}

// Unsubscribe 取消订阅事件
func (bus *MemoryEventBus) Unsubscribe(handler events.EventHandler) error {
	if handler == nil {
		return fmt.Errorf("事件处理器不能为空")
	}

	eventTypes := handler.EventTypes()
	if len(eventTypes) == 0 {
		return nil // 没有事件类型，无需取消订阅
	}

	bus.mu.Lock()
	defer bus.mu.Unlock()

	for _, eventType := range eventTypes {
		handlers, exists := bus.handlers[eventType]
		if !exists {
			continue
		}

		// 移除指定的处理器
		var newHandlers []events.EventHandler
		for _, h := range handlers {
			if h != handler {
				newHandlers = append(newHandlers, h)
			}
		}

		if len(newHandlers) == 0 {
			delete(bus.handlers, eventType)
		} else {
			bus.handlers[eventType] = newHandlers
		}
	}

	return nil
}

// GetHandlerCount 获取指定事件类型的处理器数量（用于测试）
func (bus *MemoryEventBus) GetHandlerCount(eventType string) int {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	handlers, exists := bus.handlers[eventType]
	if !exists {
		return 0
	}

	return len(handlers)
}

// Clear 清空所有处理器（用于测试）
func (bus *MemoryEventBus) Clear() {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	bus.handlers = make(map[string][]events.EventHandler)
}

// 确保实现了接口
var (
	_ events.EventBus       = (*MemoryEventBus)(nil)
	_ events.EventPublisher = (*MemoryEventBus)(nil)
)

// LoggingEventHandler 日志事件处理器示例
type LoggingEventHandler struct {
	eventTypes []string
}

// NewLoggingEventHandler 创建新的日志事件处理器
func NewLoggingEventHandler(eventTypes ...string) *LoggingEventHandler {
	return &LoggingEventHandler{
		eventTypes: eventTypes,
	}
}

// Handle 处理事件
func (h *LoggingEventHandler) Handle(ctx context.Context, event events.DomainEvent) error {
	log.Printf("处理事件: 类型=%s, 聚合ID=%s, 时间=%s",
		event.EventType(),
		event.AggregateID(),
		event.OccurredAt().Format("2006-01-02 15:04:05"))
	return nil
}

// EventTypes 返回处理的事件类型
func (h *LoggingEventHandler) EventTypes() []string {
	return h.eventTypes
}

// 确保实现了接口
var _ events.EventHandler = (*LoggingEventHandler)(nil)
