package main

import (
	"fmt"
	"testing"

	"github.com/justinwongcn/ant"
)

func TestAnt(t *testing.T) {
	// 创建一个新的 HTTP 服务器
	server := ant.NewHTTPServer()

	// 注册处理 /hello 路径的处理函数
	server.Handle("GET /hello", func(ctx *ant.Context) {
		// 获取查询参数中的名字
		nameVal := ctx.QueryValue("name")
		name, err := nameVal.String()
		if err != nil || name == "" {
			name = "World"
		}

		// 设置响应头
		ctx.Resp.Header().Set("Content-Type", "text/plain; charset=utf-8")

		// 返回问候语
		fmt.Fprintf(ctx.Resp, "Hello, %s!", name)
	})

	// 启动服务器
	fmt.Println("Server is starting on http://localhost:8080")
	fmt.Println("Try visiting: http://localhost:8080/hello?name=YourName")
	if err := server.Run(":8080"); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}