package ant

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandleRegistration 测试路由注册和请求处理
func TestHandleRegistration(t *testing.T) {
	server := NewHTTPServer()
	var isCalled bool

	testHandler := func(ctx *Context) {
		isCalled = true
		ctx.Resp.WriteHeader(http.StatusOK)
	}

	// 注册测试路由
	server.Handle("GET /test", testHandler)

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	// 触发请求处理
	server.ServeHTTP(rec, req)

	// 验证处理函数被调用
	if !isCalled {
		t.Fatal("Handler was not triggered")
	}

	// 验证状态码
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

// TestNotFound 测试未注册路由返回404
func TestNotFound(t *testing.T) {
	server := NewHTTPServer()

	req := httptest.NewRequest("GET", "/unknown", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rec.Code)
	}
}

// TestContextPassing 验证Context正确传递请求对象
func TestContextPassing(t *testing.T) {
	server := NewHTTPServer()
	var capturedReq *http.Request

	testHandler := func(ctx *Context) {
		capturedReq = ctx.Req
		ctx.Resp.WriteHeader(http.StatusOK)
	}

	server.Handle("GET /ctx", testHandler)

	req := httptest.NewRequest("GET", "/ctx", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	// 验证请求对象一致性
	if capturedReq != req {
		t.Error("Context did not receive correct request object")
	}
}

// TestParametricRoute 测试参数化路由（需Go 1.22+）
func TestParametricRoute(t *testing.T) {
	server := NewHTTPServer()
	var capturedID string

	testHandler := func(ctx *Context) {
		capturedID = ctx.Req.PathValue("id")
		ctx.Resp.WriteHeader(http.StatusOK)
	}

	server.Handle("GET /users/{id}", testHandler)

	req := httptest.NewRequest("GET", "/users/123", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if capturedID != "123" {
		t.Errorf("Expected param id=123, got %q", capturedID)
	}
}

// TestInvalidAddress 测试无效监听地址
func TestInvalidAddress(t *testing.T) {
	server := NewHTTPServer()

	// 使用不可能成功的监听地址
	err := server.Run("invalid-address:999999")

	if err == nil {
		t.Error("Expected error for invalid address")
	}
}

// TestUseMiddleware 测试中间件注册
func TestUseMiddleware(t *testing.T) {
	server := NewHTTPServer()

	if len(server.middlewares) != 0 {
		t.Error("Expected empty middleware slice initially")
	}

	// 测试单个中间件注册
	mw1 := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
		}
	}
	server.Use(mw1)

	if len(server.middlewares) != 1 {
		t.Error("Expected one middleware after registration")
	}

	// 测试多个中间件注册
	mw2 := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
		}
	}
	mw3 := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
		}
	}
	server.Use(mw2, mw3)

	if len(server.middlewares) != 3 {
		t.Error("Expected three middlewares after multiple registration")
	}
}

// TestMiddlewareChain 测试中间件链的构建和执行顺序
func TestMiddlewareChain(t *testing.T) {
	server := NewHTTPServer()
	executionOrder := make([]int, 0)

	// 创建三个测试中间件，每个都会记录自己的执行顺序
	mw1 := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			executionOrder = append(executionOrder, 1)
			next(ctx)
			executionOrder = append(executionOrder, 6)
		}
	}

	mw2 := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			executionOrder = append(executionOrder, 2)
			next(ctx)
			executionOrder = append(executionOrder, 5)
		}
	}

	mw3 := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			executionOrder = append(executionOrder, 3)
			next(ctx)
			executionOrder = append(executionOrder, 4)
		}
	}

	server.Use(mw1, mw2, mw3)

	handler := func(ctx *Context) {}

	chainedHandler := server.buildMiddlewareChain(handler)

	chainedHandler(&Context{
		Req:  httptest.NewRequest("GET", "/test", nil),
		Resp: httptest.NewRecorder(),
	})

	// 验证执行顺序
	expectedOrder := []int{1, 2, 3, 4, 5, 6}
	if len(executionOrder) != len(expectedOrder) {
		t.Errorf("Expected %d middleware executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, v := range expectedOrder {
		if executionOrder[i] != v {
			t.Errorf("Expected execution order %v, got %v", expectedOrder, executionOrder)
			break
		}
	}
}

// TestUseWithNilMiddlewares 测试 middlewares 为 nil 时的中间件注册
func TestUseWithNilMiddlewares(t *testing.T) {
	server := &HTTPServer{}
	// 确保初始状态下 middlewares 为 nil
	if server.middlewares != nil {
		t.Error("Expected nil middlewares initially")
	}

	// 测试注册单个中间件
	mw := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
		}
	}
	server.Use(mw)

	// 验证中间件被正确初始化和注册
	if len(server.middlewares) != 1 {
		t.Error("Expected one middleware after registration")
	}

	// 测试追加注册中间件
	mw2 := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
		}
	}
	server.Use(mw2)

	// 验证中间件被正确追加
	if len(server.middlewares) != 2 {
		t.Error("Expected two middlewares after second registration")
	}
}

func TestFlashResp(t *testing.T) {
	// 创建一个HTTPServer实例
	server := &HTTPServer{}

	// 创建一个模拟的Context
	ctx := & Context{
		Resp: &httptest.ResponseRecorder{
			Body: bytes.NewBuffer(nil), // 创建一个空的bytes.Buffer用于捕获写入的数据
		},
		RespStatusCode: 200,
		RespData:       []byte("Hello, World!"),
	}

	// 调用flashResp方法
	server.flashResp(ctx)

	// 检查状态码是否正确设置
	if status := ctx.Resp.(*httptest.ResponseRecorder).Code; status != 200 {
		t.Errorf("handler returned wrong status code: got %v want %v", status, 200)
	}

	// 检查响应体是否正确
	expected := "Hello, World!"
	if recBody := ctx.Resp.(*httptest.ResponseRecorder).Body.String(); recBody != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recBody, expected)
	}
}