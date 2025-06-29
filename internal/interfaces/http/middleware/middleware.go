// Package middleware 提供HTTP中间件的接口层实现
package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// LoggingMiddleware 日志中间件
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 创建响应写入器包装器来捕获状态码
		wrapper := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// 处理请求
		next.ServeHTTP(wrapper, r)

		// 记录请求日志
		duration := time.Since(start)
		log.Printf("[%s] %s %s - %d - %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			wrapper.statusCode,
			duration,
		)
	})
}

// RecoveryMiddleware 恢复中间件，捕获panic
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())

				// 返回500错误
				http.Error(w, "内部服务器错误", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware CORS中间件
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头部
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecurityMiddleware 安全中间件
func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置安全头部
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		next.ServeHTTP(w, r)
	})
}

// ContentTypeMiddleware 内容类型中间件
func ContentTypeMiddleware(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware 简单的速率限制中间件
func RateLimitMiddleware(requestsPerSecond int) func(http.Handler) http.Handler {
	// 简化实现，实际项目中应该使用更复杂的算法
	lastRequest := make(map[string]time.Time)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.RemoteAddr
			now := time.Now()

			if lastTime, exists := lastRequest[clientIP]; exists {
				if now.Sub(lastTime) < time.Second/time.Duration(requestsPerSecond) {
					http.Error(w, "请求过于频繁", http.StatusTooManyRequests)
					return
				}
			}

			lastRequest[clientIP] = now
			next.ServeHTTP(w, r)
		})
	}
}

// AuthMiddleware 认证中间件
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 简化实现，检查Authorization头部
		authHeader := r.Header.Get("Authorization")

		// 跳过健康检查和根路径的认证
		if r.URL.Path == "/health" || r.URL.Path == "/ready" || r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		// 简单的token验证（实际项目中应该使用JWT或其他安全方案）
		if authHeader == "" || authHeader != "Bearer valid-token" {
			http.Error(w, "未授权访问", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CompressionMiddleware 压缩中间件（简化版本）
func CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查客户端是否支持gzip
		if !containsGzip(r.Header.Get("Accept-Encoding")) {
			next.ServeHTTP(w, r)
			return
		}

		// 设置压缩头部
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		// 这里应该创建gzip写入器，简化实现
		next.ServeHTTP(w, r)
	})
}

// TimeoutMiddleware 超时中间件
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, timeout, "请求超时")
	}
}

// responseWriter 响应写入器包装器
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 写入状态码
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// containsGzip 检查是否包含gzip编码
func containsGzip(acceptEncoding string) bool {
	return len(acceptEncoding) > 0 && (acceptEncoding == "gzip" ||
		len(acceptEncoding) > 4 && acceptEncoding[:4] == "gzip")
}

// MiddlewareChain 中间件链
type MiddlewareChain struct {
	middlewares []func(http.Handler) http.Handler
}

// NewMiddlewareChain 创建新的中间件链
func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

// Use 添加中间件到链中
func (c *MiddlewareChain) Use(middleware func(http.Handler) http.Handler) *MiddlewareChain {
	c.middlewares = append(c.middlewares, middleware)
	return c
}

// Then 应用中间件链到处理器
func (c *MiddlewareChain) Then(handler http.Handler) http.Handler {
	// 从后往前应用中间件
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i](handler)
	}
	return handler
}

// DefaultMiddlewareChain 创建默认的中间件链
func DefaultMiddlewareChain() *MiddlewareChain {
	return NewMiddlewareChain().
		Use(RecoveryMiddleware).
		Use(LoggingMiddleware).
		Use(SecurityMiddleware).
		Use(CORSMiddleware).
		Use(TimeoutMiddleware(30 * time.Second))
}

// APIMiddlewareChain 创建API专用的中间件链
func APIMiddlewareChain() *MiddlewareChain {
	return NewMiddlewareChain().
		Use(RecoveryMiddleware).
		Use(LoggingMiddleware).
		Use(SecurityMiddleware).
		Use(CORSMiddleware).
		Use(ContentTypeMiddleware("application/json")).
		Use(RateLimitMiddleware(100)). // 每秒100个请求
		Use(TimeoutMiddleware(30 * time.Second))
}
