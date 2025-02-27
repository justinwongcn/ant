package accesslog

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/justinwongcn/ant"
)

// 辅助函数：创建测试上下文
func createTestContext(method, path string) (*ant.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	ctx := &ant.Context{
		Req:  req,
		Resp: rec,
	}
	return ctx, rec
}

// 辅助函数：捕获标准日志输出
func captureStdLog() (func() string, func()) {
	originalOutput := log.Writer()
	r, w, _ := os.Pipe()
	log.SetOutput(w)

	// 返回获取日志内容的函数
	getLog := func() string {
		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		return buf.String()
	}

	// 返回清理函数
	cleanup := func() {
		log.SetOutput(originalOutput)
	}

	return getLog, cleanup
}

// 辅助函数：捕获自定义日志输出
func captureCustomLog() (func(string), *bytes.Buffer) {
	var logBuffer bytes.Buffer
	logFn := func(s string) {
		logBuffer.WriteString(s)
	}
	return logFn, &logBuffer
}

// 辅助函数用于简化断言
func assertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("预期 %v，实际得到 %v", expected, actual)
	}
}

// 辅助函数：验证日志条目字段
func verifyLogEntry(t *testing.T, logData []byte, expectedMethod, expectedPath string, minDuration time.Duration) {
	t.Helper()
	var logEntry accessLog
	if err := json.Unmarshal(logData, &logEntry); err != nil {
		t.Fatalf("日志解析失败: %v", err)
	}

	// 验证字段准确性
	assertEqual(t, expectedMethod, logEntry.HTTPMethod)
	assertEqual(t, expectedPath, logEntry.Path)
	if logEntry.Duration < minDuration {
		t.Errorf("处理时间计算不准确，预期≥%v，实际得到 %v", minDuration, logEntry.Duration)
	}
}

// TestAccessLog 综合测试访问日志中间件的各种功能
func TestAccessLog(t *testing.T) {
	t.Run("默认日志记录", func(t *testing.T) {
		// 创建测试上下文
		ctx, _ := createTestContext("GET", "/test")

		// 创建带缓冲区的日志捕获
		logFn, logBuffer := captureCustomLog()

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
		if logBuffer.Len() == 0 {
			t.Fatal("未生成访问日志")
		}

		// 验证日志内容
		verifyLogEntry(t, logBuffer.Bytes(), "GET", "/test", 10*time.Millisecond)
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
		ctx, _ := createTestContext("POST", "/api")
		handler := md(func(ctx *ant.Context) {})
		handler(ctx)

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
		ctx, _ := createTestContext("PUT", "/resource/1")

		// 模拟业务处理
		start := time.Now()
		handler := md(func(ctx *ant.Context) {
			time.Sleep(5 * time.Millisecond)
			ctx.RespStatusCode = http.StatusAccepted
		})
		handler(ctx)

		// 验证日志内容
		verifyLogEntry(t, []byte(logContent), "PUT", "/resource/1", 5*time.Millisecond)

		// 验证时间戳记录准确性
		var logEntry accessLog
		json.Unmarshal([]byte(logContent), &logEntry)
		if time.Since(start)-logEntry.Duration > 1*time.Millisecond {
			t.Error("时间戳记录不准确")
		}
	})

	t.Run("默认日志函数", func(t *testing.T) {
		// 捕获标准日志输出
		getLog, cleanup := captureStdLog()
		defer cleanup()

		// 创建一个测试日志内容
		testLog := "测试日志内容"

		// 使用默认构建器的默认日志函数
		builder := NewBuilder()
		builder.logFunc(testLog)

		// 检查日志输出是否包含测试内容
		logOutput := getLog()
		if !strings.Contains(logOutput, testLog) {
			t.Errorf("默认日志函数未正确记录日志，预期包含 %q，实际得到 %q", testLog, logOutput)
		}
	})

	t.Run("AccessLog工厂函数", func(t *testing.T) {
		// 测试AccessLog函数是否返回有效的中间件
		middleware := AccessLog()

		// 验证返回的是否是一个有效的中间件函数
		if middleware == nil {
			t.Fatal("AccessLog()返回了nil，预期返回一个有效的中间件函数")
		}

		// 创建测试上下文
		ctx, _ := createTestContext("GET", "/test-access-log")

		// 捕获标准日志输出
		getLog, cleanup := captureStdLog()
		defer cleanup()

		// 使用中间件包装一个处理函数
		handler := middleware(func(ctx *ant.Context) {
			// 模拟处理逻辑
			ctx.RespStatusCode = http.StatusOK
		})

		// 执行处理函数
		handler(ctx)

		// 获取日志输出
		logOutput := getLog()

		// 验证日志输出包含请求信息
		if !strings.Contains(logOutput, "/test-access-log") {
			t.Errorf("AccessLog中间件未正确记录请求路径，日志输出: %s", logOutput)
		}
		if !strings.Contains(logOutput, "GET") {
			t.Errorf("AccessLog中间件未正确记录HTTP方法，日志输出: %s", logOutput)
		}
	})
}
