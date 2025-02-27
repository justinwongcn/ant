package errhandle

import (
	"github.com/justinwongcn/ant"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorHandleMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		handler        ant.HandleFunc
		registerCode   int
		registerResp   []byte
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "未注册的错误码使用原始响应",
			handler: func(ctx *ant.Context) {
				ctx.RespStatusCode = http.StatusNotFound
				ctx.RespData = []byte("original response")
			},
			registerCode:   http.StatusBadRequest,
			registerResp:   []byte("Bad Request"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "original response",
		},
		{
			name: "已注册的错误码使用预设响应",
			handler: func(ctx *ant.Context) {
				ctx.RespStatusCode = http.StatusBadRequest
				ctx.RespData = []byte("original response")
			},
			registerCode:   http.StatusBadRequest,
			registerResp:   []byte("Bad Request"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad Request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mb := NewMiddlewareBuilder()
			mb.RegisterError(tt.registerCode, tt.registerResp)

			middleware := mb.Build()
			handler := middleware(tt.handler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			resp := httptest.NewRecorder()
			ctx := &ant.Context{
				Req:  req,
				Resp: resp,
			}

			handler(ctx)

			if ctx.RespStatusCode != tt.expectedStatus {
				t.Errorf("期望状态码 %d, 实际获得 %d", tt.expectedStatus, ctx.RespStatusCode)
			}

			if string(ctx.RespData) != tt.expectedBody {
				t.Errorf("期望响应体 %s, 实际获得 %s", tt.expectedBody, string(ctx.RespData))
			}
		})
	}
}

func TestErrorHandleMiddlewareChaining(t *testing.T) {
	mb := NewMiddlewareBuilder()
	mb.RegisterError(http.StatusBadRequest, []byte("Bad Request"))
		mb.RegisterError(http.StatusNotFound, []byte("Not Found"))

	middleware := mb.Build()
	handler := middleware(func(ctx *ant.Context) {
		ctx.RespStatusCode = http.StatusNotFound
		ctx.RespData = []byte("original response")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	ctx := &ant.Context{
		Req:  req,
		Resp: resp,
	}

	handler(ctx)

	if ctx.RespStatusCode != http.StatusNotFound {
		t.Errorf("期望状态码 %d, 实际获得 %d", http.StatusNotFound, ctx.RespStatusCode)
	}

	if string(ctx.RespData) != "Not Found" {
		t.Errorf("期望响应体 %s, 实际获得 %s", "Not Found", string(ctx.RespData))
	}
}