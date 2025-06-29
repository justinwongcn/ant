// Package entities 包含Web服务器域的实体
package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/justinwongcn/ant/internal/domain/shared/errors"
)

// MiddlewareID 表示中间件的唯一标识符
type MiddlewareID string

// NewMiddlewareID 创建一个新的MiddlewareID
func NewMiddlewareID() MiddlewareID {
	return MiddlewareID(uuid.New().String())
}

// String 返回MiddlewareID的字符串表示
func (m MiddlewareID) String() string {
	return string(m)
}

// MiddlewareFunc 表示中间件函数
type MiddlewareFunc interface {
	Process(next HandlerFunc) HandlerFunc
	Name() string
}

// MiddlewareType 表示中间件类型
type MiddlewareType int

const (
	MiddlewareTypeGeneral MiddlewareType = iota
	MiddlewareTypeAuth
	MiddlewareTypeLogging
	MiddlewareTypeRecovery
	MiddlewareTypeRateLimit
	MiddlewareTypeCORS
	MiddlewareTypeCompression
	MiddlewareTypeCache
)

// String 返回MiddlewareType的字符串表示
func (m MiddlewareType) String() string {
	switch m {
	case MiddlewareTypeGeneral:
		return "general"
	case MiddlewareTypeAuth:
		return "auth"
	case MiddlewareTypeLogging:
		return "logging"
	case MiddlewareTypeRecovery:
		return "recovery"
	case MiddlewareTypeRateLimit:
		return "ratelimit"
	case MiddlewareTypeCORS:
		return "cors"
	case MiddlewareTypeCompression:
		return "compression"
	case MiddlewareTypeCache:
		return "cache"
	default:
		return "unknown"
	}
}

// Middleware 表示Web服务器域中的中间件实体
type Middleware struct {
	id             MiddlewareID
	name           string
	description    string
	middlewareType MiddlewareType
	function       MiddlewareFunc
	priority       int
	enabled        bool
	createdAt      time.Time
	updatedAt      time.Time
	metadata       map[string]interface{}
}

// NewMiddleware 创建一个新的Middleware实体
func NewMiddleware(name string, middlewareType MiddlewareType, function MiddlewareFunc, priority int) (*Middleware, error) {
	if name == "" {
		return nil, errors.NewValidationError("name", "中间件名称不能为空")
	}

	if function == nil {
		return nil, errors.NewValidationError("function", "中间件函数不能为空")
	}

	now := time.Now()
	return &Middleware{
		id:             NewMiddlewareID(),
		name:           name,
		description:    "",
		middlewareType: middlewareType,
		function:       function,
		priority:       priority,
		enabled:        true,
		createdAt:      now,
		updatedAt:      now,
		metadata:       make(map[string]interface{}),
	}, nil
}

// ID 返回中间件ID
func (m *Middleware) ID() MiddlewareID {
	return m.id
}

// Name 返回中间件名称
func (m *Middleware) Name() string {
	return m.name
}

// SetName 设置中间件名称
func (m *Middleware) SetName(name string) error {
	if name == "" {
		return errors.NewValidationError("name", "中间件名称不能为空")
	}

	m.name = name
	m.updatedAt = time.Now()
	return nil
}

// Description 返回中间件描述
func (m *Middleware) Description() string {
	return m.description
}

// SetDescription 设置中间件描述
func (m *Middleware) SetDescription(description string) {
	m.description = description
	m.updatedAt = time.Now()
}

// Type 返回中间件类型
func (m *Middleware) Type() MiddlewareType {
	return m.middlewareType
}

// Function 返回中间件函数
func (m *Middleware) Function() MiddlewareFunc {
	return m.function
}

// Priority 返回中间件优先级
func (m *Middleware) Priority() int {
	return m.priority
}

// SetPriority 设置中间件优先级
func (m *Middleware) SetPriority(priority int) {
	m.priority = priority
	m.updatedAt = time.Now()
}

// IsEnabled 返回中间件是否启用
func (m *Middleware) IsEnabled() bool {
	return m.enabled
}

// Enable 启用中间件
func (m *Middleware) Enable() {
	if !m.enabled {
		m.enabled = true
		m.updatedAt = time.Now()
	}
}

// Disable 禁用中间件
func (m *Middleware) Disable() {
	if m.enabled {
		m.enabled = false
		m.updatedAt = time.Now()
	}
}

// CreatedAt 返回创建时间
func (m *Middleware) CreatedAt() time.Time {
	return m.createdAt
}

// UpdatedAt 返回最后更新时间
func (m *Middleware) UpdatedAt() time.Time {
	return m.updatedAt
}

// Metadata 返回中间件元数据
func (m *Middleware) Metadata() map[string]interface{} {
	// Return a copy to prevent external modification
	metadata := make(map[string]interface{})
	for k, v := range m.metadata {
		metadata[k] = v
	}
	return metadata
}

// SetMetadata 设置元数据值
func (m *Middleware) SetMetadata(key string, value interface{}) error {
	if key == "" {
		return errors.NewValidationError("key", "元数据键不能为空")
	}

	m.metadata[key] = value
	m.updatedAt = time.Now()
	return nil
}

// GetMetadata 获取元数据值
func (m *Middleware) GetMetadata(key string) (interface{}, bool) {
	value, exists := m.metadata[key]
	return value, exists
}

// RemoveMetadata 移除元数据值
func (m *Middleware) RemoveMetadata(key string) {
	if _, exists := m.metadata[key]; exists {
		delete(m.metadata, key)
		m.updatedAt = time.Now()
	}
}

// Process 使用此中间件处理请求
func (m *Middleware) Process(next HandlerFunc) HandlerFunc {
	if !m.enabled {
		return next
	}

	return m.function.Process(next)
}

// Equals 检查两个中间件是否相等
func (m *Middleware) Equals(other *Middleware) bool {
	if other == nil {
		return false
	}
	return m.id == other.id
}

// String 返回中间件的字符串表示
func (m *Middleware) String() string {
	return m.name
}

// MiddlewareBuilder 提供构建中间件的流畅接口
type MiddlewareBuilder struct {
	name           string
	description    string
	middlewareType MiddlewareType
	function       MiddlewareFunc
	priority       int
	metadata       map[string]interface{}
}

// NewMiddlewareBuilder 创建一个新的MiddlewareBuilder
func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		middlewareType: MiddlewareTypeGeneral,
		priority:       0,
		metadata:       make(map[string]interface{}),
	}
}

// WithName 设置中间件名称
func (b *MiddlewareBuilder) WithName(name string) *MiddlewareBuilder {
	b.name = name
	return b
}

// WithDescription 设置中间件描述
func (b *MiddlewareBuilder) WithDescription(description string) *MiddlewareBuilder {
	b.description = description
	return b
}

// WithType 设置中间件类型
func (b *MiddlewareBuilder) WithType(middlewareType MiddlewareType) *MiddlewareBuilder {
	b.middlewareType = middlewareType
	return b
}

// WithFunction 设置中间件函数
func (b *MiddlewareBuilder) WithFunction(function MiddlewareFunc) *MiddlewareBuilder {
	b.function = function
	return b
}

// WithPriority 设置中间件优先级
func (b *MiddlewareBuilder) WithPriority(priority int) *MiddlewareBuilder {
	b.priority = priority
	return b
}

// WithMetadata 添加元数据
func (b *MiddlewareBuilder) WithMetadata(key string, value interface{}) *MiddlewareBuilder {
	b.metadata[key] = value
	return b
}

// Build 构建中间件
func (b *MiddlewareBuilder) Build() (*Middleware, error) {
	middleware, err := NewMiddleware(b.name, b.middlewareType, b.function, b.priority)
	if err != nil {
		return nil, err
	}

	if b.description != "" {
		middleware.SetDescription(b.description)
	}

	for k, v := range b.metadata {
		if err := middleware.SetMetadata(k, v); err != nil {
			return nil, err
		}
	}

	return middleware, nil
}

// PredefinedPriorities 定义常见的中间件优先级
const (
	PriorityRecovery    = 1000 // 最高优先级 - 应该最先运行
	PriorityLogging     = 900
	PriorityCORS        = 800
	PriorityAuth        = 700
	PriorityRateLimit   = 600
	PriorityCompression = 500
	PriorityCache       = 400
	PriorityGeneral     = 100 // 最低优先级 - 应该最后运行
)
