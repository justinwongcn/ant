// Package aggregates 包含Web服务器域的聚合根
package aggregates

import (
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/justinwongcn/ant/internal/domain/shared/errors"
	"github.com/justinwongcn/ant/internal/domain/shared/events"
	"github.com/justinwongcn/ant/internal/domain/webserver/entities"
	"github.com/justinwongcn/ant/internal/domain/webserver/valueobjects"
)

// ServerID 表示Web服务器的唯一标识符
type ServerID string

// NewServerID 创建一个新的ServerID
func NewServerID() ServerID {
	return ServerID(uuid.New().String())
}

// String 返回ServerID的字符串表示
func (s ServerID) String() string {
	return string(s)
}

// ServerState 表示Web服务器的状态
type ServerState int

const (
	ServerStateStopped ServerState = iota
	ServerStateStarting
	ServerStateRunning
	ServerStateStopping
)

// String 返回ServerState的字符串表示
func (s ServerState) String() string {
	switch s {
	case ServerStateStopped:
		return "stopped"
	case ServerStateStarting:
		return "starting"
	case ServerStateRunning:
		return "running"
	case ServerStateStopping:
		return "stopping"
	default:
		return "unknown"
	}
}

// WebServer 表示Web服务器聚合根
type WebServer struct {
	id          ServerID
	name        string
	address     string
	state       ServerState
	routes      map[string]*entities.Route // key: pattern string
	middlewares []*entities.Middleware
	createdAt   time.Time
	updatedAt   time.Time
	startedAt   *time.Time
	stoppedAt   *time.Time
	mu          sync.RWMutex
	events      []events.DomainEvent
}

// NewWebServer 创建一个新的WebServer聚合
func NewWebServer(name, address string) (*WebServer, error) {
	if name == "" {
		return nil, errors.NewValidationError("name", "服务器名称不能为空")
	}

	if address == "" {
		return nil, errors.NewValidationError("address", "服务器地址不能为空")
	}

	now := time.Now()
	server := &WebServer{
		id:          NewServerID(),
		name:        name,
		address:     address,
		state:       ServerStateStopped,
		routes:      make(map[string]*entities.Route),
		middlewares: make([]*entities.Middleware, 0),
		createdAt:   now,
		updatedAt:   now,
		events:      make([]events.DomainEvent, 0),
	}

	return server, nil
}

// ID 返回服务器ID
func (w *WebServer) ID() ServerID {
	return w.id
}

// Name 返回服务器名称
func (w *WebServer) Name() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.name
}

// SetName 设置服务器名称
func (w *WebServer) SetName(name string) error {
	if name == "" {
		return errors.NewValidationError("name", "服务器名称不能为空")
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	w.name = name
	w.updatedAt = time.Now()
	return nil
}

// Address 返回服务器地址
func (w *WebServer) Address() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.address
}

// State 返回当前服务器状态
func (w *WebServer) State() ServerState {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.state
}

// IsRunning 如果服务器正在运行则返回true
func (w *WebServer) IsRunning() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.state == ServerStateRunning
}

// CreatedAt 返回创建时间
func (w *WebServer) CreatedAt() time.Time {
	return w.createdAt
}

// UpdatedAt 返回最后更新时间
func (w *WebServer) UpdatedAt() time.Time {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.updatedAt
}

// StartedAt 返回启动时间
func (w *WebServer) StartedAt() *time.Time {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.startedAt
}

// StoppedAt 返回停止时间
func (w *WebServer) StoppedAt() *time.Time {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.stoppedAt
}

// RegisterRoute 注册一个新路由
func (w *WebServer) RegisterRoute(route *entities.Route) error {
	if route == nil {
		return errors.NewValidationError("route", "路由不能为空")
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	patternKey := route.Pattern().String()

	// 检查路由是否已存在
	if _, exists := w.routes[patternKey]; exists {
		return errors.NewDomainErrorWithCause(
			errors.ErrCodeRouteAlreadyExists,
			"路由已存在",
			errors.NewValidationError("pattern", patternKey),
		)
	}

	w.routes[patternKey] = route
	w.updatedAt = time.Now()

	// 添加领域事件
	event := events.NewRouteRegisteredEvent(
		w.id.String(),
		route.Pattern().Method().Value(),
		route.Pattern().Path(),
	)
	w.events = append(w.events, event)

	return nil
}

// UnregisterRoute 移除一个路由
func (w *WebServer) UnregisterRoute(pattern *valueobjects.URLPattern) error {
	if pattern == nil {
		return errors.NewValidationError("pattern", "路由模式不能为空")
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	patternKey := pattern.String()

	if _, exists := w.routes[patternKey]; !exists {
		return errors.NewDomainError(errors.ErrCodeInvalidRoute, "路由不存在")
	}

	delete(w.routes, patternKey)
	w.updatedAt = time.Now()

	return nil
}

// GetRoute 通过模式获取路由
func (w *WebServer) GetRoute(pattern *valueobjects.URLPattern) (*entities.Route, error) {
	if pattern == nil {
		return nil, errors.NewValidationError("pattern", "路由模式不能为空")
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	patternKey := pattern.String()
	route, exists := w.routes[patternKey]
	if !exists {
		return nil, errors.NewDomainError(errors.ErrCodeInvalidRoute, "路由不存在")
	}

	return route, nil
}

// GetRoutes 返回按优先级排序的所有路由
func (w *WebServer) GetRoutes() []*entities.Route {
	w.mu.RLock()
	defer w.mu.RUnlock()

	routes := make([]*entities.Route, 0, len(w.routes))
	for _, route := range w.routes {
		routes = append(routes, route)
	}

	// 按优先级排序(优先级高的在前)
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Priority() > routes[j].Priority()
	})

	return routes
}

// FindMatchingRoute 查找匹配给定方法和路径的第一个路由
func (w *WebServer) FindMatchingRoute(method *valueobjects.HTTPMethod, path string) (*entities.Route, map[string]string, error) {
	if method == nil {
		return nil, nil, errors.NewValidationError("method", "HTTP方法不能为空")
	}

	if path == "" {
		return nil, nil, errors.NewValidationError("path", "路径不能为空")
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	// 获取按优先级排序的路由
	routes := w.getSortedRoutes()

	for _, route := range routes {
		if route.Matches(method, path) {
			params := route.ExtractParameters(path)
			return route, params, nil
		}
	}

	return nil, nil, errors.NewDomainError(errors.ErrCodeInvalidRoute, "未找到匹配的路由")
}

// AddMiddleware 添加一个中间件
func (w *WebServer) AddMiddleware(middleware *entities.Middleware) error {
	if middleware == nil {
		return errors.NewValidationError("middleware", "中间件不能为空")
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	w.middlewares = append(w.middlewares, middleware)
	w.updatedAt = time.Now()

	// 按优先级排序中间件
	w.sortMiddlewares()

	return nil
}

// GetMiddlewares 返回按优先级排序的所有中间件
func (w *WebServer) GetMiddlewares() []*entities.Middleware {
	w.mu.RLock()
	defer w.mu.RUnlock()

	middlewares := make([]*entities.Middleware, len(w.middlewares))
	copy(middlewares, w.middlewares)
	return middlewares
}

// Start 启动Web服务器
func (w *WebServer) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.state == ServerStateRunning {
		return errors.ErrServerAlreadyRunning
	}

	w.state = ServerStateStarting
	now := time.Now()
	w.startedAt = &now
	w.stoppedAt = nil
	w.updatedAt = now

	// 切换到运行状态
	w.state = ServerStateRunning

	// 添加领域事件
	event := events.NewServerStartedEvent(w.id.String(), w.address)
	w.events = append(w.events, event)

	return nil
}

// Stop 停止Web服务器
func (w *WebServer) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.state != ServerStateRunning {
		return errors.ErrServerNotRunning
	}

	w.state = ServerStateStopping
	now := time.Now()
	w.stoppedAt = &now
	w.updatedAt = now

	// 切换到停止状态
	w.state = ServerStateStopped

	// 添加领域事件
	event := events.NewServerStoppedEvent(w.id.String())
	w.events = append(w.events, event)

	return nil
}

// GetEvents 返回并清除领域事件
func (w *WebServer) GetEvents() []events.DomainEvent {
	w.mu.Lock()
	defer w.mu.Unlock()

	events := make([]events.DomainEvent, len(w.events))
	copy(events, w.events)
	w.events = w.events[:0] // 清空事件

	return events
}

// getSortedRoutes 返回按优先级排序的路由(内部方法，假设已持有锁)
func (w *WebServer) getSortedRoutes() []*entities.Route {
	routes := make([]*entities.Route, 0, len(w.routes))
	for _, route := range w.routes {
		routes = append(routes, route)
	}

	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Priority() > routes[j].Priority()
	})

	return routes
}

// sortMiddlewares 按优先级排序中间件(内部方法，假设已持有锁)
func (w *WebServer) sortMiddlewares() {
	sort.Slice(w.middlewares, func(i, j int) bool {
		return w.middlewares[i].Priority() > w.middlewares[j].Priority()
	})
}
