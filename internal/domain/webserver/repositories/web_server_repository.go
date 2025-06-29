// Package repositories 定义了web服务器领域的仓库接口
package repositories

import (
	"context"

	"github.com/justinwongcn/ant/internal/domain/webserver/aggregates"
	"github.com/justinwongcn/ant/internal/domain/webserver/entities"
	"github.com/justinwongcn/ant/internal/domain/webserver/valueobjects"
)

// WebServerRepository 定义了web服务器持久化的接口
type WebServerRepository interface {
	// Save 保存web服务器聚合
	Save(ctx context.Context, server *aggregates.WebServer) error

	// FindByID 通过ID查找web服务器
	FindByID(ctx context.Context, id aggregates.ServerID) (*aggregates.WebServer, error)

	// FindByName 通过名称查找web服务器
	FindByName(ctx context.Context, name string) (*aggregates.WebServer, error)

	// FindAll 查找所有web服务器
	FindAll(ctx context.Context) ([]*aggregates.WebServer, error)

	// FindRunning 查找所有运行中的web服务器
	FindRunning(ctx context.Context) ([]*aggregates.WebServer, error)

	// Remove 删除web服务器
	Remove(ctx context.Context, id aggregates.ServerID) error

	// Exists 检查web服务器是否存在
	Exists(ctx context.Context, id aggregates.ServerID) (bool, error)

	// Count 返回web服务器的总数
	Count(ctx context.Context) (int, error)
}

// RouteRepository 定义了路由持久化的接口
type RouteRepository interface {
	// Save 保存路由实体
	Save(ctx context.Context, serverID aggregates.ServerID, route *entities.Route) error

	// FindByID 通过ID查找路由
	FindByID(ctx context.Context, id entities.RouteID) (*entities.Route, error)

	// FindByPattern 通过模式查找路由
	FindByPattern(ctx context.Context, serverID aggregates.ServerID, pattern *valueobjects.URLPattern) (*entities.Route, error)

	// FindByServerID 查找服务器所有路由
	FindByServerID(ctx context.Context, serverID aggregates.ServerID) ([]*entities.Route, error)

	// FindMatching 查找匹配给定方法和路径的路由
	FindMatching(ctx context.Context, serverID aggregates.ServerID, method *valueobjects.HTTPMethod, path string) ([]*entities.Route, error)

	// Remove 删除路由
	Remove(ctx context.Context, id entities.RouteID) error

	// RemoveByPattern 通过模式删除路由
	RemoveByPattern(ctx context.Context, serverID aggregates.ServerID, pattern *valueobjects.URLPattern) error

	// RemoveByServerID 删除服务器所有路由
	RemoveByServerID(ctx context.Context, serverID aggregates.ServerID) error

	// Exists 检查路由是否存在
	Exists(ctx context.Context, id entities.RouteID) (bool, error)

	// ExistsByPattern 检查路由是否通过模式存在
	ExistsByPattern(ctx context.Context, serverID aggregates.ServerID, pattern *valueobjects.URLPattern) (bool, error)

	// Count 返回服务器的路由总数
	Count(ctx context.Context, serverID aggregates.ServerID) (int, error)
}

// MiddlewareRepository 定义了中间件持久化的接口
type MiddlewareRepository interface {
	// Save 保存中间件实体
	Save(ctx context.Context, serverID aggregates.ServerID, middleware *entities.Middleware) error

	// FindByID 通过ID查找中间件
	FindByID(ctx context.Context, id entities.MiddlewareID) (*entities.Middleware, error)

	// FindByName 通过名称查找中间件
	FindByName(ctx context.Context, serverID aggregates.ServerID, name string) (*entities.Middleware, error)

	// FindByServerID 查找服务器所有中间件
	FindByServerID(ctx context.Context, serverID aggregates.ServerID) ([]*entities.Middleware, error)

	// FindByType 通过类型查找中间件
	FindByType(ctx context.Context, serverID aggregates.ServerID, middlewareType entities.MiddlewareType) ([]*entities.Middleware, error)

	// FindEnabled 查找服务器所有启用的中间件
	FindEnabled(ctx context.Context, serverID aggregates.ServerID) ([]*entities.Middleware, error)

	// Remove 删除中间件
	Remove(ctx context.Context, id entities.MiddlewareID) error

	// RemoveByName 通过名称删除中间件
	RemoveByName(ctx context.Context, serverID aggregates.ServerID, name string) error

	// RemoveByServerID 删除服务器所有中间件
	RemoveByServerID(ctx context.Context, serverID aggregates.ServerID) error

	// Exists 检查中间件是否存在
	Exists(ctx context.Context, id entities.MiddlewareID) (bool, error)

	// ExistsByName 检查中间件是否通过名称存在
	ExistsByName(ctx context.Context, serverID aggregates.ServerID, name string) (bool, error)

	// Count 返回服务器的中间件总数
	Count(ctx context.Context, serverID aggregates.ServerID) (int, error)
}

// WebServerQueryRepository 定义了web服务器数据的只读查询接口
type WebServerQueryRepository interface {
	// GetServerStats 返回服务器统计信息
	GetServerStats(ctx context.Context, serverID aggregates.ServerID) (*ServerStats, error)

	// GetAllServerStats 返回所有服务器统计信息
	GetAllServerStats(ctx context.Context) ([]*ServerStats, error)

	// SearchRoutes 根据条件搜索路由
	SearchRoutes(ctx context.Context, criteria *RouteSearchCriteria) ([]*entities.Route, error)

	// SearchMiddlewares 根据条件搜索中间件
	SearchMiddlewares(ctx context.Context, criteria *MiddlewareSearchCriteria) ([]*entities.Middleware, error)
}

// ServerStats 表示web服务器的统计信息
type ServerStats struct {
	ServerID        aggregates.ServerID
	ServerName      string
	State           aggregates.ServerState
	RouteCount      int
	MiddlewareCount int
	CreatedAt       string
	StartedAt       *string
	Uptime          *string
}

// RouteSearchCriteria 表示路由的搜索条件
type RouteSearchCriteria struct {
	ServerID    *aggregates.ServerID
	Method      *valueobjects.HTTPMethod
	PathPattern string
	Enabled     *bool
	Limit       int
	Offset      int
	SortBy      string
	SortOrder   string
}

// MiddlewareSearchCriteria 表示中间件的搜索条件
type MiddlewareSearchCriteria struct {
	ServerID    *aggregates.ServerID
	Name        string
	Type        *entities.MiddlewareType
	Enabled     *bool
	MinPriority *int
	MaxPriority *int
	Limit       int
	Offset      int
	SortBy      string
	SortOrder   string
}

// RepositoryManager 管理所有web服务器仓库
type RepositoryManager interface {
	// WebServer 返回web服务器仓库
	WebServer() WebServerRepository

	// Route 返回路由仓库
	Route() RouteRepository

	// Middleware 返回中间件仓库
	Middleware() MiddlewareRepository

	// Query 返回查询仓库
	Query() WebServerQueryRepository

	// BeginTransaction 开始一个新事务
	BeginTransaction(ctx context.Context) (Transaction, error)
}

// Transaction 表示数据库事务
type Transaction interface {
	// Commit 提交事务
	Commit() error

	// Rollback 回滚事务
	Rollback() error

	// WebServer 返回此事务中的web服务器仓库
	WebServer() WebServerRepository

	// Route 返回此事务中的路由仓库
	Route() RouteRepository

	// Middleware 返回此事务中的中间件仓库
	Middleware() MiddlewareRepository
}

// 事务相关错误
var (
	ErrTransactionAlreadyCommitted = &TransactionError{Code: "TRANSACTION_ALREADY_COMMITTED", Message: "事务已提交"}
	ErrTransactionRolledBack       = &TransactionError{Code: "TRANSACTION_ROLLED_BACK", Message: "事务已回滚"}
)

// TransactionError 事务错误
type TransactionError struct {
	Code    string
	Message string
}

func (e *TransactionError) Error() string {
	return e.Message
}
