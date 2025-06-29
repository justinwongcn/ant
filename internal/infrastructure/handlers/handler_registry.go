// Package handlers 提供处理器注册表的基础设施实现
package handlers

import (
	"fmt"
	"sync"

	"github.com/justinwongcn/ant/internal/domain/webserver/entities"
)

// HandlerRegistry 处理器注册表实现
type HandlerRegistry struct {
	routeHandlers      map[string]entities.HandlerFunc
	middlewareHandlers map[string]entities.MiddlewareFunc
	mu                 sync.RWMutex
}

// NewHandlerRegistry 创建新的处理器注册表
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		routeHandlers:      make(map[string]entities.HandlerFunc),
		middlewareHandlers: make(map[string]entities.MiddlewareFunc),
	}
}

// RegisterRouteHandler 注册路由处理器
func (r *HandlerRegistry) RegisterRouteHandler(name string, handler entities.HandlerFunc) error {
	if name == "" {
		return fmt.Errorf("处理器名称不能为空")
	}
	if handler == nil {
		return fmt.Errorf("处理器不能为空")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.routeHandlers[name]; exists {
		return fmt.Errorf("路由处理器已存在: %s", name)
	}

	r.routeHandlers[name] = handler
	return nil
}

// RegisterMiddlewareHandler 注册中间件处理器
func (r *HandlerRegistry) RegisterMiddlewareHandler(name string, handler entities.MiddlewareFunc) error {
	if name == "" {
		return fmt.Errorf("处理器名称不能为空")
	}
	if handler == nil {
		return fmt.Errorf("处理器不能为空")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.middlewareHandlers[name]; exists {
		return fmt.Errorf("中间件处理器已存在: %s", name)
	}

	r.middlewareHandlers[name] = handler
	return nil
}

// GetRouteHandler 获取路由处理器
func (r *HandlerRegistry) GetRouteHandler(name string) (entities.HandlerFunc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handler, exists := r.routeHandlers[name]
	if !exists {
		return nil, fmt.Errorf("路由处理器不存在: %s", name)
	}

	return handler, nil
}

// GetMiddlewareHandler 获取中间件处理器
func (r *HandlerRegistry) GetMiddlewareHandler(name string) (entities.MiddlewareFunc, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handler, exists := r.middlewareHandlers[name]
	if !exists {
		return nil, fmt.Errorf("中间件处理器不存在: %s", name)
	}

	return handler, nil
}

// ListRouteHandlers 列出所有路由处理器名称
func (r *HandlerRegistry) ListRouteHandlers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.routeHandlers))
	for name := range r.routeHandlers {
		names = append(names, name)
	}

	return names
}

// ListMiddlewareHandlers 列出所有中间件处理器名称
func (r *HandlerRegistry) ListMiddlewareHandlers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.middlewareHandlers))
	for name := range r.middlewareHandlers {
		names = append(names, name)
	}

	return names
}

// UnregisterRouteHandler 注销路由处理器
func (r *HandlerRegistry) UnregisterRouteHandler(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.routeHandlers[name]; !exists {
		return fmt.Errorf("路由处理器不存在: %s", name)
	}

	delete(r.routeHandlers, name)
	return nil
}

// UnregisterMiddlewareHandler 注销中间件处理器
func (r *HandlerRegistry) UnregisterMiddlewareHandler(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.middlewareHandlers[name]; !exists {
		return fmt.Errorf("中间件处理器不存在: %s", name)
	}

	delete(r.middlewareHandlers, name)
	return nil
}

// Clear 清空所有处理器（用于测试）
func (r *HandlerRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.routeHandlers = make(map[string]entities.HandlerFunc)
	r.middlewareHandlers = make(map[string]entities.MiddlewareFunc)
}

// RouteHandlerCount 返回路由处理器数量
func (r *HandlerRegistry) RouteHandlerCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.routeHandlers)
}

// MiddlewareHandlerCount 返回中间件处理器数量
func (r *HandlerRegistry) MiddlewareHandlerCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.middlewareHandlers)
}

// 示例处理器实现

// SimpleRouteHandler 简单路由处理器
type SimpleRouteHandler struct {
	name     string
	response string
}

// NewSimpleRouteHandler 创建简单路由处理器
func NewSimpleRouteHandler(name, response string) *SimpleRouteHandler {
	return &SimpleRouteHandler{
		name:     name,
		response: response,
	}
}

// Handle 处理请求
func (h *SimpleRouteHandler) Handle(ctx entities.RequestContext) error {
	// 简单实现：设置响应内容
	ctx.SetBody([]byte(h.response))
	return nil
}

// Name 返回处理器名称
func (h *SimpleRouteHandler) Name() string {
	return h.name
}

// SimpleMiddlewareHandler 简单中间件处理器
type SimpleMiddlewareHandler struct {
	name string
}

// NewSimpleMiddlewareHandler 创建简单中间件处理器
func NewSimpleMiddlewareHandler(name string) *SimpleMiddlewareHandler {
	return &SimpleMiddlewareHandler{
		name: name,
	}
}

// Process 处理中间件逻辑
func (h *SimpleMiddlewareHandler) Process(next entities.HandlerFunc) entities.HandlerFunc {
	return &WrappedHandler{
		name:       fmt.Sprintf("%s-wrapped", h.name),
		next:       next,
		middleware: h.name,
	}
}

// Name 返回中间件名称
func (h *SimpleMiddlewareHandler) Name() string {
	return h.name
}

// WrappedHandler 包装的处理器
type WrappedHandler struct {
	name       string
	next       entities.HandlerFunc
	middleware string
}

// Handle 处理请求
func (h *WrappedHandler) Handle(ctx entities.RequestContext) error {
	// 中间件前置处理
	// 这里可以添加日志、认证、限流等逻辑

	// 调用下一个处理器
	if h.next != nil {
		if err := h.next.Handle(ctx); err != nil {
			return err
		}
	}

	// 中间件后置处理
	// 这里可以添加响应修改、清理等逻辑

	return nil
}

// Name 返回处理器名称
func (h *WrappedHandler) Name() string {
	return h.name
}

// 确保实现了接口
var (
	_ entities.HandlerFunc    = (*SimpleRouteHandler)(nil)
	_ entities.MiddlewareFunc = (*SimpleMiddlewareHandler)(nil)
	_ entities.HandlerFunc    = (*WrappedHandler)(nil)
)
