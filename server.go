package ant

import (
	"fmt"
	"log"
	"net/http"
)

// HandleFunc 定义HTTP请求处理函数类型
// ctx: 封装了HTTP请求上下文的Context对象
// 该函数类型用于注册路由处理器和中间件
type HandleFunc func(ctx *Context)

// Server 定义HTTP服务器接口
// 提供路由注册和服务器启动的核心功能
type Server interface {
	http.Handler
	// Handle 注册路由。pattern 支持 Go 1.22 的新路由语法，例如：
	// GET /users/{id}
	// POST /users
	//
	// pattern: 路由模式，支持HTTP方法和路径参数
	// handler: 处理该路由的处理函数
	Handle(pattern string, handler HandleFunc)

	// Run 启动服务器
	// addr: 监听地址，格式为 "host:port"。如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8082"
	// 返回值: 服务器运行过程中发生的错误
	Run(addr string) error
}

// 确保 HTTPServer 实现了 Server 接口
var _ Server = (*HTTPServer)(nil)

// HTTPServer HTTP服务器的具体实现
type HTTPServer struct {
	mux            *http.ServeMux // 底层路由复用器
	middlewares    []Middleware   // 已注册的中间件列表
	TemplateEngine TemplateEngine // 模板引擎
}

// ServerOption 定义服务器配置选项函数类型
// server: 需要配置的HTTP服务器实例
type ServerOption func(server *HTTPServer)

// ServerWithTemplateEngine 创建设置模板引擎的配置选项
// engine: 要使用的模板引擎实例
// 返回值: 配置函数
func ServerWithTemplateEngine(engine TemplateEngine) ServerOption {
	return func(server *HTTPServer) {
		server.TemplateEngine = engine
	}
}

// NewHTTPServer 创建一个新的HTTP服务器实例
// opts: 可选的服务器配置选项
// 返回值: 初始化后的HTTPServer指针
// 注意：默认不包含任何中间件，需要通过Use方法注册
func NewHTTPServer(opts ...ServerOption) *HTTPServer {
	server := &HTTPServer{
		mux:         http.NewServeMux(),
		middlewares: make([]Middleware, 0),
	}
	// 应用所有配置选项
	for _, opt := range opts {
		opt(server)
	}
	return server
}

// Use 注册中间件
// mdls: 要注册的中间件列表，支持同时注册多个
// 注意：中间件的调用顺序与注册顺序相反
func (s *HTTPServer) Use(mdls ...Middleware) {
	if s.middlewares == nil {
		s.middlewares = mdls
		return
	}
	s.middlewares = append(s.middlewares, mdls...)
}

// Handle 注册路由处理函数
// pattern: 路由模式，支持Go 1.22新路由语法
// handler: 该路由的处理函数
// 注意：每个请求都会创建新的Context实例
func (s *HTTPServer) Handle(pattern string, handler HandleFunc) {
	s.mux.Handle(pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 创建请求上下文
		ctx := &Context{
			Req:            r,
			Resp:           w,
			TemplateEngine: s.TemplateEngine, // 将服务器的模板引擎传递给Context
		}
		// 构建并执行中间件链
		middlewareChain := s.buildMiddlewareChain(handler)
		middlewareChain(ctx)
	}))
}

// buildMiddlewareChain 使用迭代器模式构建中间件调用链
// handler: 最终地请求处理函数
// 返回值: 包含所有中间件的处理函数
// 注意：中间件的执行顺序与注册顺序相反
func (s *HTTPServer) buildMiddlewareChain(handler HandleFunc) HandleFunc {
	// 创建中间件切片副本
	middlewares := s.middlewares

	// 返回包含完整中间件链的闭包
	return func(ctx *Context) {
		// 定义递归函数来依次调用中间件
		var next HandleFunc = handler // 链的终点是实际的处理器
		for i := len(middlewares) - 1; i >= 0; i-- {
			middleware := middlewares[i]
			next = middleware(next)
		}
		// 启动中间件链
		next(ctx)
		// 在所有中间件执行完成后写入响应
		s.writeResponse(ctx)
	}
}

// ServeHTTP 实现http.Handler接口
// 作为HTTP服务器的请求处理入口
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// writeResponse 将Context中缓存的响应数据写入HTTP响应
// ctx: 请求上下文
// 注意：会自动处理状态码和响应体的写入
func (s *HTTPServer) writeResponse(ctx *Context) {
	if ctx.RespStatusCode > 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}

	// 写入响应体
	_, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		log.Printf("回写响应失败: %v", err)
	}
}

// Run 启动HTTP服务器
// addr: 服务器监听地址
// 返回值: 服务器运行过程中的错误
// 注意：这是一个阻塞调用，服务器会一直运行直到出错
func (s *HTTPServer) Run(addr string) error {
	fmt.Printf("Server is running on %s\n", addr)
	return http.ListenAndServe(addr, s)
}
