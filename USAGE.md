# Ant 框架使用文档

`ant` 是一个基于 Go 语言的、采用领域驱动设计思想的现代化 Web 框架。它旨在帮助开发者构建结构清晰、可维护、可扩展的 Web 应用。

## 核心特性

*   **轻量级核心**：提供构建 Web 应用所需的核心组件，包括路由、中间件、上下文管理等。
*   **强大的路由**：支持 RESTful 风格的路由，包括静态路由和参数路由。
*   **中间件生态**：支持全局和局部的中间件，内置访问日志、错误处理、异常恢复等常用中间件。
*   **会话管理**：内置基于 Cookie 和内存的会话管理机制。
*   **领域驱动设计 (DDD)**：项目结构遵循 DDD 思想，鼓励开发者编写高内聚、低耦合的代码。
*   **可扩展性**：清晰的分层架构使得替换或增加自定义实现变得容易。

## 安装

```bash
go get github.com/justinwongcn/ant
```
*注意：由于我无法得知确切的 git 仓库地址，请在上方命令中替换为真实的项目路径。*

## 快速入门

下面是一个最简单的 "Hello, World" 示例，展示了如何启动一个 `ant` 服务器。

```go
package main

import (
	"github.com/justinwongcn/ant" // 替换为实际的 import 路径
	"net/http"
)

func main() {
	// 创建一个新的 Ant 服务器实例
	server := ant.NewServer()

	// 注册一个 GET 路由
	// 当用户访问 "/" 时，执行回调函数
	server.GET("/", func(c *ant.Context) {
		c.String(http.StatusOK, "Hello, Ant!")
	})

	// 启动服务器，监听 8080 端口
	// 如果启动失败，将打印错误信息
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
```

## 路由

`ant` 提供了简单直观的方式来定义路由。

### 基本路由

你可以为所有标准的 HTTP 方法注册路由。

```go
server.GET("/ping", func(c *ant.Context) {
    c.String(http.StatusOK, "pong")
})

server.POST("/users", func(c *ant.Context) {
    // 创建用户的逻辑
    c.String(http.StatusCreated, "User created")
})
```

### 路由参数

通过 `:name` 的形式可以定义路由参数，在处理器中通过 `c.Param("name")` 来获取。

```go
server.GET("/users/:id", func(c *ant.Context) {
	// 从 URL 路径中获取 "id" 参数
	id := c.Param("id")
	c.String(http.StatusOK, "Hello, user %s", id)
})
```

## 中间件

中间件可用于在请求处理链中插入逻辑，例如日志、认证、错误处理等。

### 使用内置中间件

可以使用 `server.Use()` 来为所有路由注册一个或多个全局中间件。

```go
package main

import (
	"github.com/justinwongcn/ant"
	"github.com/justinwongcn/ant/middleware/accesslog"
	"net/http"
)

func main() {
	server := ant.NewServer()

	// 使用访问日志中间件
	server.Use(accesslog.NewBuilder().Build())

	server.GET("/", func(c *ant.Context) {
		c.String(http.StatusOK, "This request will be logged.")
	})

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
```
