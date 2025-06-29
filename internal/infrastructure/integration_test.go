package infrastructure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/justinwongcn/ant/internal/application/dto"
	"github.com/justinwongcn/ant/internal/application/handlers"
	"github.com/justinwongcn/ant/internal/application/services"
	"github.com/justinwongcn/ant/internal/domain/shared/events"
	"github.com/justinwongcn/ant/internal/domain/webserver/aggregates"
	infraEvents "github.com/justinwongcn/ant/internal/infrastructure/events"
	infraHandlers "github.com/justinwongcn/ant/internal/infrastructure/handlers"
	"github.com/justinwongcn/ant/internal/infrastructure/repositories/memory"
)

// TestInfrastructureIntegration 测试基础设施层的集成
func TestInfrastructureIntegration(t *testing.T) {
	// 创建基础设施组件
	repoManager := memory.NewRepositoryManager()
	eventBus := infraEvents.NewMemoryEventBus()
	handlerRegistry := infraHandlers.NewHandlerRegistry()

	// 注册示例处理器
	routeHandler := infraHandlers.NewSimpleRouteHandler("hello-handler", "Hello, World!")
	err := handlerRegistry.RegisterRouteHandler("hello-handler", routeHandler)
	require.NoError(t, err)

	middlewareHandler := infraHandlers.NewSimpleMiddlewareHandler("logging-middleware")
	err = handlerRegistry.RegisterMiddlewareHandler("logging-middleware", middlewareHandler)
	require.NoError(t, err)

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

	appService := services.NewWebServerApplicationService(commandHandler, queryHandler)

	ctx := context.Background()

	// 测试创建服务器
	t.Run("创建Web服务器", func(t *testing.T) {
		createReq := &dto.CreateWebServerRequest{
			Name:    "test-server",
			Address: ":8080",
		}

		response, err := appService.CreateServer(ctx, createReq)
		require.NoError(t, err)
		assert.Equal(t, "test-server", response.Name)
		assert.Equal(t, ":8080", response.Address)
		assert.Equal(t, "stopped", response.State)
	})

	// 测试获取服务器
	t.Run("获取服务器信息", func(t *testing.T) {
		// 先创建服务器
		createReq := &dto.CreateWebServerRequest{
			Name:    "test-server-2",
			Address: ":8081",
		}

		createResp, err := appService.CreateServer(ctx, createReq)
		require.NoError(t, err)

		// 获取服务器信息
		getReq := &dto.GetServerRequest{
			ServerID: createResp.ServerID,
		}

		response, err := appService.GetServer(ctx, getReq)
		require.NoError(t, err)
		assert.Equal(t, createResp.ServerID, response.ServerID)
		assert.Equal(t, "test-server-2", response.Name)
		assert.Equal(t, ":8081", response.Address)
		assert.Equal(t, "stopped", response.State)
	})

	// 测试注册路由
	t.Run("注册路由", func(t *testing.T) {
		// 先创建服务器
		createReq := &dto.CreateWebServerRequest{
			Name:    "test-server-3",
			Address: ":8082",
		}

		createResp, err := appService.CreateServer(ctx, createReq)
		require.NoError(t, err)

		// 注册路由
		routeReq := &dto.RegisterRouteRequest{
			ServerID:    createResp.ServerID,
			Method:      "GET",
			Path:        "/hello",
			HandlerName: "hello-handler",
			Name:        "hello-route",
			Description: "Hello route",
		}

		routeResp, err := appService.RegisterRoute(ctx, routeReq)
		require.NoError(t, err)
		assert.Equal(t, createResp.ServerID, routeResp.ServerID)
		assert.Equal(t, "GET", routeResp.Method)
		assert.Equal(t, "/hello", routeResp.Path)
		assert.Equal(t, "hello-route", routeResp.Name)
	})

	// 测试添加中间件
	t.Run("添加中间件", func(t *testing.T) {
		// 先创建服务器
		createReq := &dto.CreateWebServerRequest{
			Name:    "test-server-4",
			Address: ":8083",
		}

		createResp, err := appService.CreateServer(ctx, createReq)
		require.NoError(t, err)

		// 添加中间件
		middlewareReq := &dto.AddMiddlewareRequest{
			ServerID:    createResp.ServerID,
			Name:        "logging",
			Type:        "logging",
			Priority:    100,
			HandlerName: "logging-middleware",
			Description: "Logging middleware",
		}

		middlewareResp, err := appService.AddMiddleware(ctx, middlewareReq)
		require.NoError(t, err)
		assert.Equal(t, createResp.ServerID, middlewareResp.ServerID)
		assert.Equal(t, "logging", middlewareResp.Name)
		assert.Equal(t, "logging", middlewareResp.Type)
		assert.Equal(t, 100, middlewareResp.Priority)
	})

	// 测试启动和停止服务器
	t.Run("启动和停止服务器", func(t *testing.T) {
		// 先创建服务器
		createReq := &dto.CreateWebServerRequest{
			Name:    "test-server-5",
			Address: ":8084",
		}

		createResp, err := appService.CreateServer(ctx, createReq)
		require.NoError(t, err)

		// 启动服务器
		startReq := &dto.StartServerRequest{
			ServerID: createResp.ServerID,
		}

		startResp, err := appService.StartServer(ctx, startReq)
		require.NoError(t, err)
		assert.Equal(t, createResp.ServerID, startResp.ServerID)
		assert.Equal(t, "running", startResp.State)

		// 停止服务器
		stopReq := &dto.StopServerRequest{
			ServerID: createResp.ServerID,
		}

		stopResp, err := appService.StopServer(ctx, stopReq)
		require.NoError(t, err)
		assert.Equal(t, createResp.ServerID, stopResp.ServerID)
		assert.Equal(t, "stopped", stopResp.State)
	})

	// 测试查找匹配路由
	t.Run("查找匹配路由", func(t *testing.T) {
		// 先创建服务器
		createReq := &dto.CreateWebServerRequest{
			Name:    "test-server-6",
			Address: ":8085",
		}

		createResp, err := appService.CreateServer(ctx, createReq)
		require.NoError(t, err)

		// 注册路由
		routeReq := &dto.RegisterRouteRequest{
			ServerID:    createResp.ServerID,
			Method:      "GET",
			Path:        "/users/{id}",
			HandlerName: "hello-handler",
			Name:        "user-route",
		}

		_, err = appService.RegisterRoute(ctx, routeReq)
		require.NoError(t, err)

		// 查找匹配路由
		routeInfo, params, err := appService.FindMatchingRoute(ctx, createResp.ServerID, "GET", "/users/123")
		require.NoError(t, err)
		assert.Equal(t, "GET", routeInfo.Method)
		assert.Equal(t, "/users/{id}", routeInfo.Path)
		assert.Equal(t, "user-route", routeInfo.Name)
		assert.Equal(t, "123", params["id"])
	})
}

// TestEventHandling 测试事件处理
func TestEventHandling(t *testing.T) {
	eventBus := infraEvents.NewMemoryEventBus()

	// 创建事件处理器
	loggingHandler := infraEvents.NewLoggingEventHandler(
		events.ServerStartedEventType,
		events.RouteRegisteredEventType,
	)

	// 订阅事件
	err := eventBus.Subscribe(loggingHandler)
	require.NoError(t, err)

	// 验证处理器数量
	assert.Equal(t, 1, eventBus.GetHandlerCount(events.ServerStartedEventType))
	assert.Equal(t, 1, eventBus.GetHandlerCount(events.RouteRegisteredEventType))

	// 发布事件
	ctx := context.Background()
	serverStartedEvent := events.NewServerStartedEvent("test-server", ":8080")
	routeRegisteredEvent := events.NewRouteRegisteredEvent("test-server", "GET", "/hello")

	err = eventBus.Publish(ctx, serverStartedEvent, routeRegisteredEvent)
	require.NoError(t, err)

	// 取消订阅
	err = eventBus.Unsubscribe(loggingHandler)
	require.NoError(t, err)

	// 验证处理器已移除
	assert.Equal(t, 0, eventBus.GetHandlerCount(events.ServerStartedEventType))
	assert.Equal(t, 0, eventBus.GetHandlerCount(events.RouteRegisteredEventType))
}

// TestRepositoryOperations 测试仓储操作
func TestRepositoryOperations(t *testing.T) {
	repoManager := memory.NewRepositoryManager()
	ctx := context.Background()

	// 测试Web服务器仓储
	t.Run("Web服务器仓储操作", func(t *testing.T) {
		server, err := aggregates.NewWebServer("test-server", ":8080")
		require.NoError(t, err)

		// 保存服务器
		err = repoManager.WebServer().Save(ctx, server)
		require.NoError(t, err)

		// 查找服务器
		foundServer, err := repoManager.WebServer().FindByID(ctx, server.ID())
		require.NoError(t, err)
		assert.Equal(t, server.ID(), foundServer.ID())
		assert.Equal(t, server.Name(), foundServer.Name())

		// 检查存在性
		exists, err := repoManager.WebServer().Exists(ctx, server.ID())
		require.NoError(t, err)
		assert.True(t, exists)

		// 计数
		count, err := repoManager.WebServer().Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	// 测试事务
	t.Run("事务操作", func(t *testing.T) {
		tx, err := repoManager.BeginTransaction(ctx)
		require.NoError(t, err)

		// 在事务中创建服务器
		server, err := aggregates.NewWebServer("tx-server", ":8081")
		require.NoError(t, err)

		err = tx.WebServer().Save(ctx, server)
		require.NoError(t, err)

		// 提交事务
		err = tx.Commit()
		require.NoError(t, err)

		// 验证服务器已保存
		foundServer, err := repoManager.WebServer().FindByID(ctx, server.ID())
		require.NoError(t, err)
		assert.Equal(t, server.ID(), foundServer.ID())
	})
}
