package ant

import (
	"fmt"
)

// Middleware 定义中间件类型
// 中间件函数接收下一个处理器，返回一个新的处理器
type Middleware func(next HandleFunc) HandleFunc

// Chain 构建中间件链，按照传入顺序依次执行中间件
func Chain(ms ...Middleware) Middleware {
	return func(next HandleFunc) HandleFunc {
		for i := len(ms) - 1; i >= 0; i-- {
			next = ms[i](next)
		}
		return next
	}
}

// Recovery 错误恢复中间件
// 捕获处理器中的 panic，避免服务器崩溃
func Recovery() Middleware {
	return func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			defer func() {
				if err := recover(); err != nil {
					// 记录错误信息
					fmt.Printf("panic recovered: %v\n", err)
					// 返回 500 错误响应
					ctx.RespJSON(500, map[string]string{
						"error": "Internal Server Error",
					})
				}
			}()
			next(ctx)
		}
	}
}
