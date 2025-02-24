package ant

import (
	"fmt"
	"net/http"
)

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
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
var _ Server = &HTTPServer{}

type HTTPServer struct {
	mux *http.ServeMux
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		mux: http.NewServeMux(),
	}
}

// Handle 注册路由
func (s *HTTPServer) Handle(pattern string, handler HandleFunc) {
	s.mux.Handle(pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 创建一个空的参数映射
		params := make(map[string]string)
		ctx := &Context{
			Req:    r,
			Resp:   w,
			Params: func(key string) string {
				return params[key]
			},
		}
		handler(ctx)
	}))
}

// ServeHTTP HTTPServer 处理请求的入口
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Run 启动服务器
func (s *HTTPServer) Run(addr string) error {
	fmt.Printf("Server is running on %s\n", addr)
	return http.ListenAndServe(addr, s)
}
