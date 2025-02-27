package accesslog

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/justinwongcn/ant"
)

func TestAccessLogMiddleware(t *testing.T) {
	t.Run("默认日志记录", func(t *testing.T) {
		// 创建模拟请求
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		ctx := &ant.Context{
			Req:  req,
			Resp: rec,
		}

		// 创建带缓冲区的日志捕获
		var logOutput bytes.Buffer
		logFn := func(s string) {
			logOutput.WriteString(s)
		}

		// 构建中间件
		builder := NewBuilder().LogFunc(logFn)
		md := builder.Build()

		// 模拟处理链
		handler := md(func(ctx *ant.Context) {
			time.Sleep(10 * time.Millisecond) // 模拟处理耗时
			ctx.RespStatusCode = http.StatusOK
		})

		handler(ctx)

		// 验证日志输出
		if logOutput.Len() == 0 {
			t.Fatal("未生成访问日志")
		}

		var logEntry accessLog
		if err := json.Unmarshal(logOutput.Bytes(), &logEntry); err != nil {
			t.Fatalf("日志解析失败: %v", err)
		}

		// 验证字段准确性
		assertEqual(t, "GET", logEntry.HTTPMethod)
		assertEqual(t, "/test", logEntry.Path)
		assertEqual(t, req.Host, logEntry.Host)
		if logEntry.Duration < 10*time.Millisecond {
			t.Errorf("处理时间计算不准确，预期≥10ms，实际得到 %v", logEntry.Duration)
		}
	})

	t.Run("自定义日志处理", func(t *testing.T) {
		called := false
		customLogFn := func(s string) {
			called = true
		}

		// 构建自定义中间件
		builder := NewBuilder().LogFunc(customLogFn)
		md := builder.Build()

		// 模拟请求处理
		handler := md(func(ctx *ant.Context) {})
		handler(&ant.Context{
			Req:  httptest.NewRequest("POST", "/api", nil),
			Resp: httptest.NewRecorder(),
		})

		if !called {
			t.Error("自定义日志函数未被调用")
		}
	})

	t.Run("完整请求生命周期", func(t *testing.T) {
		var logContent string
		captureLog := func(s string) {
			logContent = s
		}

		// 初始化中间件链
		md := NewBuilder().LogFunc(captureLog).Build()

		// 创建测试上下文
		req := httptest.NewRequest("PUT", "/resource/1", nil)
		ctx := &ant.Context{
			Req:  req,
			Resp: httptest.NewRecorder(),
		}

		// 模拟业务处理
		start := time.Now()
		handler := md(func(ctx *ant.Context) {
			time.Sleep(5 * time.Millisecond)
			ctx.RespStatusCode = http.StatusAccepted
		})
		handler(ctx)

		// 验证日志内容
		var logEntry accessLog
		if err := json.Unmarshal([]byte(logContent), &logEntry); err != nil {
			t.Fatal("日志解析错误:", err)
		}

		// 断言关键字段
		assertEqual(t, "PUT", logEntry.HTTPMethod)
		assertEqual(t, "/resource/1", logEntry.Path)
		if logEntry.Duration < 5*time.Millisecond {
			t.Errorf("处理时间异常，预期≥5ms，实际得到 %v", logEntry.Duration)
		}
		if time.Since(start)-logEntry.Duration > 1*time.Millisecond {
			t.Error("时间戳记录不准确")
		}
	})
}

// 辅助函数用于简化断言
func assertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("预期 %v，实际得到 %v", expected, actual)
	}
}
