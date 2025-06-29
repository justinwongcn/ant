// Package entities 包含Web服务器域的实体
package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/justinwongcn/ant/internal/domain/shared/errors"
	"github.com/justinwongcn/ant/internal/domain/webserver/valueobjects"
)

// RouteID 表示路由的唯一标识符
type RouteID string

// NewRouteID 创建一个新的RouteID
func NewRouteID() RouteID {
	return RouteID(uuid.New().String())
}

// String 返回RouteID的字符串表示
func (r RouteID) String() string {
	return string(r)
}

// HandlerFunc represents a request handler function
type HandlerFunc interface {
	Handle(ctx RequestContext) error
	Name() string
}

// RequestContext represents the context for handling a request
type RequestContext interface {
	Method() *valueobjects.HTTPMethod
	Path() string
	Parameters() map[string]string
	QueryParams() map[string]string
	Headers() map[string]string
	Body() []byte
	SetStatusCode(code *valueobjects.StatusCode)
	SetHeader(name, value string)
	SetBody(body []byte)
	Write(data []byte) error
}

// Route 表示Web服务器域中的路由实体
type Route struct {
	id          RouteID
	pattern     *valueobjects.URLPattern
	handler     HandlerFunc
	name        string
	description string
	createdAt   time.Time
	updatedAt   time.Time
	enabled     bool
	metadata    map[string]interface{}
}

// NewRoute 创建一个新的Route实体
func NewRoute(pattern *valueobjects.URLPattern, handler HandlerFunc) (*Route, error) {
	if pattern == nil {
		return nil, errors.NewValidationError("pattern", "路由模式不能为空")
	}

	if handler == nil {
		return nil, errors.NewValidationError("handler", "处理器不能为空")
	}

	now := time.Now()
	return &Route{
		id:          NewRouteID(),
		pattern:     pattern,
		handler:     handler,
		name:        generateRouteName(pattern),
		description: "",
		createdAt:   now,
		updatedAt:   now,
		enabled:     true,
		metadata:    make(map[string]interface{}),
	}, nil
}

// ID 返回路由ID
func (r *Route) ID() RouteID {
	return r.id
}

// Pattern 返回URL模式
func (r *Route) Pattern() *valueobjects.URLPattern {
	return r.pattern
}

// Handler 返回处理器函数
func (r *Route) Handler() HandlerFunc {
	return r.handler
}

// Name 返回路由名称
func (r *Route) Name() string {
	return r.name
}

// SetName 设置路由名称
func (r *Route) SetName(name string) error {
	if name == "" {
		return errors.NewValidationError("name", "路由名称不能为空")
	}

	r.name = name
	r.updatedAt = time.Now()
	return nil
}

// Description 返回路由描述
func (r *Route) Description() string {
	return r.description
}

// SetDescription 设置路由描述
func (r *Route) SetDescription(description string) {
	r.description = description
	r.updatedAt = time.Now()
}

// CreatedAt 返回创建时间
func (r *Route) CreatedAt() time.Time {
	return r.createdAt
}

// UpdatedAt 返回最后更新时间
func (r *Route) UpdatedAt() time.Time {
	return r.updatedAt
}

// IsEnabled 返回路由是否启用
func (r *Route) IsEnabled() bool {
	return r.enabled
}

// Enable 启用路由
func (r *Route) Enable() {
	if !r.enabled {
		r.enabled = true
		r.updatedAt = time.Now()
	}
}

// Disable 禁用路由
func (r *Route) Disable() {
	if r.enabled {
		r.enabled = false
		r.updatedAt = time.Now()
	}
}

// Metadata 返回路由元数据
func (r *Route) Metadata() map[string]interface{} {
	// Return a copy to prevent external modification
	metadata := make(map[string]interface{})
	for k, v := range r.metadata {
		metadata[k] = v
	}
	return metadata
}

// SetMetadata 设置元数据值
func (r *Route) SetMetadata(key string, value interface{}) error {
	if key == "" {
		return errors.NewValidationError("key", "元数据键不能为空")
	}

	r.metadata[key] = value
	r.updatedAt = time.Now()
	return nil
}

// GetMetadata 获取元数据值
func (r *Route) GetMetadata(key string) (interface{}, bool) {
	value, exists := r.metadata[key]
	return value, exists
}

// RemoveMetadata 移除元数据值
func (r *Route) RemoveMetadata(key string) {
	if _, exists := r.metadata[key]; exists {
		delete(r.metadata, key)
		r.updatedAt = time.Now()
	}
}

// Matches 检查路由是否匹配给定的方法和路径
func (r *Route) Matches(method *valueobjects.HTTPMethod, path string) bool {
	if !r.enabled {
		return false
	}
	return r.pattern.Matches(method, path)
}

// ExtractParameters 从给定路径提取路径参数
func (r *Route) ExtractParameters(path string) map[string]string {
	if !r.enabled {
		return nil
	}
	return r.pattern.ExtractParameters(path)
}

// Handle 使用此路由处理请求
func (r *Route) Handle(ctx RequestContext) error {
	if !r.enabled {
		return errors.NewDomainError(errors.ErrCodeInvalidRoute, "路由已禁用")
	}

	return r.handler.Handle(ctx)
}

// Priority 返回路由匹配的优先级
func (r *Route) Priority() int {
	return r.pattern.Priority()
}

// Equals 检查两个路由是否相等
func (r *Route) Equals(other *Route) bool {
	if other == nil {
		return false
	}
	return r.id == other.id
}

// String 返回路由的字符串表示
func (r *Route) String() string {
	return r.pattern.String()
}

// generateRouteName 为路由生成默认名称
func generateRouteName(pattern *valueobjects.URLPattern) string {
	return pattern.String()
}

// RouteBuilder 提供构建路由的流畅接口
type RouteBuilder struct {
	pattern     *valueobjects.URLPattern
	handler     HandlerFunc
	name        string
	description string
	metadata    map[string]interface{}
}

// NewRouteBuilder 创建一个新的RouteBuilder
func NewRouteBuilder() *RouteBuilder {
	return &RouteBuilder{
		metadata: make(map[string]interface{}),
	}
}

// WithPattern 设置URL模式
func (b *RouteBuilder) WithPattern(pattern *valueobjects.URLPattern) *RouteBuilder {
	b.pattern = pattern
	return b
}

// WithHandler 设置处理器
func (b *RouteBuilder) WithHandler(handler HandlerFunc) *RouteBuilder {
	b.handler = handler
	return b
}

// WithName 设置路由名称
func (b *RouteBuilder) WithName(name string) *RouteBuilder {
	b.name = name
	return b
}

// WithDescription 设置路由描述
func (b *RouteBuilder) WithDescription(description string) *RouteBuilder {
	b.description = description
	return b
}

// WithMetadata 添加元数据
func (b *RouteBuilder) WithMetadata(key string, value interface{}) *RouteBuilder {
	b.metadata[key] = value
	return b
}

// Build 构建路由
func (b *RouteBuilder) Build() (*Route, error) {
	route, err := NewRoute(b.pattern, b.handler)
	if err != nil {
		return nil, err
	}

	if b.name != "" {
		if err := route.SetName(b.name); err != nil {
			return nil, err
		}
	}

	if b.description != "" {
		route.SetDescription(b.description)
	}

	for k, v := range b.metadata {
		if err := route.SetMetadata(k, v); err != nil {
			return nil, err
		}
	}

	return route, nil
}
