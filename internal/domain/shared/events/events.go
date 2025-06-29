// Package events 定义Ant Web Framework的领域事件
package events

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// DomainEvent 表示一个领域事件
type DomainEvent interface {
	EventID() string
	EventType() string
	AggregateID() string
	AggregateType() string
	OccurredAt() time.Time
	Version() int
	Data() interface{}
}

// BaseDomainEvent 为领域事件提供基础实现
type BaseDomainEvent struct {
	eventID       string
	eventType     string
	aggregateID   string
	aggregateType string
	occurredAt    time.Time
	version       int
	data          interface{}
}

// NewBaseDomainEvent 创建一个新的基础领域事件
func NewBaseDomainEvent(eventType, aggregateID, aggregateType string, version int, data interface{}) *BaseDomainEvent {
	return &BaseDomainEvent{
		eventID:       uuid.New().String(),
		eventType:     eventType,
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
		occurredAt:    time.Now(),
		version:       version,
		data:          data,
	}
}

func (e *BaseDomainEvent) EventID() string       { return e.eventID }
func (e *BaseDomainEvent) EventType() string     { return e.eventType }
func (e *BaseDomainEvent) AggregateID() string   { return e.aggregateID }
func (e *BaseDomainEvent) AggregateType() string { return e.aggregateType }
func (e *BaseDomainEvent) OccurredAt() time.Time { return e.occurredAt }
func (e *BaseDomainEvent) Version() int          { return e.version }
func (e *BaseDomainEvent) Data() interface{}     { return e.data }

// EventPublisher 发布领域事件
type EventPublisher interface {
	Publish(ctx context.Context, events ...DomainEvent) error
}

// EventHandler 处理领域事件
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
	EventTypes() []string
}

// EventBus 管理事件发布和处理
type EventBus interface {
	EventPublisher
	Subscribe(handler EventHandler) error
	Unsubscribe(handler EventHandler) error
}

// =============================================================================
// Web Server Domain Events
// =============================================================================

const (
	ServerStartedEventType    = "server.started"
	ServerStoppedEventType    = "server.stopped"
	RouteRegisteredEventType  = "route.registered"
	MiddlewareAddedEventType  = "middleware.added"
	RequestReceivedEventType  = "request.received"
	RequestProcessedEventType = "request.processed"
)

// ServerStartedEvent 表示服务器启动事件
type ServerStartedEvent struct {
	*BaseDomainEvent
	ServerID string
	Address  string
}

// NewServerStartedEvent 创建一个新的服务器启动事件
func NewServerStartedEvent(serverID, address string) *ServerStartedEvent {
	return &ServerStartedEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			ServerStartedEventType,
			serverID,
			"WebServer",
			1,
			map[string]interface{}{
				"serverID": serverID,
				"address":  address,
			},
		),
		ServerID: serverID,
		Address:  address,
	}
}

// ServerStoppedEvent 表示服务器停止事件
type ServerStoppedEvent struct {
	*BaseDomainEvent
	ServerID string
}

// NewServerStoppedEvent 创建一个新的服务器停止事件
func NewServerStoppedEvent(serverID string) *ServerStoppedEvent {
	return &ServerStoppedEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			ServerStoppedEventType,
			serverID,
			"WebServer",
			1,
			map[string]interface{}{
				"serverID": serverID,
			},
		),
		ServerID: serverID,
	}
}

// RouteRegisteredEvent 表示路由注册事件
type RouteRegisteredEvent struct {
	*BaseDomainEvent
	ServerID string
	Method   string
	Pattern  string
}

// NewRouteRegisteredEvent 创建一个新的路由注册事件
func NewRouteRegisteredEvent(serverID, method, pattern string) *RouteRegisteredEvent {
	return &RouteRegisteredEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			RouteRegisteredEventType,
			serverID,
			"WebServer",
			1,
			map[string]interface{}{
				"serverID": serverID,
				"method":   method,
				"pattern":  pattern,
			},
		),
		ServerID: serverID,
		Method:   method,
		Pattern:  pattern,
	}
}

// =============================================================================
// Session领域事件
// =============================================================================

const (
	SessionCreatedEventType   = "session.created"
	SessionExpiredEventType   = "session.expired"
	SessionDestroyedEventType = "session.destroyed"
	SessionRefreshedEventType = "session.refreshed"
)

// SessionCreatedEvent 表示会话创建事件
type SessionCreatedEvent struct {
	*BaseDomainEvent
	SessionID string
	UserID    string
	ExpiresAt time.Time
}

// NewSessionCreatedEvent 创建一个新的会话创建事件
func NewSessionCreatedEvent(sessionID, userID string, expiresAt time.Time) *SessionCreatedEvent {
	return &SessionCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			SessionCreatedEventType,
			sessionID,
			"Session",
			1,
			map[string]interface{}{
				"sessionID": sessionID,
				"userID":    userID,
				"expiresAt": expiresAt,
			},
		),
		SessionID: sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
}

// SessionExpiredEvent 表示会话过期事件
type SessionExpiredEvent struct {
	*BaseDomainEvent
	SessionID string
}

// NewSessionExpiredEvent 创建一个新的会话过期事件
func NewSessionExpiredEvent(sessionID string) *SessionExpiredEvent {
	return &SessionExpiredEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			SessionExpiredEventType,
			sessionID,
			"Session",
			1,
			map[string]interface{}{
				"sessionID": sessionID,
			},
		),
		SessionID: sessionID,
	}
}

// =============================================================================
// 文件领域事件
// =============================================================================

const (
	FileUploadedEventType   = "file.uploaded"
	FileDownloadedEventType = "file.downloaded"
	FileDeletedEventType    = "file.deleted"
)

// FileUploadedEvent 表示文件上传事件
type FileUploadedEvent struct {
	*BaseDomainEvent
	FilePath     string
	OriginalName string
	Size         int64
	ContentType  string
}

// NewFileUploadedEvent 创建一个新的文件上传事件
func NewFileUploadedEvent(filePath, originalName string, size int64, contentType string) *FileUploadedEvent {
	return &FileUploadedEvent{
		BaseDomainEvent: NewBaseDomainEvent(
			FileUploadedEventType,
			filePath,
			"FileResource",
			1,
			map[string]interface{}{
				"filePath":     filePath,
				"originalName": originalName,
				"size":         size,
				"contentType":  contentType,
			},
		),
		FilePath:     filePath,
		OriginalName: originalName,
		Size:         size,
		ContentType:  contentType,
	}
}
