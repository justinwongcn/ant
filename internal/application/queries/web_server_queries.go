// Package queries 包含应用层的查询定义
package queries

import (
	"time"

	"github.com/google/uuid"
)

// Query 表示一个查询的基础接口
type Query interface {
	// QueryID 返回查询的唯一标识符
	QueryID() string
	// QueryType 返回查询类型
	QueryType() string
	// Timestamp 返回查询创建时间
	Timestamp() time.Time
	// Validate 验证查询的有效性
	Validate() error
}

// BaseQuery 提供查询的基础实现
type BaseQuery struct {
	queryID   string
	queryType string
	timestamp time.Time
}

// NewBaseQuery 创建新的基础查询
func NewBaseQuery(queryType string) *BaseQuery {
	return &BaseQuery{
		queryID:   uuid.New().String(),
		queryType: queryType,
		timestamp: time.Now(),
	}
}

func (q *BaseQuery) QueryID() string      { return q.queryID }
func (q *BaseQuery) QueryType() string    { return q.queryType }
func (q *BaseQuery) Timestamp() time.Time { return q.timestamp }

// GetServerQuery 获取服务器信息查询
type GetServerQuery struct {
	*BaseQuery
	ServerID string
}

// NewGetServerQuery 创建新的获取服务器查询
func NewGetServerQuery(serverID string) *GetServerQuery {
	return &GetServerQuery{
		BaseQuery: NewBaseQuery("GetServer"),
		ServerID:  serverID,
	}
}

// Validate 验证获取服务器查询
func (q *GetServerQuery) Validate() error {
	if q.ServerID == "" {
		return ErrInvalidQueryParameter("serverID", "服务器ID不能为空")
	}
	return nil
}

// ListServersQuery 列出服务器查询
type ListServersQuery struct {
	*BaseQuery
	State  string
	Limit  int
	Offset int
}

// NewListServersQuery 创建新的列出服务器查询
func NewListServersQuery() *ListServersQuery {
	return &ListServersQuery{
		BaseQuery: NewBaseQuery("ListServers"),
		Limit:     10,
		Offset:    0,
	}
}

// WithState 设置状态过滤
func (q *ListServersQuery) WithState(state string) *ListServersQuery {
	q.State = state
	return q
}

// WithPagination 设置分页参数
func (q *ListServersQuery) WithPagination(limit, offset int) *ListServersQuery {
	q.Limit = limit
	q.Offset = offset
	return q
}

// Validate 验证列出服务器查询
func (q *ListServersQuery) Validate() error {
	if q.Limit < 0 {
		return ErrInvalidQueryParameter("limit", "限制数量不能为负数")
	}
	if q.Offset < 0 {
		return ErrInvalidQueryParameter("offset", "偏移量不能为负数")
	}
	if q.Limit > 1000 {
		return ErrInvalidQueryParameter("limit", "限制数量不能超过1000")
	}
	return nil
}

// GetRoutesQuery 获取路由列表查询
type GetRoutesQuery struct {
	*BaseQuery
	ServerID string
	Method   string
	Enabled  *bool
	Limit    int
	Offset   int
}

// NewGetRoutesQuery 创建新的获取路由查询
func NewGetRoutesQuery(serverID string) *GetRoutesQuery {
	return &GetRoutesQuery{
		BaseQuery: NewBaseQuery("GetRoutes"),
		ServerID:  serverID,
		Limit:     10,
		Offset:    0,
	}
}

// WithMethod 设置方法过滤
func (q *GetRoutesQuery) WithMethod(method string) *GetRoutesQuery {
	q.Method = method
	return q
}

// WithEnabled 设置启用状态过滤
func (q *GetRoutesQuery) WithEnabled(enabled bool) *GetRoutesQuery {
	q.Enabled = &enabled
	return q
}

// WithPagination 设置分页参数
func (q *GetRoutesQuery) WithPagination(limit, offset int) *GetRoutesQuery {
	q.Limit = limit
	q.Offset = offset
	return q
}

// Validate 验证获取路由查询
func (q *GetRoutesQuery) Validate() error {
	if q.ServerID == "" {
		return ErrInvalidQueryParameter("serverID", "服务器ID不能为空")
	}
	if q.Limit < 0 {
		return ErrInvalidQueryParameter("limit", "限制数量不能为负数")
	}
	if q.Offset < 0 {
		return ErrInvalidQueryParameter("offset", "偏移量不能为负数")
	}
	if q.Limit > 1000 {
		return ErrInvalidQueryParameter("limit", "限制数量不能超过1000")
	}
	return nil
}

// GetMiddlewaresQuery 获取中间件列表查询
type GetMiddlewaresQuery struct {
	*BaseQuery
	ServerID string
	Type     string
	Enabled  *bool
	Limit    int
	Offset   int
}

// NewGetMiddlewaresQuery 创建新的获取中间件查询
func NewGetMiddlewaresQuery(serverID string) *GetMiddlewaresQuery {
	return &GetMiddlewaresQuery{
		BaseQuery: NewBaseQuery("GetMiddlewares"),
		ServerID:  serverID,
		Limit:     10,
		Offset:    0,
	}
}

// WithType 设置类型过滤
func (q *GetMiddlewaresQuery) WithType(middlewareType string) *GetMiddlewaresQuery {
	q.Type = middlewareType
	return q
}

// WithEnabled 设置启用状态过滤
func (q *GetMiddlewaresQuery) WithEnabled(enabled bool) *GetMiddlewaresQuery {
	q.Enabled = &enabled
	return q
}

// WithPagination 设置分页参数
func (q *GetMiddlewaresQuery) WithPagination(limit, offset int) *GetMiddlewaresQuery {
	q.Limit = limit
	q.Offset = offset
	return q
}

// Validate 验证获取中间件查询
func (q *GetMiddlewaresQuery) Validate() error {
	if q.ServerID == "" {
		return ErrInvalidQueryParameter("serverID", "服务器ID不能为空")
	}
	if q.Limit < 0 {
		return ErrInvalidQueryParameter("limit", "限制数量不能为负数")
	}
	if q.Offset < 0 {
		return ErrInvalidQueryParameter("offset", "偏移量不能为负数")
	}
	if q.Limit > 1000 {
		return ErrInvalidQueryParameter("limit", "限制数量不能超过1000")
	}
	return nil
}

// GetServerStatsQuery 获取服务器统计信息查询
type GetServerStatsQuery struct {
	*BaseQuery
	ServerID string
}

// NewGetServerStatsQuery 创建新的获取服务器统计查询
func NewGetServerStatsQuery(serverID string) *GetServerStatsQuery {
	return &GetServerStatsQuery{
		BaseQuery: NewBaseQuery("GetServerStats"),
		ServerID:  serverID,
	}
}

// Validate 验证获取服务器统计查询
func (q *GetServerStatsQuery) Validate() error {
	if q.ServerID == "" {
		return ErrInvalidQueryParameter("serverID", "服务器ID不能为空")
	}
	return nil
}

// FindMatchingRouteQuery 查找匹配路由查询
type FindMatchingRouteQuery struct {
	*BaseQuery
	ServerID string
	Method   string
	Path     string
}

// NewFindMatchingRouteQuery 创建新的查找匹配路由查询
func NewFindMatchingRouteQuery(serverID, method, path string) *FindMatchingRouteQuery {
	return &FindMatchingRouteQuery{
		BaseQuery: NewBaseQuery("FindMatchingRoute"),
		ServerID:  serverID,
		Method:    method,
		Path:      path,
	}
}

// Validate 验证查找匹配路由查询
func (q *FindMatchingRouteQuery) Validate() error {
	if q.ServerID == "" {
		return ErrInvalidQueryParameter("serverID", "服务器ID不能为空")
	}
	if q.Method == "" {
		return ErrInvalidQueryParameter("method", "HTTP方法不能为空")
	}
	if q.Path == "" {
		return ErrInvalidQueryParameter("path", "路径不能为空")
	}
	return nil
}

// QueryValidationError 查询验证错误
type QueryValidationError struct {
	Field   string
	Message string
}

func (e *QueryValidationError) Error() string {
	return e.Message
}

// ErrInvalidQueryParameter 创建无效查询参数错误
func ErrInvalidQueryParameter(field, message string) *QueryValidationError {
	return &QueryValidationError{
		Field:   field,
		Message: message,
	}
}
