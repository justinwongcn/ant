package ant

import (
	"fmt"
	"log"
	"net/http"
)

type HandleFunc func(ctx *Context)

type Server interface {
	// Handle 注册路由。pattern 支持 Go 1.22 的新路由语法，例如：
	// GET /users/{id}
	// POST /users
	Handle(pattern string, handler HandleFunc)
	// Run 启动服务器
	// addr 是监听地址。如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8082"
	Run(addr string) error
}

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = (*HTTPServer)(nil)

type HTTPServer struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		mux:         http.NewServeMux(),
		middlewares: make([]Middleware, 0),
	}
}

// Use 注册中间件
func (s *HTTPServer) Use(mdls ...Middleware) {
	if s.middlewares == nil {
		s.middlewares = mdls
		return
	}
	s.middlewares = append(s.middlewares, mdls...)
}

// Handle 注册路由
func (s *HTTPServer) Handle(pattern string, handler HandleFunc) {
	s.mux.Handle(pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 创建一个空的参数映射
		ctx := &Context{
			Req:  r,
			Resp: w,
		}
		handler(ctx)
	}))
}

// buildMiddlewareChain 构建中间件处理链
func (s *HTTPServer) buildMiddlewareChain(handler HandleFunc) HandleFunc {
	root := handler
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		root = s.middlewares[i](root)
	}
	return root
}

// ServeHTTP HTTPServer 处理请求的入口
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Req:  r,
		Resp: w,
	}

	// 使用公共函数构建中间件链
	handler := s.buildMiddlewareChain(s.serve)
	handler(ctx)
}

// serve 是最终的请求处理函数
func (s *HTTPServer) serve(ctx *Context) {
	s.mux.ServeHTTP(ctx.Resp, ctx.Req)
}

// flashResp 将Context中缓存的响应数据写入到HTTP响应中
func (s *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode > 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}

	// 将响应数据写入响应体
	_, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		log.Printf("回写响应失败: %v", err)
	}
}

// Run 启动服务器
func (s *HTTPServer) Run(addr string) error {
	fmt.Printf("Server is running on %s\n", addr)
	return http.ListenAndServe(addr, s)
}
