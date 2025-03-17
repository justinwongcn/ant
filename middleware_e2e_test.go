package ant

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestMiddlewareE2E 测试中间件的端到端功能
func TestMiddlewareE2E(t *testing.T) {
	// 创建HTTP服务器
	server := NewHTTPServer()

	// 记录中间件执行顺序
	order := make([]string, 0)

	// 创建测试中间件
	mw1 := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			order = append(order, "mw1 before")
			next(ctx)
			order = append(order, "mw1 after")
		}
	}

	mw2 := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			order = append(order, "mw2 before")
			next(ctx)
			order = append(order, "mw2 after")
		}
	}

	// 注册中间件
	server.Use(mw1, mw2)

	// 注册测试路由
	server.Handle("GET /test", func(ctx *Context) {
		order = append(order, "handler")
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("success")
	})

	// 测试正常请求处理
	t.Run("正常请求处理", func(t *testing.T) {
		// 重置执行顺序记录
		order = make([]string, 0)

		// 创建测试请求
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		// 处理请求
		server.ServeHTTP(rec, req)

		// 验证状态码
		if rec.Code != http.StatusOK {
			t.Errorf("期望状态码200，得到：%d", rec.Code)
		}

		// 验证响应内容
		if rec.Body.String() != "success" {
			t.Errorf("期望响应体'success'，得到：%s", rec.Body.String())
		}

		// 验证中间件执行顺序
		expectedOrder := []string{"mw1 before", "mw2 before", "handler", "mw2 after", "mw1 after"}
		if len(order) != len(expectedOrder) {
			t.Errorf("期望执行顺序长度%d，得到：%d", len(expectedOrder), len(order))
		}
		for i, v := range expectedOrder {
			if order[i] != v {
				t.Errorf("期望执行顺序第%d步为%s，得到：%s", i+1, v, order[i])
			}
		}
	})

	// 测试错误处理
	t.Run("错误处理", func(t *testing.T) {
		// 创建新的服务器实例用于错误处理测试
		errorServer := NewHTTPServer()

		// 创建一个会触发错误的中间件
		errorMw := func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				ctx.RespStatusCode = http.StatusInternalServerError
				ctx.RespData = []byte("middleware error")
				// 不调用next，直接返回错误
			}
		}

		// 注册错误中间件
		errorServer.Use(errorMw)

		// 注册测试路由
		errorServer.Handle("GET /error", func(ctx *Context) {
			// 这个处理器不应该被执行
			t.Error("handler should not be called")
		})

		// 创建测试请求
		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		rec := httptest.NewRecorder()

		// 处理请求
		errorServer.ServeHTTP(rec, req)

		// 验证状态码
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("期望状态码500，得到：%d", rec.Code)
		}

		// 验证错误消息
		if rec.Body.String() != "middleware error" {
			t.Errorf("期望错误消息'middleware error'，得到：%s", rec.Body.String())
		}
	})

	// 测试中间件链的执行顺序
	t.Run("中间件链执行顺序", func(t *testing.T) {
		// 创建新的服务器实例
		chainServer := NewHTTPServer()

		// 用于记录执行顺序的切片
		executionOrder := make([]string, 0)

		// 创建多个中间件
		middlewares := make([]Middleware, 3)
		for i := range 3 {
			index := i // 捕获索引值
			middlewares[i] = func(next HandleFunc) HandleFunc {
				return func(ctx *Context) {
					executionOrder = append(executionOrder, fmt.Sprintf("mw%d before", index+1))
					next(ctx)
					executionOrder = append(executionOrder, fmt.Sprintf("mw%d after", index+1))
				}
			}
		}

		// 注册中间件
		chainServer.Use(middlewares...)

		// 注册测试路由
		chainServer.Handle("GET /chain", func(ctx *Context) {
			executionOrder = append(executionOrder, "handler")
			ctx.RespStatusCode = http.StatusOK
			ctx.RespData = []byte("chain test")
		})

		// 创建测试请求
		req := httptest.NewRequest(http.MethodGet, "/chain", nil)
		rec := httptest.NewRecorder()

		// 处理请求
		chainServer.ServeHTTP(rec, req)

		// 验证状态码
		if rec.Code != http.StatusOK {
			t.Errorf("期望状态码200，得到：%d", rec.Code)
		}

		// 验证响应内容
		if rec.Body.String() != "chain test" {
			t.Errorf("期望响应体'chain test'，得到：%s", rec.Body.String())
		}

		// 验证中间件执行顺序
		expectedOrder := []string{
			"mw1 before", "mw2 before", "mw3 before",
			"handler",
			"mw3 after", "mw2 after", "mw1 after",
		}
		if len(executionOrder) != len(expectedOrder) {
			t.Errorf("期望执行顺序长度%d，得到：%d", len(expectedOrder), len(executionOrder))
		}
		for i, v := range expectedOrder {
			if executionOrder[i] != v {
				t.Errorf("期望执行顺序第%d步为%s，得到：%s", i+1, v, executionOrder[i])
			}
		}
	})

	// 测试中间件修改请求和响应
	t.Run("中间件修改请求和响应", func(t *testing.T) {
		// 创建新的服务器实例
		modifyServer := NewHTTPServer()

		// 创建修改请求的中间件
		requestModifier := func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				// 修改请求头
				ctx.Req.Header.Set("X-Modified-By", "middleware")
				next(ctx)
			}
		}

		// 创建修改响应的中间件
		responseModifier := func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				next(ctx)
				// 修改响应头
				ctx.Resp.Header().Set("X-Response-Modified", "true")
				// 修改响应体
				if strings.Contains(string(ctx.RespData), "original") {
					ctx.RespData = []byte("modified response")
				}
			}
		}

		// 注册中间件
		modifyServer.Use(requestModifier, responseModifier)

		// 注册测试路由
		modifyServer.Handle("GET /modify", func(ctx *Context) {
			// 验证请求是否被修改
			if ctx.Req.Header.Get("X-Modified-By") != "middleware" {
				t.Error("请求头未被正确修改")
			}
			ctx.RespStatusCode = http.StatusOK
			ctx.RespData = []byte("original response")
		})

		// 创建测试请求
		req := httptest.NewRequest(http.MethodGet, "/modify", nil)
		rec := httptest.NewRecorder()

		// 处理请求
		modifyServer.ServeHTTP(rec, req)

		// 验证状态码
		if rec.Code != http.StatusOK {
			t.Errorf("期望状态码200，得到：%d", rec.Code)
		}

		// 验证响应头
		if rec.Header().Get("X-Response-Modified") != "true" {
			t.Error("响应头未被正确修改")
		}

		// 验证响应体是否被修改
		if rec.Body.String() != "modified response" {
			t.Errorf("期望响应体被修改为'modified response'，得到：%s", rec.Body.String())
		}
	})
}
