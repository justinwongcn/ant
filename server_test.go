package ant

import (
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
