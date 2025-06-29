// Package main 演示完整的DDD架构HTTP服务器
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justinwongcn/ant/internal/application/dto"
	"github.com/justinwongcn/ant/internal/application/handlers"
	"github.com/justinwongcn/ant/internal/application/services"
	infraEvents "github.com/justinwongcn/ant/internal/infrastructure/events"
	infraHandlers "github.com/justinwongcn/ant/internal/infrastructure/handlers"
	"github.com/justinwongcn/ant/internal/infrastructure/repositories/memory"
	"github.com/justinwongcn/ant/internal/interfaces/http/middleware"
	"github.com/justinwongcn/ant/internal/interfaces/http/server"
)

func main() {
	log.Println("=== 启动Ant Web框架 - DDD架构HTTP服务器 ===")

	// 1. 创建基础设施层组件
	log.Println("1. 初始化基础设施层...")

	// 仓储管理器
	repoManager := memory.NewRepositoryManager()

	// 事件总线
	eventBus := infraEvents.NewMemoryEventBus()

	// 处理器注册表
	handlerRegistry := infraHandlers.NewHandlerRegistry()

	// 2. 注册示例处理器
	log.Println("2. 注册示例处理器...")

	// 注册路由处理器
	helloHandler := infraHandlers.NewSimpleRouteHandler("hello-handler", "Hello, World!")
	if err := handlerRegistry.RegisterRouteHandler("hello-handler", helloHandler); err != nil {
		log.Fatalf("注册hello处理器失败: %v", err)
	}

	userHandler := infraHandlers.NewSimpleRouteHandler("user-handler", "User API")
	if err := handlerRegistry.RegisterRouteHandler("user-handler", userHandler); err != nil {
		log.Fatalf("注册user处理器失败: %v", err)
	}

	// 注册中间件处理器
	loggingMiddleware := infraHandlers.NewSimpleMiddlewareHandler("logging-middleware")
	if err := handlerRegistry.RegisterMiddlewareHandler("logging-middleware", loggingMiddleware); err != nil {
		log.Fatalf("注册logging中间件失败: %v", err)
	}

	authMiddleware := infraHandlers.NewSimpleMiddlewareHandler("auth-middleware")
	if err := handlerRegistry.RegisterMiddlewareHandler("auth-middleware", authMiddleware); err != nil {
		log.Fatalf("注册auth中间件失败: %v", err)
	}

	// 3. 创建应用层组件
	log.Println("3. 初始化应用层...")

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

	// 4. 创建接口层组件
	log.Println("4. 初始化接口层...")

	// HTTP服务器配置
	config := &server.Config{
		Address:         ":8080",
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		IdleTimeout:     120 * time.Second,
		ShutdownTimeout: 30 * time.Second,
	}

	// 创建HTTP服务器
	httpServer := server.NewHTTPServer(config, appService)

	// 添加中间件
	middlewareChain := middleware.DefaultMiddlewareChain()
	httpServer.AddMiddleware(middlewareChain.Then)

	// 5. 创建示例数据
	log.Println("5. 创建示例数据...")

	ctx := context.Background()

	// 创建示例Web服务器
	createReq := &dto.CreateWebServerRequest{
		Name:    "demo-server",
		Address: ":9090",
	}

	serverResp, err := appService.CreateServer(ctx, createReq)
	if err != nil {
		log.Fatalf("创建示例服务器失败: %v", err)
	}
	log.Printf("创建示例服务器成功: %s", serverResp.ServerID)

	// 注册示例路由
	routeReq := &dto.RegisterRouteRequest{
		ServerID:    serverResp.ServerID,
		Method:      "GET",
		Path:        "/hello",
		HandlerName: "hello-handler",
		Name:        "hello-route",
		Description: "Hello world route",
	}

	routeResp, err := appService.RegisterRoute(ctx, routeReq)
	if err != nil {
		log.Fatalf("注册示例路由失败: %v", err)
	}
	log.Printf("注册示例路由成功: %s", routeResp.RouteID)

	// 添加示例中间件
	middlewareReq := &dto.AddMiddlewareRequest{
		ServerID:    serverResp.ServerID,
		Name:        "logging",
		Type:        "logging",
		Priority:    100,
		HandlerName: "logging-middleware",
		Description: "Request logging middleware",
	}

	middlewareResp, err := appService.AddMiddleware(ctx, middlewareReq)
	if err != nil {
		log.Fatalf("添加示例中间件失败: %v", err)
	}
	log.Printf("添加示例中间件成功: %s", middlewareResp.MiddlewareID)

	// 6. 启动HTTP服务器
	log.Println("6. 启动HTTP服务器...")

	if err := httpServer.Start(); err != nil {
		log.Fatalf("启动HTTP服务器失败: %v", err)
	}

	log.Printf("HTTP服务器已启动，监听地址: %s", httpServer.Address())
	log.Println("API端点:")
	log.Println("  GET  /                     - 根路径")
	log.Println("  GET  /health               - 健康检查")
	log.Println("  GET  /ready                - 就绪检查")
	log.Println("  GET  /api/servers          - 获取服务器列表")
	log.Println("  POST /api/servers          - 创建服务器")
	log.Printf("  GET  /api/servers/%s       - 获取服务器信息", serverResp.ServerID)
	log.Printf("  GET  /api/servers/%s/routes - 获取路由列表", serverResp.ServerID)
	log.Printf("  POST /api/servers/%s/routes - 注册路由", serverResp.ServerID)
	log.Printf("  POST /api/servers/%s/start  - 启动服务器", serverResp.ServerID)
	log.Printf("  POST /api/servers/%s/stop   - 停止服务器", serverResp.ServerID)

	// 7. 等待中断信号
	log.Println("7. 服务器运行中，按Ctrl+C停止...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("收到停止信号，正在关闭服务器...")

	// 8. 优雅关闭
	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Stop(shutdownCtx); err != nil {
		log.Printf("关闭HTTP服务器失败: %v", err)
	}

	log.Println("服务器已关闭")
	log.Println("=== Ant Web框架演示结束 ===")
}
