// Package e2e 包含端到端测试
package e2e

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/justinwongcn/ant/internal/application/dto"
	"github.com/justinwongcn/ant/internal/application/handlers"
	"github.com/justinwongcn/ant/internal/application/services"
	"github.com/justinwongcn/ant/internal/domain/shared/events"
	infraEvents "github.com/justinwongcn/ant/internal/infrastructure/events"
	infraHandlers "github.com/justinwongcn/ant/internal/infrastructure/handlers"
	"github.com/justinwongcn/ant/internal/infrastructure/repositories/memory"
)

// TestDDDArchitectureE2E 端到端测试DDD架构
func TestDDDArchitectureE2E(t *testing.T) {
	// 1. 设置完整的DDD架构
	t.Log("设置DDD架构组件...")

	// 基础设施层
	repoManager := memory.NewRepositoryManager()
	eventBus := infraEvents.NewMemoryEventBus()
	handlerRegistry := infraHandlers.NewHandlerRegistry()

	// 注册事件处理器
	loggingHandler := infraEvents.NewLoggingEventHandler(
		events.ServerStartedEventType,
		events.ServerStoppedEventType,
		events.RouteRegisteredEventType,
	)
	err := eventBus.Subscribe(loggingHandler)
	require.NoError(t, err)

	// 注册业务处理器
	helloHandler := infraHandlers.NewSimpleRouteHandler("hello-handler", "Hello, World!")
	err = handlerRegistry.RegisterRouteHandler("hello-handler", helloHandler)
	require.NoError(t, err)

	userHandler := infraHandlers.NewSimpleRouteHandler("user-handler", "User API")
	err = handlerRegistry.RegisterRouteHandler("user-handler", userHandler)
	require.NoError(t, err)

	loggingMiddleware := infraHandlers.NewSimpleMiddlewareHandler("logging-middleware")
	err = handlerRegistry.RegisterMiddlewareHandler("logging-middleware", loggingMiddleware)
	require.NoError(t, err)

	authMiddleware := infraHandlers.NewSimpleMiddlewareHandler("auth-middleware")
	err = handlerRegistry.RegisterMiddlewareHandler("auth-middleware", authMiddleware)
	require.NoError(t, err)

	// 应用层
	commandHandler := handlers.NewWebServerCommandHandler(
		repoManager.WebServer(),
		repoManager.Route(),
		repoManager.Middleware(),
		eventBus,
		handlerRegistry,
	)

	queryHandler := handlers.NewWebServerQueryHandler(
		repoManager.WebServer(),
		repoManager.Route(),
		repoManager.Middleware(),
		repoManager.Query(),
	)

	appService := services.NewWebServerApplicationService(commandHandler, queryHandler)

	ctx := context.Background()

	// 2. 测试服务器生命周期管理
	t.Log("测试服务器生命周期管理...")

	// 创建服务器
	createReq := &dto.CreateWebServerRequest{
		Name:    "e2e-test-server",
		Address: ":8888",
	}

	serverResp, err := appService.CreateServer(ctx, createReq)
	require.NoError(t, err)
	assert.NotEmpty(t, serverResp.ServerID)
	assert.Equal(t, "e2e-test-server", serverResp.Name)
	assert.Equal(t, ":8888", serverResp.Address)
	assert.Equal(t, "stopped", serverResp.State)

	serverID := serverResp.ServerID

	// 3. 测试路由管理
	t.Log("测试路由管理...")

	// 注册多个路由
	routes := []struct {
		method      string
		path        string
		handlerName string
		name        string
	}{
		{"GET", "/", "hello-handler", "root"},
		{"GET", "/hello", "hello-handler", "hello"},
		{"GET", "/users", "user-handler", "users-list"},
		{"GET", "/users/{id}", "user-handler", "user-detail"},
		{"POST", "/users", "user-handler", "user-create"},
		{"PUT", "/users/{id}", "user-handler", "user-update"},
		{"DELETE", "/users/{id}", "user-handler", "user-delete"},
	}

	routeIDs := make([]string, 0, len(routes))
	for _, route := range routes {
		routeReq := &dto.RegisterRouteRequest{
			ServerID:    serverID,
			Method:      route.method,
			Path:        route.path,
			HandlerName: route.handlerName,
			Name:        route.name,
			Description: fmt.Sprintf("%s route for %s", route.method, route.path),
		}

		routeResp, err := appService.RegisterRoute(ctx, routeReq)
		require.NoError(t, err)
		assert.NotEmpty(t, routeResp.RouteID)
		routeIDs = append(routeIDs, routeResp.RouteID)
	}

	// 验证路由数量
	getRoutesReq := &dto.GetRoutesRequest{
		ServerID: serverID,
		Limit:    20,
		Offset:   0,
	}

	routesResp, err := appService.GetRoutes(ctx, getRoutesReq)
	require.NoError(t, err)
	assert.Equal(t, len(routes), routesResp.Total)
	assert.Len(t, routesResp.Routes, len(routes))

	// 4. 测试中间件管理
	t.Log("测试中间件管理...")

	// 添加多个中间件
	middlewares := []struct {
		name           string
		middlewareType string
		handlerName    string
		priority       int
	}{
		{"logging", "logging", "logging-middleware", 900},
		{"auth", "auth", "auth-middleware", 700},
	}

	middlewareIDs := make([]string, 0, len(middlewares))
	for _, middleware := range middlewares {
		middlewareReq := &dto.AddMiddlewareRequest{
			ServerID:    serverID,
			Name:        middleware.name,
			Type:        middleware.middlewareType,
			Priority:    middleware.priority,
			HandlerName: middleware.handlerName,
			Description: fmt.Sprintf("%s middleware", middleware.name),
		}

		middlewareResp, err := appService.AddMiddleware(ctx, middlewareReq)
		require.NoError(t, err)
		assert.NotEmpty(t, middlewareResp.MiddlewareID)
		middlewareIDs = append(middlewareIDs, middlewareResp.MiddlewareID)
	}

	// 验证中间件数量
	getMiddlewaresReq := &dto.GetMiddlewaresRequest{
		ServerID: serverID,
		Limit:    20,
		Offset:   0,
	}

	middlewaresResp, err := appService.GetMiddlewares(ctx, getMiddlewaresReq)
	require.NoError(t, err)
	assert.Equal(t, len(middlewares), middlewaresResp.Total)
	assert.Len(t, middlewaresResp.Middlewares, len(middlewares))

	// 5. 测试服务器状态管理
	t.Log("测试服务器状态管理...")

	// 启动服务器
	startReq := &dto.StartServerRequest{
		ServerID: serverID,
	}

	startResp, err := appService.StartServer(ctx, startReq)
	require.NoError(t, err)
	assert.Equal(t, serverID, startResp.ServerID)
	assert.Equal(t, "running", startResp.State)
	assert.NotNil(t, startResp.StartedAt)

	// 验证服务器状态
	getServerReq := &dto.GetServerRequest{
		ServerID: serverID,
	}

	getServerResp, err := appService.GetServer(ctx, getServerReq)
	require.NoError(t, err)
	assert.Equal(t, "running", getServerResp.State)
	assert.Equal(t, len(routes), getServerResp.RouteCount)
	assert.Equal(t, len(middlewares), getServerResp.MiddlewareCount)

	// 6. 测试路由匹配
	t.Log("测试路由匹配...")

	testCases := []struct {
		method         string
		path           string
		expectedRoute  string
		expectedParams map[string]string
		shouldMatch    bool
	}{
		{"GET", "/", "root", map[string]string{}, true},
		{"GET", "/hello", "hello", map[string]string{}, true},
		{"GET", "/users", "users-list", map[string]string{}, true},
		{"GET", "/users/123", "user-detail", map[string]string{"id": "123"}, true},
		{"POST", "/users", "user-create", map[string]string{}, true},
		{"PUT", "/users/456", "user-update", map[string]string{"id": "456"}, true},
		{"DELETE", "/users/789", "user-delete", map[string]string{"id": "789"}, true},
		{"GET", "/nonexistent", "", nil, false},
		{"PATCH", "/users/123", "", nil, false},
	}

	for _, tc := range testCases {
		routeInfo, params, err := appService.FindMatchingRoute(ctx, serverID, tc.method, tc.path)

		if tc.shouldMatch {
			require.NoError(t, err, "应该找到匹配的路由: %s %s", tc.method, tc.path)
			assert.Equal(t, tc.expectedRoute, routeInfo.Name)
			assert.Equal(t, tc.method, routeInfo.Method)

			if len(tc.expectedParams) > 0 {
				assert.Equal(t, tc.expectedParams, params)
			}
		} else {
			assert.Error(t, err, "不应该找到匹配的路由: %s %s", tc.method, tc.path)
		}
	}

	// 7. 测试服务器统计信息
	t.Log("测试服务器统计信息...")

	statsResp, err := appService.GetServerStats(ctx, serverID)
	require.NoError(t, err)
	assert.Equal(t, serverID, statsResp.ServerID)
	assert.Equal(t, "e2e-test-server", statsResp.ServerName)
	assert.Equal(t, "running", statsResp.State)
	assert.Equal(t, len(routes), statsResp.RouteCount)
	assert.Equal(t, len(middlewares), statsResp.MiddlewareCount)
	assert.NotEmpty(t, statsResp.StartedAt)
	assert.NotEmpty(t, statsResp.Uptime)

	// 8. 测试服务器停止
	t.Log("测试服务器停止...")

	stopReq := &dto.StopServerRequest{
		ServerID: serverID,
	}

	stopResp, err := appService.StopServer(ctx, stopReq)
	require.NoError(t, err)
	assert.Equal(t, serverID, stopResp.ServerID)
	assert.Equal(t, "stopped", stopResp.State)
	assert.NotNil(t, stopResp.StoppedAt)

	// 验证最终状态
	finalServerResp, err := appService.GetServer(ctx, getServerReq)
	require.NoError(t, err)
	assert.Equal(t, "stopped", finalServerResp.State)

	// 9. 测试多服务器场景
	t.Log("测试多服务器场景...")

	// 创建第二个服务器
	createReq2 := &dto.CreateWebServerRequest{
		Name:    "e2e-test-server-2",
		Address: ":9999",
	}

	_, err = appService.CreateServer(ctx, createReq2)
	require.NoError(t, err)

	// 列出所有服务器
	listReq := &dto.ListServersRequest{
		Limit:  10,
		Offset: 0,
	}

	listResp, err := appService.ListServers(ctx, listReq)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, listResp.Total, 2)
	assert.GreaterOrEqual(t, len(listResp.Servers), 2)

	t.Log("DDD架构端到端测试完成！")
}
