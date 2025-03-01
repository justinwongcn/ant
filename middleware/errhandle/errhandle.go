package errhandle

import "github.com/justinwongcn/ant"

// MiddlewareBuilder 用于构建错误处理中间件
type MiddlewareBuilder struct {
	resp map[int][]byte
}

// NewMiddlewareBuilder 创建一个新的MiddlewareBuilder实例
func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		resp: make(map[int][]byte, 64),
	}
}

// RegisterError 注册错误码对应的响应内容
// code: HTTP状态码
// resp: 对应的响应内容
func (m *MiddlewareBuilder) RegisterError(code int, resp []byte) *MiddlewareBuilder {
	m.resp[code] = resp
	return m
}

// Build 构建错误处理中间件
// 该中间件会检查响应状态码，如果匹配已注册的错误码，则使用预设的响应内容
func (m *MiddlewareBuilder) Build() ant.Middleware {
	return func(next ant.HandleFunc) ant.HandleFunc {
		return func(ctx *ant.Context) {
			// 先执行后续的处理函数
			next(ctx)
			
			// 检查状态码是否匹配预设的错误响应
			resp, ok := m.resp[ctx.RespStatusCode]
			if ok {
				// 使用预设的错误响应替换原有响应
				ctx.RespData = resp
			}
		}
	}
}