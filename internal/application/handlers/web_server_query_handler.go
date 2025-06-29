// Package handlers 包含应用层的查询处理器
package handlers

import (
	"context"
	"fmt"

	"github.com/justinwongcn/ant/internal/application/dto"
	"github.com/justinwongcn/ant/internal/application/queries"
	"github.com/justinwongcn/ant/internal/domain/webserver/aggregates"
	"github.com/justinwongcn/ant/internal/domain/webserver/entities"
	"github.com/justinwongcn/ant/internal/domain/webserver/repositories"
	"github.com/justinwongcn/ant/internal/domain/webserver/valueobjects"
)

// WebServerQueryHandler 处理Web服务器相关的查询
type WebServerQueryHandler struct {
	webServerRepo  repositories.WebServerRepository
	routeRepo      repositories.RouteRepository
	middlewareRepo repositories.MiddlewareRepository
	queryRepo      repositories.WebServerQueryRepository
}

// NewWebServerQueryHandler 创建新的Web服务器查询处理器
func NewWebServerQueryHandler(
	webServerRepo repositories.WebServerRepository,
	routeRepo repositories.RouteRepository,
	middlewareRepo repositories.MiddlewareRepository,
	queryRepo repositories.WebServerQueryRepository,
) *WebServerQueryHandler {
	return &WebServerQueryHandler{
		webServerRepo:  webServerRepo,
		routeRepo:      routeRepo,
		middlewareRepo: middlewareRepo,
		queryRepo:      queryRepo,
	}
}

// HandleGetServer 处理获取服务器查询
func (h *WebServerQueryHandler) HandleGetServer(ctx context.Context, query *queries.GetServerQuery) (*dto.GetServerResponse, error) {
	// 验证查询
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("查询验证失败: %w", err)
	}

	// 获取服务器聚合
	serverID := aggregates.ServerID(query.ServerID)
	server, err := h.webServerRepo.FindByID(ctx, serverID)
	if err != nil {
		return nil, fmt.Errorf("获取Web服务器失败: %w", err)
	}

	// 获取路由数量
	routeCount, err := h.routeRepo.Count(ctx, serverID)
	if err != nil {
		routeCount = 0 // 如果获取失败，设为0
	}

	// 获取中间件数量
	middlewareCount, err := h.middlewareRepo.Count(ctx, serverID)
	if err != nil {
		middlewareCount = 0 // 如果获取失败，设为0
	}

	// 构建响应
	response := &dto.GetServerResponse{
		ServerID:        server.ID().String(),
		Name:            server.Name(),
		Address:         server.Address(),
		State:           dto.ConvertServerStateToString(server.State()),
		RouteCount:      routeCount,
		MiddlewareCount: middlewareCount,
		CreatedAt:       server.CreatedAt(),
		UpdatedAt:       server.UpdatedAt(),
		StartedAt:       server.StartedAt(),
		StoppedAt:       server.StoppedAt(),
	}

	return response, nil
}

// HandleListServers 处理列出服务器查询
func (h *WebServerQueryHandler) HandleListServers(ctx context.Context, query *queries.ListServersQuery) (*dto.ListServersResponse, error) {
	// 验证查询
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("查询验证失败: %w", err)
	}

	var servers []*aggregates.WebServer
	var err error

	// 根据状态过滤
	if query.State != "" {
		servers, err = h.webServerRepo.FindRunning(ctx)
		if err != nil {
			return nil, fmt.Errorf("获取运行中的服务器失败: %w", err)
		}
	} else {
		servers, err = h.webServerRepo.FindAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("获取所有服务器失败: %w", err)
		}
	}

	// 应用分页
	total := len(servers)
	start := query.Offset
	end := start + query.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedServers := servers[start:end]

	// 构建响应
	serverResponses := make([]dto.GetServerResponse, len(pagedServers))
	for i, server := range pagedServers {
		// 获取路由数量
		routeCount, _ := h.routeRepo.Count(ctx, server.ID())
		// 获取中间件数量
		middlewareCount, _ := h.middlewareRepo.Count(ctx, server.ID())

		serverResponses[i] = dto.GetServerResponse{
			ServerID:        server.ID().String(),
			Name:            server.Name(),
			Address:         server.Address(),
			State:           dto.ConvertServerStateToString(server.State()),
			RouteCount:      routeCount,
			MiddlewareCount: middlewareCount,
			CreatedAt:       server.CreatedAt(),
			UpdatedAt:       server.UpdatedAt(),
			StartedAt:       server.StartedAt(),
			StoppedAt:       server.StoppedAt(),
		}
	}

	response := &dto.ListServersResponse{
		Servers: serverResponses,
		Total:   total,
		Limit:   query.Limit,
		Offset:  query.Offset,
	}

	return response, nil
}

// HandleGetRoutes 处理获取路由列表查询
func (h *WebServerQueryHandler) HandleGetRoutes(ctx context.Context, query *queries.GetRoutesQuery) (*dto.GetRoutesResponse, error) {
	// 验证查询
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("查询验证失败: %w", err)
	}

	// 获取路由列表
	serverID := aggregates.ServerID(query.ServerID)
	routes, err := h.routeRepo.FindByServerID(ctx, serverID)
	if err != nil {
		return nil, fmt.Errorf("获取路由列表失败: %w", err)
	}

	// 应用过滤
	filteredRoutes := routes
	if query.Method != "" {
		filteredRoutes = make([]*entities.Route, 0)
		for _, route := range routes {
			if route.Pattern().Method().Value() == query.Method {
				filteredRoutes = append(filteredRoutes, route)
			}
		}
	}

	if query.Enabled != nil {
		enabledRoutes := make([]*entities.Route, 0)
		for _, route := range filteredRoutes {
			if route.IsEnabled() == *query.Enabled {
				enabledRoutes = append(enabledRoutes, route)
			}
		}
		filteredRoutes = enabledRoutes
	}

	// 应用分页
	total := len(filteredRoutes)
	start := query.Offset
	end := start + query.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedRoutes := filteredRoutes[start:end]

	// 构建响应
	routeInfos := make([]dto.RouteInfo, len(pagedRoutes))
	for i, route := range pagedRoutes {
		routeInfos[i] = dto.RouteInfo{
			RouteID:     route.ID().String(),
			Method:      route.Pattern().Method().Value(),
			Path:        route.Pattern().Path(),
			Name:        route.Name(),
			Description: route.Description(),
			Priority:    route.Priority(),
			Enabled:     route.IsEnabled(),
			Metadata:    route.Metadata(),
			CreatedAt:   route.CreatedAt(),
			UpdatedAt:   route.UpdatedAt(),
		}
	}

	response := &dto.GetRoutesResponse{
		Routes: routeInfos,
		Total:  total,
		Limit:  query.Limit,
		Offset: query.Offset,
	}

	return response, nil
}

// HandleGetMiddlewares 处理获取中间件列表查询
func (h *WebServerQueryHandler) HandleGetMiddlewares(ctx context.Context, query *queries.GetMiddlewaresQuery) (*dto.GetMiddlewaresResponse, error) {
	// 验证查询
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("查询验证失败: %w", err)
	}

	// 获取中间件列表
	serverID := aggregates.ServerID(query.ServerID)
	middlewares, err := h.middlewareRepo.FindByServerID(ctx, serverID)
	if err != nil {
		return nil, fmt.Errorf("获取中间件列表失败: %w", err)
	}

	// 应用过滤
	filteredMiddlewares := middlewares
	if query.Type != "" {
		filteredMiddlewares = make([]*entities.Middleware, 0)
		for _, middleware := range middlewares {
			if middleware.Type().String() == query.Type {
				filteredMiddlewares = append(filteredMiddlewares, middleware)
			}
		}
	}

	if query.Enabled != nil {
		enabledMiddlewares := make([]*entities.Middleware, 0)
		for _, middleware := range filteredMiddlewares {
			if middleware.IsEnabled() == *query.Enabled {
				enabledMiddlewares = append(enabledMiddlewares, middleware)
			}
		}
		filteredMiddlewares = enabledMiddlewares
	}

	// 应用分页
	total := len(filteredMiddlewares)
	start := query.Offset
	end := start + query.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedMiddlewares := filteredMiddlewares[start:end]

	// 构建响应
	middlewareInfos := make([]dto.MiddlewareInfo, len(pagedMiddlewares))
	for i, middleware := range pagedMiddlewares {
		middlewareInfos[i] = dto.MiddlewareInfo{
			MiddlewareID: middleware.ID().String(),
			Name:         middleware.Name(),
			Type:         middleware.Type().String(),
			Priority:     middleware.Priority(),
			Description:  middleware.Description(),
			Enabled:      middleware.IsEnabled(),
			Metadata:     middleware.Metadata(),
			CreatedAt:    middleware.CreatedAt(),
			UpdatedAt:    middleware.UpdatedAt(),
		}
	}

	response := &dto.GetMiddlewaresResponse{
		Middlewares: middlewareInfos,
		Total:       total,
		Limit:       query.Limit,
		Offset:      query.Offset,
	}

	return response, nil
}

// HandleGetServerStats 处理获取服务器统计信息查询
func (h *WebServerQueryHandler) HandleGetServerStats(ctx context.Context, query *queries.GetServerStatsQuery) (*dto.ServerStatsResponse, error) {
	// 验证查询
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("查询验证失败: %w", err)
	}

	// 获取服务器统计信息
	serverID := aggregates.ServerID(query.ServerID)
	stats, err := h.queryRepo.GetServerStats(ctx, serverID)
	if err != nil {
		return nil, fmt.Errorf("获取服务器统计信息失败: %w", err)
	}

	// 构建响应
	response := &dto.ServerStatsResponse{
		ServerID:        stats.ServerID.String(),
		ServerName:      stats.ServerName,
		State:           stats.State.String(),
		RouteCount:      stats.RouteCount,
		MiddlewareCount: stats.MiddlewareCount,
		CreatedAt:       stats.CreatedAt,
	}

	if stats.StartedAt != nil {
		response.StartedAt = *stats.StartedAt
	}

	if stats.Uptime != nil {
		response.Uptime = *stats.Uptime
	}

	return response, nil
}

// HandleFindMatchingRoute 处理查找匹配路由查询
func (h *WebServerQueryHandler) HandleFindMatchingRoute(ctx context.Context, query *queries.FindMatchingRouteQuery) (*dto.RouteInfo, map[string]string, error) {
	// 验证查询
	if err := query.Validate(); err != nil {
		return nil, nil, fmt.Errorf("查询验证失败: %w", err)
	}

	// 获取服务器聚合
	serverID := aggregates.ServerID(query.ServerID)
	server, err := h.webServerRepo.FindByID(ctx, serverID)
	if err != nil {
		return nil, nil, fmt.Errorf("获取Web服务器失败: %w", err)
	}

	// 创建HTTP方法值对象
	method, err := valueobjects.NewHTTPMethod(query.Method)
	if err != nil {
		return nil, nil, fmt.Errorf("创建HTTP方法失败: %w", err)
	}

	// 查找匹配的路由
	route, params, err := server.FindMatchingRoute(method, query.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("查找匹配路由失败: %w", err)
	}

	// 构建响应
	routeInfo := &dto.RouteInfo{
		RouteID:     route.ID().String(),
		Method:      route.Pattern().Method().Value(),
		Path:        route.Pattern().Path(),
		Name:        route.Name(),
		Description: route.Description(),
		Priority:    route.Priority(),
		Enabled:     route.IsEnabled(),
		Metadata:    route.Metadata(),
		CreatedAt:   route.CreatedAt(),
		UpdatedAt:   route.UpdatedAt(),
	}

	return routeInfo, params, nil
}
