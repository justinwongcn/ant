package recovery

import (
	"github.com/justinwongcn/ant"
)

// MiddlewareBuilder 用于构建panic恢复中间件
type MiddlewareBuilder struct {
	// StatusCode 发生panic时返回的HTTP状态码
	StatusCode int
	// ErrMsg 发生panic时返回的错误信息
	ErrMsg string
	// LogFunc 用于记录panic信息的日志函数
	LogFunc func(ctx *ant.Context)
}

// NewMiddlewareBuilder 创建一个新的MiddlewareBuilder实例
// 默认使用500状态码和通用错误消息
func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		StatusCode: 500,
		ErrMsg:     "Internal Server Error",
		LogFunc:    func(ctx *ant.Context) {},
	}
}

// Build 构建panic恢复中间件
// 该中间件会捕获处理器中的panic，设置自定义的响应状态码和错误信息，
// 并通过用户定义的日志函数记录错误信息
func (m *MiddlewareBuilder) Build() ant.Middleware {
	return func(next ant.HandleFunc) ant.HandleFunc {
		return func(ctx *ant.Context) {
			defer func() {
				if err := recover(); err != nil {
					// 设置响应状态码和错误信息
					ctx.RespStatusCode = m.StatusCode
					ctx.Resp.WriteHeader(m.StatusCode)
					ctx.RespData = []byte(m.ErrMsg)
					// 调用日志函数记录错误信息
					m.LogFunc(ctx)
				}
			}()
			next(ctx)
		}
	}
}
