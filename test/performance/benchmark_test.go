// Package performance 包含性能测试
package performance

import (
	"context"
	"fmt"
	"testing"

	"github.com/justinwongcn/ant/internal/application/dto"
	"github.com/justinwongcn/ant/internal/application/handlers"
	"github.com/justinwongcn/ant/internal/application/services"
	infraEvents "github.com/justinwongcn/ant/internal/infrastructure/events"
	infraHandlers "github.com/justinwongcn/ant/internal/infrastructure/handlers"
	"github.com/justinwongcn/ant/internal/infrastructure/repositories/memory"
)

// setupBenchmarkService 设置基准测试服务
func setupBenchmarkService() services.WebServerService {
	repoManager := memory.NewRepositoryManager()
	eventBus := infraEvents.NewMemoryEventBus()
	handlerRegistry := infraHandlers.NewHandlerRegistry()

	// 注册处理器
	handler := infraHandlers.NewSimpleRouteHandler("bench-handler", "Benchmark Response")
	handlerRegistry.RegisterRouteHandler("bench-handler", handler)

	middleware := infraHandlers.NewSimpleMiddlewareHandler("bench-middleware")
	handlerRegistry.RegisterMiddlewareHandler("bench-middleware", middleware)

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

	return services.NewWebServerApplicationService(commandHandler, queryHandler)
}

// BenchmarkCreateServer 基准测试创建服务器
func BenchmarkCreateServer(b *testing.B) {
	service := setupBenchmarkService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &dto.CreateWebServerRequest{
			Name:    fmt.Sprintf("bench-server-%d", i),
			Address: fmt.Sprintf(":80%02d", i%100),
		}

		_, err := service.CreateServer(ctx, req)
		if err != nil {
			b.Fatalf("创建服务器失败: %v", err)
		}
	}
}

// BenchmarkRegisterRoute 基准测试注册路由
func BenchmarkRegisterRoute(b *testing.B) {
	service := setupBenchmarkService()
	ctx := context.Background()

	// 预先创建服务器
	createReq := &dto.CreateWebServerRequest{
		Name:    "bench-server",
		Address: ":8080",
	}

	serverResp, err := service.CreateServer(ctx, createReq)
	if err != nil {
		b.Fatalf("创建服务器失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &dto.RegisterRouteRequest{
			ServerID:    serverResp.ServerID,
			Method:      "GET",
			Path:        fmt.Sprintf("/api/bench/%d", i),
			HandlerName: "bench-handler",
			Name:        fmt.Sprintf("bench-route-%d", i),
		}

		_, err := service.RegisterRoute(ctx, req)
		if err != nil {
			b.Fatalf("注册路由失败: %v", err)
		}
	}
}

// BenchmarkFindMatchingRoute 基准测试查找匹配路由
func BenchmarkFindMatchingRoute(b *testing.B) {
	service := setupBenchmarkService()
	ctx := context.Background()

	// 预先创建服务器和路由
	createReq := &dto.CreateWebServerRequest{
		Name:    "bench-server",
		Address: ":8080",
	}

	serverResp, err := service.CreateServer(ctx, createReq)
	if err != nil {
		b.Fatalf("创建服务器失败: %v", err)
	}

	// 注册多个路由
	routes := []string{
		"/",
		"/api/users",
		"/api/users/{id}",
		"/api/users/{id}/posts",
		"/api/users/{id}/posts/{postId}",
		"/api/products",
		"/api/products/{id}",
		"/api/orders",
		"/api/orders/{id}",
		"/api/categories/{category}/products",
	}

	for i, path := range routes {
		req := &dto.RegisterRouteRequest{
			ServerID:    serverResp.ServerID,
			Method:      "GET",
			Path:        path,
			HandlerName: "bench-handler",
			Name:        fmt.Sprintf("route-%d", i),
		}

		_, err := service.RegisterRoute(ctx, req)
		if err != nil {
			b.Fatalf("注册路由失败: %v", err)
		}
	}

	// 测试路径
	testPaths := []string{
		"/",
		"/api/users",
		"/api/users/123",
		"/api/users/456/posts",
		"/api/users/789/posts/abc",
		"/api/products",
		"/api/products/xyz",
		"/api/orders",
		"/api/orders/order123",
		"/api/categories/electronics/products",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path := testPaths[i%len(testPaths)]
		_, _, err := service.FindMatchingRoute(ctx, serverResp.ServerID, "GET", path)
		if err != nil {
			b.Fatalf("查找匹配路由失败: %v", err)
		}
	}
}

// BenchmarkGetServer 基准测试获取服务器
func BenchmarkGetServer(b *testing.B) {
	service := setupBenchmarkService()
	ctx := context.Background()

	// 预先创建服务器
	createReq := &dto.CreateWebServerRequest{
		Name:    "bench-server",
		Address: ":8080",
	}

	serverResp, err := service.CreateServer(ctx, createReq)
	if err != nil {
		b.Fatalf("创建服务器失败: %v", err)
	}

	getReq := &dto.GetServerRequest{
		ServerID: serverResp.ServerID,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetServer(ctx, getReq)
		if err != nil {
			b.Fatalf("获取服务器失败: %v", err)
		}
	}
}

// BenchmarkConcurrentOperations 基准测试并发操作
func BenchmarkConcurrentOperations(b *testing.B) {
	service := setupBenchmarkService()
	ctx := context.Background()

	// 预先创建服务器
	createReq := &dto.CreateWebServerRequest{
		Name:    "bench-server",
		Address: ":8080",
	}

	serverResp, err := service.CreateServer(ctx, createReq)
	if err != nil {
		b.Fatalf("创建服务器失败: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// 混合操作：注册路由和查找路由
			if i%2 == 0 {
				// 注册路由
				req := &dto.RegisterRouteRequest{
					ServerID:    serverResp.ServerID,
					Method:      "GET",
					Path:        fmt.Sprintf("/concurrent/%d", i),
					HandlerName: "bench-handler",
					Name:        fmt.Sprintf("concurrent-route-%d", i),
				}

				_, err := service.RegisterRoute(ctx, req)
				if err != nil {
					b.Fatalf("注册路由失败: %v", err)
				}
			} else {
				// 获取服务器信息
				getReq := &dto.GetServerRequest{
					ServerID: serverResp.ServerID,
				}

				_, err := service.GetServer(ctx, getReq)
				if err != nil {
					b.Fatalf("获取服务器失败: %v", err)
				}
			}
			i++
		}
	})
}

// BenchmarkMemoryUsage 基准测试内存使用
func BenchmarkMemoryUsage(b *testing.B) {
	service := setupBenchmarkService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 创建服务器
		createReq := &dto.CreateWebServerRequest{
			Name:    fmt.Sprintf("memory-test-server-%d", i),
			Address: fmt.Sprintf(":90%02d", i%100),
		}

		serverResp, err := service.CreateServer(ctx, createReq)
		if err != nil {
			b.Fatalf("创建服务器失败: %v", err)
		}

		// 注册多个路由
		for j := 0; j < 10; j++ {
			routeReq := &dto.RegisterRouteRequest{
				ServerID:    serverResp.ServerID,
				Method:      "GET",
				Path:        fmt.Sprintf("/memory/test/%d/%d", i, j),
				HandlerName: "bench-handler",
				Name:        fmt.Sprintf("memory-route-%d-%d", i, j),
			}

			_, err := service.RegisterRoute(ctx, routeReq)
			if err != nil {
				b.Fatalf("注册路由失败: %v", err)
			}
		}

		// 添加中间件
		middlewareReq := &dto.AddMiddlewareRequest{
			ServerID:    serverResp.ServerID,
			Name:        fmt.Sprintf("memory-middleware-%d", i),
			Type:        "general",
			Priority:    100,
			HandlerName: "bench-middleware",
		}

		_, err = service.AddMiddleware(ctx, middlewareReq)
		if err != nil {
			b.Fatalf("添加中间件失败: %v", err)
		}
	}
}
