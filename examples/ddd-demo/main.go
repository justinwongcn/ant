// Package main demonstrates the new DDD architecture of Ant Web Framework
package main

import (
	"fmt"
	"log"

	sessionvo "github.com/justinwongcn/ant/internal/domain/session/valueobjects"
	"github.com/justinwongcn/ant/internal/domain/webserver/aggregates"
	"github.com/justinwongcn/ant/internal/domain/webserver/entities"
	"github.com/justinwongcn/ant/internal/domain/webserver/valueobjects"
)

// SimpleHandler implements the HandlerFunc interface
type SimpleHandler struct {
	name     string
	response string
}

func NewSimpleHandler(name, response string) *SimpleHandler {
	return &SimpleHandler{
		name:     name,
		response: response,
	}
}

func (h *SimpleHandler) Handle(ctx entities.RequestContext) error {
	// In a real implementation, this would write to the response
	fmt.Printf("Handler '%s' processing request: %s %s\n",
		h.name, ctx.Method().Value(), ctx.Path())
	fmt.Printf("Response: %s\n", h.response)
	return nil
}

func (h *SimpleHandler) Name() string {
	return h.name
}

// MockRequestContext implements the RequestContext interface for demo
type MockRequestContext struct {
	method     *valueobjects.HTTPMethod
	path       string
	parameters map[string]string
}

func NewMockRequestContext(method *valueobjects.HTTPMethod, path string) *MockRequestContext {
	return &MockRequestContext{
		method:     method,
		path:       path,
		parameters: make(map[string]string),
	}
}

func (m *MockRequestContext) Method() *valueobjects.HTTPMethod            { return m.method }
func (m *MockRequestContext) Path() string                                { return m.path }
func (m *MockRequestContext) Parameters() map[string]string               { return m.parameters }
func (m *MockRequestContext) QueryParams() map[string]string              { return make(map[string]string) }
func (m *MockRequestContext) Headers() map[string]string                  { return make(map[string]string) }
func (m *MockRequestContext) Body() []byte                                { return nil }
func (m *MockRequestContext) SetStatusCode(code *valueobjects.StatusCode) {}
func (m *MockRequestContext) SetHeader(name, value string)                {}
func (m *MockRequestContext) SetBody(body []byte)                         {}
func (m *MockRequestContext) Write(data []byte) error                     { return nil }

func main() {
	fmt.Println("=== Ant Web Framework - DDD Architecture Demo ===\n")

	// 1. 创建Web服务器聚合
	fmt.Println("1. 创建Web服务器聚合")
	server, err := aggregates.NewWebServer("demo-server", ":8080")
	if err != nil {
		log.Fatal("创建服务器失败:", err)
	}
	fmt.Printf("服务器创建成功: ID=%s, Name=%s, Address=%s\n\n",
		server.ID(), server.Name(), server.Address())

	// 2. 创建HTTP方法值对象
	fmt.Println("2. 创建HTTP方法值对象")
	getMethod := valueobjects.HTTPMethodGET
	postMethod := valueobjects.HTTPMethodPOST
	fmt.Printf("GET方法: %s (安全: %t, 幂等: %t)\n",
		getMethod.Value(), getMethod.IsSafe(), getMethod.IsIdempotent())
	fmt.Printf("POST方法: %s (安全: %t, 幂等: %t)\n\n",
		postMethod.Value(), postMethod.IsSafe(), postMethod.IsIdempotent())

	// 3. 创建URL模式值对象
	fmt.Println("3. 创建URL模式值对象")
	helloPattern, err := valueobjects.NewURLPattern(getMethod, "/hello")
	if err != nil {
		log.Fatal("创建URL模式失败:", err)
	}

	userPattern, err := valueobjects.NewURLPattern(getMethod, "/users/{id}")
	if err != nil {
		log.Fatal("创建URL模式失败:", err)
	}

	fmt.Printf("Hello模式: %s (优先级: %d)\n", helloPattern.String(), helloPattern.Priority())
	fmt.Printf("User模式: %s (优先级: %d, 有参数: %t)\n\n",
		userPattern.String(), userPattern.Priority(), userPattern.HasParameters())

	// 4. 创建处理器
	fmt.Println("4. 创建处理器")
	helloHandler := NewSimpleHandler("hello-handler", "Hello, World!")
	userHandler := NewSimpleHandler("user-handler", "User details")

	// 5. 创建路由实体
	fmt.Println("5. 创建路由实体")
	helloRoute, err := entities.NewRoute(helloPattern, helloHandler)
	if err != nil {
		log.Fatal("创建路由失败:", err)
	}

	userRoute, err := entities.NewRoute(userPattern, userHandler)
	if err != nil {
		log.Fatal("创建路由失败:", err)
	}

	fmt.Printf("Hello路由: ID=%s, Pattern=%s\n", helloRoute.ID(), helloRoute.Pattern().String())
	fmt.Printf("User路由: ID=%s, Pattern=%s\n\n", userRoute.ID(), userRoute.Pattern().String())

	// 6. 注册路由到服务器
	fmt.Println("6. 注册路由到服务器")
	if err := server.RegisterRoute(helloRoute); err != nil {
		log.Fatal("注册路由失败:", err)
	}

	if err := server.RegisterRoute(userRoute); err != nil {
		log.Fatal("注册路由失败:", err)
	}

	fmt.Printf("路由注册成功，服务器现有 %d 个路由\n\n", len(server.GetRoutes()))

	// 7. 测试路由匹配
	fmt.Println("7. 测试路由匹配")
	testRequests := []struct {
		method string
		path   string
	}{
		{"GET", "/hello"},
		{"GET", "/users/123"},
		{"GET", "/notfound"},
	}

	for _, req := range testRequests {
		method := valueobjects.MustNewHTTPMethod(req.method)
		route, params, err := server.FindMatchingRoute(method, req.path)

		fmt.Printf("请求: %s %s\n", req.method, req.path)
		if err != nil {
			fmt.Printf("  结果: 未找到匹配的路由 (%v)\n", err)
		} else {
			fmt.Printf("  匹配路由: %s\n", route.Pattern().String())
			if len(params) > 0 {
				fmt.Printf("  参数: %v\n", params)
			}

			// 模拟处理请求
			ctx := NewMockRequestContext(method, req.path)
			if err := route.Handle(ctx); err != nil {
				fmt.Printf("  处理错误: %v\n", err)
			}
		}
		fmt.Println()
	}

	// 8. 演示会话ID值对象
	fmt.Println("8. 演示会话ID值对象")
	sessionID, err := sessionvo.GenerateSessionID()
	if err != nil {
		log.Fatal("生成会话ID失败:", err)
	}

	fmt.Printf("生成的会话ID: %s\n", sessionID.Value())
	fmt.Printf("会话ID长度: %d\n", sessionID.Length())
	fmt.Printf("会话ID掩码: %s\n", sessionID.Mask())
	fmt.Printf("会话ID有效: %t\n\n", sessionID.IsValid())

	// 9. 演示领域事件
	fmt.Println("9. 演示领域事件")
	if err := server.Start(); err != nil {
		log.Fatal("启动服务器失败:", err)
	}

	events := server.GetEvents()
	fmt.Printf("服务器启动后产生了 %d 个领域事件:\n", len(events))
	for i, event := range events {
		fmt.Printf("  事件 %d: %s (聚合: %s, 时间: %s)\n",
			i+1, event.EventType(), event.AggregateID(), event.OccurredAt().Format("15:04:05"))
	}

	fmt.Println("\n=== Demo 完成 ===")
	fmt.Println("这个演示展示了新的DDD架构的核心概念:")
	fmt.Println("- 值对象 (HTTPMethod, URLPattern, SessionID)")
	fmt.Println("- 实体 (Route)")
	fmt.Println("- 聚合根 (WebServer)")
	fmt.Println("- 领域事件 (ServerStarted, RouteRegistered)")
	fmt.Println("- 业务规则和不变量的封装")
}
