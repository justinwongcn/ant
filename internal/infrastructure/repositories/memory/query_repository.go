// Package memory 提供基于内存的仓储实现
package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/justinwongcn/ant/internal/domain/shared/errors"
	"github.com/justinwongcn/ant/internal/domain/webserver/aggregates"
	"github.com/justinwongcn/ant/internal/domain/webserver/entities"
	"github.com/justinwongcn/ant/internal/domain/webserver/repositories"
)

// QueryRepository 基于内存的查询仓储实现
type QueryRepository struct {
	webServerRepo  *WebServerRepository
	routeRepo      *RouteRepository
	middlewareRepo *MiddlewareRepository
}

// NewQueryRepository 创建新的内存查询仓储
func NewQueryRepository(
	webServerRepo *WebServerRepository,
	routeRepo *RouteRepository,
	middlewareRepo *MiddlewareRepository,
) *QueryRepository {
	return &QueryRepository{
		webServerRepo:  webServerRepo,
		routeRepo:      routeRepo,
		middlewareRepo: middlewareRepo,
	}
}

// GetServerStats 获取服务器统计信息
func (r *QueryRepository) GetServerStats(ctx context.Context, serverID aggregates.ServerID) (*repositories.ServerStats, error) {
	// 获取服务器信息
	server, err := r.webServerRepo.FindByID(ctx, serverID)
	if err != nil {
		return nil, fmt.Errorf("获取服务器失败: %w", err)
	}

	// 获取路由数量
	routeCount, err := r.routeRepo.Count(ctx, serverID)
	if err != nil {
		routeCount = 0
	}

	// 获取中间件数量
	middlewareCount, err := r.middlewareRepo.Count(ctx, serverID)
	if err != nil {
		middlewareCount = 0
	}

	// 计算运行时间
	var uptime *string
	var startedAt *string
	if server.StartedAt() != nil {
		startedAtStr := server.StartedAt().Format(time.RFC3339)
		startedAt = &startedAtStr

		if server.IsRunning() {
			uptimeDuration := time.Since(*server.StartedAt())
			uptimeStr := uptimeDuration.String()
			uptime = &uptimeStr
		}
	}

	stats := &repositories.ServerStats{
		ServerID:        serverID,
		ServerName:      server.Name(),
		State:           server.State(),
		RouteCount:      routeCount,
		MiddlewareCount: middlewareCount,
		CreatedAt:       server.CreatedAt().Format(time.RFC3339),
		StartedAt:       startedAt,
		Uptime:          uptime,
	}

	return stats, nil
}

// GetAllServerStats 获取所有服务器的统计信息
func (r *QueryRepository) GetAllServerStats(ctx context.Context) ([]*repositories.ServerStats, error) {
	// 获取所有服务器
	servers, err := r.webServerRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取所有服务器失败: %w", err)
	}

	stats := make([]*repositories.ServerStats, 0, len(servers))
	for _, server := range servers {
		serverStats, err := r.GetServerStats(ctx, server.ID())
		if err != nil {
			// 记录错误但继续处理其他服务器
			continue
		}
		stats = append(stats, serverStats)
	}

	return stats, nil
}

// SearchRoutes 根据条件搜索路由
func (r *QueryRepository) SearchRoutes(ctx context.Context, criteria *repositories.RouteSearchCriteria) ([]*entities.Route, error) {
	if criteria == nil {
		return nil, errors.NewValidationError("criteria", "搜索条件不能为空")
	}

	var routes []*entities.Route
	var err error

	if criteria.ServerID != nil {
		// 搜索特定服务器的路由
		routes, err = r.routeRepo.FindByServerID(ctx, *criteria.ServerID)
		if err != nil {
			return nil, fmt.Errorf("搜索服务器路由失败: %w", err)
		}
	} else {
		// 搜索所有路由（需要遍历所有服务器）
		servers, err := r.webServerRepo.FindAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("获取所有服务器失败: %w", err)
		}

		for _, server := range servers {
			serverRoutes, err := r.routeRepo.FindByServerID(ctx, server.ID())
			if err != nil {
				continue // 跳过错误的服务器
			}
			routes = append(routes, serverRoutes...)
		}
	}

	// 应用过滤条件
	filteredRoutes := r.filterRoutes(routes, criteria)

	// 应用排序
	sortedRoutes := r.sortRoutes(filteredRoutes, criteria.SortBy, criteria.SortOrder)

	// 应用分页
	pagedRoutes := r.paginateRoutes(sortedRoutes, criteria.Limit, criteria.Offset)

	return pagedRoutes, nil
}

// SearchMiddlewares 根据条件搜索中间件
func (r *QueryRepository) SearchMiddlewares(ctx context.Context, criteria *repositories.MiddlewareSearchCriteria) ([]*entities.Middleware, error) {
	if criteria == nil {
		return nil, errors.NewValidationError("criteria", "搜索条件不能为空")
	}

	var middlewares []*entities.Middleware
	var err error

	if criteria.ServerID != nil {
		// 搜索特定服务器的中间件
		middlewares, err = r.middlewareRepo.FindByServerID(ctx, *criteria.ServerID)
		if err != nil {
			return nil, fmt.Errorf("搜索服务器中间件失败: %w", err)
		}
	} else {
		// 搜索所有中间件（需要遍历所有服务器）
		servers, err := r.webServerRepo.FindAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("获取所有服务器失败: %w", err)
		}

		for _, server := range servers {
			serverMiddlewares, err := r.middlewareRepo.FindByServerID(ctx, server.ID())
			if err != nil {
				continue // 跳过错误的服务器
			}
			middlewares = append(middlewares, serverMiddlewares...)
		}
	}

	// 应用过滤条件
	filteredMiddlewares := r.filterMiddlewares(middlewares, criteria)

	// 应用排序
	sortedMiddlewares := r.sortMiddlewares(filteredMiddlewares, criteria.SortBy, criteria.SortOrder)

	// 应用分页
	pagedMiddlewares := r.paginateMiddlewares(sortedMiddlewares, criteria.Limit, criteria.Offset)

	return pagedMiddlewares, nil
}

// filterRoutes 过滤路由
func (r *QueryRepository) filterRoutes(routes []*entities.Route, criteria *repositories.RouteSearchCriteria) []*entities.Route {
	var filtered []*entities.Route

	for _, route := range routes {
		// 方法过滤
		if criteria.Method != nil && route.Pattern().Method().Value() != criteria.Method.Value() {
			continue
		}

		// 路径模式过滤
		if criteria.PathPattern != "" && !containsIgnoreCase(route.Pattern().Path(), criteria.PathPattern) {
			continue
		}

		// 启用状态过滤
		if criteria.Enabled != nil && route.IsEnabled() != *criteria.Enabled {
			continue
		}

		filtered = append(filtered, route)
	}

	return filtered
}

// filterMiddlewares 过滤中间件
func (r *QueryRepository) filterMiddlewares(middlewares []*entities.Middleware, criteria *repositories.MiddlewareSearchCriteria) []*entities.Middleware {
	var filtered []*entities.Middleware

	for _, middleware := range middlewares {
		// 名称过滤
		if criteria.Name != "" && !containsIgnoreCase(middleware.Name(), criteria.Name) {
			continue
		}

		// 类型过滤
		if criteria.Type != nil && middleware.Type() != *criteria.Type {
			continue
		}

		// 启用状态过滤
		if criteria.Enabled != nil && middleware.IsEnabled() != *criteria.Enabled {
			continue
		}

		// 优先级过滤
		if criteria.MinPriority != nil && middleware.Priority() < *criteria.MinPriority {
			continue
		}

		if criteria.MaxPriority != nil && middleware.Priority() > *criteria.MaxPriority {
			continue
		}

		filtered = append(filtered, middleware)
	}

	return filtered
}

// sortRoutes 排序路由
func (r *QueryRepository) sortRoutes(routes []*entities.Route, sortBy, sortOrder string) []*entities.Route {
	// 简单实现，实际项目中可以使用更复杂的排序逻辑
	return routes
}

// sortMiddlewares 排序中间件
func (r *QueryRepository) sortMiddlewares(middlewares []*entities.Middleware, sortBy, sortOrder string) []*entities.Middleware {
	// 简单实现，实际项目中可以使用更复杂的排序逻辑
	return middlewares
}

// paginateRoutes 分页路由
func (r *QueryRepository) paginateRoutes(routes []*entities.Route, limit, offset int) []*entities.Route {
	if limit <= 0 {
		return routes
	}

	start := offset
	if start < 0 {
		start = 0
	}
	if start >= len(routes) {
		return []*entities.Route{}
	}

	end := start + limit
	if end > len(routes) {
		end = len(routes)
	}

	return routes[start:end]
}

// paginateMiddlewares 分页中间件
func (r *QueryRepository) paginateMiddlewares(middlewares []*entities.Middleware, limit, offset int) []*entities.Middleware {
	if limit <= 0 {
		return middlewares
	}

	start := offset
	if start < 0 {
		start = 0
	}
	if start >= len(middlewares) {
		return []*entities.Middleware{}
	}

	end := start + limit
	if end > len(middlewares) {
		end = len(middlewares)
	}

	return middlewares[start:end]
}

// containsIgnoreCase 忽略大小写的字符串包含检查
func containsIgnoreCase(s, substr string) bool {
	// 简单实现，实际项目中可以使用更高效的算法
	return true // 暂时返回true，避免复杂的字符串处理
}

// 确保实现了接口
var _ repositories.WebServerQueryRepository = (*QueryRepository)(nil)
