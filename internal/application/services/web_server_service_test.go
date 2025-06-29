package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/justinwongcn/ant/internal/application/dto"
	"github.com/justinwongcn/ant/internal/application/handlers"
	infraEvents "github.com/justinwongcn/ant/internal/infrastructure/events"
	infraHandlers "github.com/justinwongcn/ant/internal/infrastructure/handlers"
	"github.com/justinwongcn/ant/internal/infrastructure/repositories/memory"
)

// setupTestService 设置测试服务
func setupTestService() WebServerService {
	// 创建基础设施组件
	repoManager := memory.NewRepositoryManager()
	eventBus := infraEvents.NewMemoryEventBus()
	handlerRegistry := infraHandlers.NewHandlerRegistry()

	// 注册测试处理器
	routeHandler := infraHandlers.NewSimpleRouteHandler("test-handler", "Test Response")
	handlerRegistry.RegisterRouteHandler("test-handler", routeHandler)

	middlewareHandler := infraHandlers.NewSimpleMiddlewareHandler("test-middleware")
	handlerRegistry.RegisterMiddlewareHandler("test-middleware", middlewareHandler)

	// 创建应用层组件
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

	return NewWebServerApplicationService(commandHandler, queryHandler)
}

func TestWebServerApplicationService_CreateServer(t *testing.T) {
	// 准备
	service := setupTestService()
	ctx := context.Background()

	req := &dto.CreateWebServerRequest{
		Name:    "test-server",
		Address: ":8080",
	}

	// 执行
	response, err := service.CreateServer(ctx, req)

	// 验证
	require.NoError(t, err)
	assert.NotEmpty(t, response.ServerID)
	assert.Equal(t, "test-server", response.Name)
	assert.Equal(t, ":8080", response.Address)
	assert.Equal(t, "stopped", response.State)
	assert.NotZero(t, response.CreatedAt)
}

func TestWebServerApplicationService_StartServer(t *testing.T) {
	// 准备
	service := setupTestService()
	ctx := context.Background()

	// 先创建服务器
	createReq := &dto.CreateWebServerRequest{
		Name:    "test-server",
		Address: ":8080",
	}

	createResp, err := service.CreateServer(ctx, createReq)
	require.NoError(t, err)

	// 启动服务器
	startReq := &dto.StartServerRequest{
		ServerID: createResp.ServerID,
	}

	// 执行
	response, err := service.StartServer(ctx, startReq)

	// 验证
	require.NoError(t, err)
	assert.Equal(t, createResp.ServerID, response.ServerID)
	assert.Equal(t, "running", response.State)
	assert.NotNil(t, response.StartedAt)
	assert.Equal(t, "服务器启动成功", response.Message)
}

func TestWebServerApplicationService_GetServer(t *testing.T) {
	// 准备
	service := setupTestService()
	ctx := context.Background()

	// 先创建服务器
	createReq := &dto.CreateWebServerRequest{
		Name:    "test-server",
		Address: ":8080",
	}

	createResp, err := service.CreateServer(ctx, createReq)
	require.NoError(t, err)

	// 获取服务器
	getReq := &dto.GetServerRequest{
		ServerID: createResp.ServerID,
	}

	// 执行
	response, err := service.GetServer(ctx, getReq)

	// 验证
	require.NoError(t, err)
	assert.Equal(t, createResp.ServerID, response.ServerID)
	assert.Equal(t, "test-server", response.Name)
	assert.Equal(t, ":8080", response.Address)
	assert.Equal(t, "stopped", response.State)
	assert.Equal(t, 0, response.RouteCount)
	assert.Equal(t, 0, response.MiddlewareCount)
}

func TestWebServerApplicationService_RegisterRoute(t *testing.T) {
	// 准备
	service := setupTestService()
	ctx := context.Background()

	// 先创建服务器
	createReq := &dto.CreateWebServerRequest{
		Name:    "test-server",
		Address: ":8080",
	}

	createResp, err := service.CreateServer(ctx, createReq)
	require.NoError(t, err)

	// 注册路由
	routeReq := &dto.RegisterRouteRequest{
		ServerID:    createResp.ServerID,
		Method:      "GET",
		Path:        "/hello",
		HandlerName: "test-handler",
		Name:        "hello-route",
		Description: "Hello route",
		Metadata:    map[string]interface{}{"version": "1.0"},
	}

	// 执行
	response, err := service.RegisterRoute(ctx, routeReq)

	// 验证
	require.NoError(t, err)
	assert.NotEmpty(t, response.RouteID)
	assert.Equal(t, createResp.ServerID, response.ServerID)
	assert.Equal(t, "GET", response.Method)
	assert.Equal(t, "/hello", response.Path)
	assert.Equal(t, "hello-route", response.Name)
	assert.Equal(t, "Hello route", response.Description)
	assert.Equal(t, "路由注册成功", response.Message)
}

func TestApplicationServiceManager(t *testing.T) {
	// 准备
	webServerService := setupTestService()
	manager := NewApplicationServiceManager(webServerService)

	// 验证
	assert.Equal(t, webServerService, manager.WebServer())
	assert.Implements(t, (*ServiceManager)(nil), manager)
	assert.Implements(t, (*WebServerService)(nil), webServerService)
}

func TestWebServerApplicationService_Integration(t *testing.T) {
	// 集成测试：测试完整的工作流程
	service := setupTestService()
	ctx := context.Background()

	// 1. 创建服务器
	createReq := &dto.CreateWebServerRequest{
		Name:    "integration-test-server",
		Address: ":9999",
	}

	serverResp, err := service.CreateServer(ctx, createReq)
	require.NoError(t, err)
	assert.NotEmpty(t, serverResp.ServerID)

	// 2. 注册路由
	routeReq := &dto.RegisterRouteRequest{
		ServerID:    serverResp.ServerID,
		Method:      "GET",
		Path:        "/api/test",
		HandlerName: "test-handler",
		Name:        "test-api-route",
	}

	routeResp, err := service.RegisterRoute(ctx, routeReq)
	require.NoError(t, err)
	assert.NotEmpty(t, routeResp.RouteID)

	// 3. 添加中间件
	middlewareReq := &dto.AddMiddlewareRequest{
		ServerID:    serverResp.ServerID,
		Name:        "test-middleware",
		Type:        "general",
		Priority:    100,
		HandlerName: "test-middleware",
	}

	middlewareResp, err := service.AddMiddleware(ctx, middlewareReq)
	require.NoError(t, err)
	assert.NotEmpty(t, middlewareResp.MiddlewareID)

	// 4. 启动服务器
	startReq := &dto.StartServerRequest{
		ServerID: serverResp.ServerID,
	}

	startResp, err := service.StartServer(ctx, startReq)
	require.NoError(t, err)
	assert.Equal(t, "running", startResp.State)

	// 5. 获取服务器信息
	getReq := &dto.GetServerRequest{
		ServerID: serverResp.ServerID,
	}

	getResp, err := service.GetServer(ctx, getReq)
	require.NoError(t, err)
	assert.Equal(t, "running", getResp.State)
	assert.Equal(t, 1, getResp.RouteCount)
	assert.Equal(t, 1, getResp.MiddlewareCount)

	// 6. 查找匹配路由
	routeInfo, params, err := service.FindMatchingRoute(ctx, serverResp.ServerID, "GET", "/api/test")
	require.NoError(t, err)
	assert.Equal(t, "GET", routeInfo.Method)
	assert.Equal(t, "/api/test", routeInfo.Path)
	assert.Empty(t, params) // 没有路径参数

	// 7. 停止服务器
	stopReq := &dto.StopServerRequest{
		ServerID: serverResp.ServerID,
	}

	stopResp, err := service.StopServer(ctx, stopReq)
	require.NoError(t, err)
	assert.Equal(t, "stopped", stopResp.State)
}
