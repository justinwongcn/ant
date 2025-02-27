package recovery

import (
	"github.com/justinwongcn/ant"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoveryMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		handler        ant.HandleFunc
		expectedStatus int
		expectedBody   string
		logCalled      bool
	}{
		{
			name: "正常请求不触发panic恢复",
			handler: func(ctx *ant.Context) {
				ctx.RespData = []byte("success")
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
			logCalled:      false,
		},
		{
			name: "处理panic并返回自定义错误",
			handler: func(ctx *ant.Context) {
				panic("test panic")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal Server Error",
			logCalled:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logCalled := false
			mb := NewMiddlewareBuilder()
			mb.LogFunc = func(ctx *ant.Context) {
				logCalled = true
			}

			middleware := mb.Build()
			handler := middleware(tt.handler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			resp := httptest.NewRecorder()
			ctx := &ant.Context{
				Req:  req,
				Resp: resp,
			}

			handler(ctx)

			if resp.Code != tt.expectedStatus {
				t.Errorf("期望状态码 %d, 实际获得 %d", tt.expectedStatus, resp.Code)
			}

			if string(ctx.RespData) != tt.expectedBody {
				t.Errorf("期望响应体 %s, 实际获得 %s", tt.expectedBody, string(ctx.RespData))
			}

			if logCalled != tt.logCalled {
				t.Errorf("期望日志函数调用状态 %v, 实际获得 %v", tt.logCalled, logCalled)
			}
		})
	}
}

func TestRecoveryMiddlewareCustomization(t *testing.T) {
	mb := NewMiddlewareBuilder()
	mb.StatusCode = http.StatusServiceUnavailable
	mb.ErrMsg = "Service Unavailable"

	var logCtx *ant.Context
	mb.LogFunc = func(ctx *ant.Context) {
		logCtx = ctx
	}

	middleware := mb.Build()
	handler := middleware(func(ctx *ant.Context) {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	ctx := &ant.Context{
		Req:  req,
		Resp: resp,
	}

	handler(ctx)

	if ctx.RespStatusCode != http.StatusServiceUnavailable {
		t.Errorf("期望状态码 %d, 实际获得 %d", http.StatusServiceUnavailable, ctx.RespStatusCode)
	}

	if string(ctx.RespData) != "Service Unavailable" {
		t.Errorf("期望响应体 %s, 实际获得 %s", "Service Unavailable", string(ctx.RespData))
	}

	if logCtx == nil {
		t.Error("日志函数未被调用")
	}
}