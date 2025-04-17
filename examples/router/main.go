package main

import (
	"fmt"

	"github.com/justinwongcn/ant"
)

func main() {
	server := ant.NewHTTPServer()

	// 1. 固定路径 - 返回 "Hello, world"
	server.Handle("GET /hello", func(ctx *ant.Context) {
		ctx.WriteString("Hello, world")
	})

	// 2. 通配符路径 - 返回 "Hello, {name}"或错误信息
	server.Handle("GET /hello/{name}", func(ctx *ant.Context) {
		val := ctx.PathValue("name")
		valStr, err := val.String()
		if err != nil {
			ctx.WriteString("Error: Invalid path parameter")
			return
		}
		ctx.WriteString(fmt.Sprintf("Hello, %s", valStr))
	})

	// 3. 查询参数 - 返回 "hello {name}"
	server.Handle("GET /hello-query", func(ctx *ant.Context) {
		name, _ := ctx.DefaultQueryValue("name", "world").String()
		ctx.WriteString(fmt.Sprintf("hello %s", name))
	})

	fmt.Println("Server is running at http://localhost:8080")
	server.Run(":8080")
}
