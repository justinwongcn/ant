# Ant Web框架用户指南

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![Architecture](https://img.shields.io/badge/Architecture-DDD-green.svg)](https://en.wikipedia.org/wiki/Domain-driven_design)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](../LICENSE)

## 目录

- [简介](#简介)
- [安装和环境要求](#安装和环境要求)
- [快速开始](#快速开始)
- [核心概念](#核心概念)
- [路由系统](#路由系统)
- [请求上下文](#请求上下文)
- [中间件系统](#中间件系统)
- [模板引擎](#模板引擎)
- [文件处理](#文件处理)
- [会话管理](#会话管理)

## 简介

Ant是一个基于领域驱动设计(DDD)架构的现代化Go Web框架，专为构建可维护、可扩展的Web应用而设计。

### 主要特性

- **🏗️ DDD架构**：严格的四层架构设计，支持CQRS模式和事件驱动架构
- **🚀 高性能**：基于Go 1.24+的最新特性，优化的路由匹配算法
- **🔧 强类型安全**：完整的类型定义和编译时检查
- **🛡️ 安全可靠**：内置安全中间件和错误处理机制
- **📦 轻量级**：核心功能精简，按需扩展
- **🔌 可扩展**：清晰的分层架构，易于定制和扩展

### 适用场景

- RESTful API服务
- Web应用后端
- 微服务架构
- 企业级应用开发
- 需要严格架构约束的项目

## 安装和环境要求

### 环境要求

- Go 1.24.0 或更高版本
- 支持Go modules的开发环境

### 安装方法

```bash
# 初始化Go模块（如果还没有）
go mod init your-project-name

# 安装Ant框架
go get github.com/justinwongcn/ant

# 安装依赖
go mod tidy
```

### 验证安装

创建一个简单的测试文件来验证安装：

```go
package main

import (
    "fmt"
    "github.com/justinwongcn/ant"
)

func main() {
    server := ant.NewHTTPServer()
    fmt.Println("Ant框架安装成功！")
    
    server.Handle("GET /test", func(ctx *ant.Context) {
        ctx.WriteString("安装验证成功！")
    })
    
    fmt.Println("测试服务器启动在 http://localhost:8080/test")
    server.Run(":8080")
}
```

## 快速开始

### Hello World示例

让我们从最简单的Hello World示例开始：

```go
package main

import (
    "github.com/justinwongcn/ant"
)

func main() {
    // 创建HTTP服务器实例
    server := ant.NewHTTPServer()
    
    // 注册路由处理器
    server.Handle("GET /hello", func(ctx *ant.Context) {
        ctx.WriteString("Hello, Ant Web Framework!")
    })
    
    // 启动服务器
    server.Run(":8080")
}
```

### 带参数的路由

```go
package main

import (
    "fmt"
    "github.com/justinwongcn/ant"
)

func main() {
    server := ant.NewHTTPServer()
    
    // 参数化路由
    server.Handle("GET /hello/{name}", func(ctx *ant.Context) {
        name := ctx.PathValue("name")
        nameStr, _ := name.String()
        ctx.WriteString(fmt.Sprintf("Hello, %s!", nameStr))
    })
    
    // 查询参数
    server.Handle("GET /greet", func(ctx *ant.Context) {
        name, _ := ctx.DefaultQueryValue("name", "World").String()
        ctx.WriteString(fmt.Sprintf("Greetings, %s!", name))
    })
    
    server.Run(":8080")
}
```

### 使用中间件

```go
package main

import (
    "log"
    "time"
    "github.com/justinwongcn/ant"
    "github.com/justinwongcn/ant/middleware/accesslog"
)

func main() {
    server := ant.NewHTTPServer()
    
    // 添加访问日志中间件
    server.Use(accesslog.NewBuilder().Build())
    
    // 添加自定义中间件
    server.Use(func(next ant.HandleFunc) ant.HandleFunc {
        return func(ctx *ant.Context) {
            start := time.Now()
            next(ctx) // 执行下一个处理器
            duration := time.Since(start)
            log.Printf("请求处理耗时: %v", duration)
        }
    })
    
    server.Handle("GET /api/users", func(ctx *ant.Context) {
        ctx.WriteString("用户列表")
    })
    
    server.Run(":8080")
}
```

## 核心概念

### DDD四层架构

Ant框架采用严格的DDD四层架构：

```
┌─────────────────────────────────────────────────────────────┐
│                        接口层 (Interfaces)                    │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   HTTP服务器     │  │   中间件        │  │   处理器        │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                        应用层 (Application)                   │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   应用服务       │  │   命令处理器     │  │   查询处理器     │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                        领域层 (Domain)                        │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   聚合根        │  │   实体          │  │   值对象        │ │
│  │   领域服务       │  │   领域事件       │  │   仓储接口       │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                      基础设施层 (Infrastructure)               │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   仓储实现       │  │   事件总线       │  │   外部服务       │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### CQRS模式

命令查询职责分离(CQRS)是Ant框架的核心设计模式：

- **命令(Commands)**：修改系统状态的操作
- **查询(Queries)**：读取系统状态的操作
- **分离关注点**：读写操作完全分离，提高系统可维护性

### 事件驱动架构

- **领域事件**：业务状态变化的通知
- **事件总线**：解耦的事件发布订阅机制
- **异步处理**：非阻塞的事件处理

## 路由系统

Ant框架的路由系统基于Go 1.22+的新路由语法，提供强大而灵活的路由匹配功能。

### 基本路由注册

```go
server := ant.NewHTTPServer()

// HTTP方法 + 路径
server.Handle("GET /users", getUsersHandler)
server.Handle("POST /users", createUserHandler)
server.Handle("PUT /users/{id}", updateUserHandler)
server.Handle("DELETE /users/{id}", deleteUserHandler)
```

### 参数化路由

```go
// 路径参数
server.Handle("GET /users/{id}", func(ctx *ant.Context) {
    id := ctx.PathValue("id")
    idStr, err := id.String()
    if err != nil {
        ctx.WriteString("无效的用户ID")
        return
    }
    ctx.WriteString(fmt.Sprintf("用户ID: %s", idStr))
})

// 多个参数
server.Handle("GET /users/{userId}/posts/{postId}", func(ctx *ant.Context) {
    userId, _ := ctx.PathValue("userId").String()
    postId, _ := ctx.PathValue("postId").String()
    ctx.WriteString(fmt.Sprintf("用户%s的文章%s", userId, postId))
})
```

### 通配符路由

```go
// 通配符匹配
server.Handle("GET /files/{pathname...}", func(ctx *ant.Context) {
    pathname, _ := ctx.PathValue("pathname").String()
    ctx.WriteString(fmt.Sprintf("文件路径: %s", pathname))
})

// 精确匹配
server.Handle("GET /posts/{$}", func(ctx *ant.Context) {
    ctx.WriteString("文章列表页面")
})
```

### 路由优先级

路由匹配遵循以下优先级规则：

1. **最具体的模式优先**：字面量路径优先于通配符
2. **方法匹配优先**：指定HTTP方法的路由优先于通用路由
3. **参数数量**：参数较少的路由优先匹配

```go
// 优先级示例
server.Handle("GET /api/users/profile", profileHandler)     // 最高优先级
server.Handle("GET /api/users/{id}", getUserHandler)        // 中等优先级
server.Handle("GET /api/{resource...}", genericHandler)     // 最低优先级
```

## 请求上下文

`Context`是Ant框架的核心类型，封装了HTTP请求和响应的所有信息。

### Context结构

```go
type Context struct {
    Req  *http.Request       // 原始HTTP请求对象
    Resp http.ResponseWriter // HTTP响应写入器

    // 缓存和状态
    cacheQueryValues url.Values // 缓存的查询参数
    RespStatusCode   int        // 响应状态码
    RespData         []byte     // 响应数据

    // 扩展功能
    TemplateEngine TemplateEngine    // 模板引擎
    UserValues     map[string]any    // 用户数据存储
}
```

### 获取请求数据

#### 路径参数

```go
server.Handle("GET /users/{id}", func(ctx *ant.Context) {
    // 获取路径参数
    id := ctx.PathValue("id")
    idStr, err := id.String()
    if err != nil {
        ctx.WriteString("无效的参数")
        return
    }
    ctx.WriteString(fmt.Sprintf("用户ID: %s", idStr))
})
```

#### 查询参数

```go
server.Handle("GET /search", func(ctx *ant.Context) {
    // 获取查询参数
    keyword := ctx.QueryValue("q")
    keywordStr, _ := keyword.String()

    // 带默认值的查询参数
    page, _ := ctx.DefaultQueryValue("page", "1").String()

    ctx.WriteString(fmt.Sprintf("搜索: %s, 页码: %s", keywordStr, page))
})
```

#### 表单数据

```go
server.Handle("POST /login", func(ctx *ant.Context) {
    // 解析表单数据
    err := ctx.Req.ParseForm()
    if err != nil {
        ctx.WriteString("表单解析失败")
        return
    }

    username := ctx.Req.FormValue("username")
    password := ctx.Req.FormValue("password")

    ctx.WriteString(fmt.Sprintf("用户名: %s", username))
})
```

#### JSON数据

```go
server.Handle("POST /api/users", func(ctx *ant.Context) {
    var user struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }

    // 读取请求体
    body, err := io.ReadAll(ctx.Req.Body)
    if err != nil {
        ctx.WriteString("读取请求体失败")
        return
    }

    // 解析JSON
    err = json.Unmarshal(body, &user)
    if err != nil {
        ctx.WriteString("JSON解析失败")
        return
    }

    ctx.WriteString(fmt.Sprintf("创建用户: %s (%s)", user.Name, user.Email))
})
```

### 设置响应数据

#### 文本响应

```go
server.Handle("GET /text", func(ctx *ant.Context) {
    ctx.WriteString("纯文本响应")
})
```

#### JSON响应

```go
server.Handle("GET /api/user", func(ctx *ant.Context) {
    user := map[string]interface{}{
        "id":   1,
        "name": "张三",
        "email": "zhangsan@example.com",
    }

    // 设置响应头
    ctx.Resp.Header().Set("Content-Type", "application/json")

    // 序列化JSON
    jsonData, _ := json.Marshal(user)
    ctx.RespData = jsonData
    ctx.RespStatusCode = http.StatusOK
})
```

#### 自定义状态码

```go
server.Handle("GET /not-found", func(ctx *ant.Context) {
    ctx.RespStatusCode = http.StatusNotFound
    ctx.WriteString("页面未找到")
})
```

## 中间件系统

中间件是Ant框架的重要特性，允许在请求处理链中插入自定义逻辑。

### 中间件概念

中间件是一个函数，它接收下一个处理器并返回一个新的处理器：

```go
type Middleware func(next HandleFunc) HandleFunc
```

### 使用内置中间件

#### 访问日志中间件

```go
import "github.com/justinwongcn/ant/middleware/accesslog"

server := ant.NewHTTPServer()

// 使用默认访问日志中间件
server.Use(accesslog.AccessLog())

// 或者使用构建器模式自定义
server.Use(accesslog.NewBuilder().
    LogFunc(func(log string) {
        fmt.Println("自定义日志:", log)
    }).
    Build())
```

#### 错误处理中间件

```go
import "github.com/justinwongcn/ant/middleware/errhandle"

server.Use(errhandle.NewBuilder().
    RegisterError(http.StatusNotFound, func(ctx *ant.Context) {
        ctx.WriteString("自定义404页面")
    }).
    Build())
```

#### 恢复中间件

```go
import "github.com/justinwongcn/ant/middleware/recovery"

server.Use(recovery.NewBuilder().
    StatusCode(http.StatusInternalServerError).
    ErrMsg("服务器内部错误").
    LogFunc(func(ctx context.Context) {
        fmt.Println("发生panic，已恢复")
    }).
    Build())
```

### 自定义中间件

#### 简单中间件

```go
// 计时中间件
func TimingMiddleware(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        start := time.Now()

        // 执行下一个处理器
        next(ctx)

        duration := time.Since(start)
        log.Printf("请求 %s %s 耗时: %v",
            ctx.Req.Method, ctx.Req.URL.Path, duration)
    }
}

// 使用中间件
server.Use(TimingMiddleware)
```

#### 认证中间件

```go
func AuthMiddleware(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        // 检查认证头
        authHeader := ctx.Req.Header.Get("Authorization")
        if authHeader == "" {
            ctx.RespStatusCode = http.StatusUnauthorized
            ctx.WriteString("需要认证")
            return
        }

        // 验证token（这里简化处理）
        if !strings.HasPrefix(authHeader, "Bearer ") {
            ctx.RespStatusCode = http.StatusUnauthorized
            ctx.WriteString("无效的认证格式")
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token != "valid-token" {
            ctx.RespStatusCode = http.StatusUnauthorized
            ctx.WriteString("无效的token")
            return
        }

        // 认证通过，继续处理
        next(ctx)
    }
}
```

#### 跨域中间件

```go
func CORSMiddleware(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        // 设置CORS头
        ctx.Resp.Header().Set("Access-Control-Allow-Origin", "*")
        ctx.Resp.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        ctx.Resp.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // 处理预检请求
        if ctx.Req.Method == "OPTIONS" {
            ctx.RespStatusCode = http.StatusOK
            return
        }

        next(ctx)
    }
}
```

### 中间件链

中间件按注册顺序执行，形成一个处理链：

```go
server := ant.NewHTTPServer()

// 中间件执行顺序：Recovery -> CORS -> Auth -> Timing -> 路由处理器
server.Use(recovery.NewBuilder().Build())  // 1. 最外层：异常恢复
server.Use(CORSMiddleware)                 // 2. 跨域处理
server.Use(AuthMiddleware)                 // 3. 认证检查
server.Use(TimingMiddleware)               // 4. 性能计时
server.Use(accesslog.AccessLog())          // 5. 访问日志

server.Handle("GET /api/data", func(ctx *ant.Context) {
    ctx.WriteString("受保护的数据")
})
```

## 模板引擎

Ant框架提供基于Go标准库`html/template`的模板引擎，支持多种模板加载方式。

### 模板引擎配置

```go
import "github.com/justinwongcn/ant"

// 创建模板引擎
templateEngine := &ant.GoTemplateEngine{}

// 从文件加载模板
err := templateEngine.LoadFromFiles("templates/index.html", "templates/layout.html")
if err != nil {
    log.Fatal("模板加载失败:", err)
}

// 从glob模式加载
err = templateEngine.LoadFromGlob("templates/*.html")
if err != nil {
    log.Fatal("模板加载失败:", err)
}

// 从嵌入式文件系统加载
//go:embed templates/*
var templateFS embed.FS

err = templateEngine.LoadFromFS(templateFS, "templates/*.html")
if err != nil {
    log.Fatal("模板加载失败:", err)
}
```

### 在Context中使用模板

```go
server := ant.NewHTTPServer()

// 设置模板引擎
server.Handle("GET /page", func(ctx *ant.Context) {
    // 设置模板引擎
    ctx.TemplateEngine = templateEngine

    // 准备模板数据
    data := struct {
        Title   string
        Content string
        User    string
    }{
        Title:   "欢迎页面",
        Content: "这是使用Ant框架渲染的页面",
        User:    "张三",
    }

    // 渲染模板
    err := ctx.RespTemplate("index.html", data)
    if err != nil {
        ctx.WriteString("模板渲染失败: " + err.Error())
    }
})
```

### 模板示例

创建模板文件 `templates/index.html`：

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <meta charset="UTF-8">
</head>
<body>
    <h1>{{.Title}}</h1>
    <p>{{.Content}}</p>
    {{if .User}}
        <p>欢迎，{{.User}}！</p>
    {{else}}
        <p>请先登录</p>
    {{end}}
</body>
</html>
```

### 模板继承

创建布局模板 `templates/layout.html`：

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}} - 我的网站</title>
    <meta charset="UTF-8">
</head>
<body>
    <header>
        <h1>我的网站</h1>
    </header>
    <main>
        {{template "content" .}}
    </main>
    <footer>
        <p>&copy; 2024 我的网站</p>
    </footer>
</body>
</html>
```

创建内容模板 `templates/home.html`：

```html
{{define "content"}}
<h2>{{.Title}}</h2>
<p>{{.Content}}</p>
{{end}}
```

使用模板继承：

```go
server.Handle("GET /home", func(ctx *ant.Context) {
    ctx.TemplateEngine = templateEngine

    data := struct {
        Title   string
        Content string
    }{
        Title:   "首页",
        Content: "欢迎来到我的网站首页",
    }

    err := ctx.RespTemplate("layout.html", data)
    if err != nil {
        ctx.WriteString("模板渲染失败: " + err.Error())
    }
})
```

## 文件处理

Ant框架提供完整的文件上传、下载和静态文件服务功能。

### 文件上传

#### 基本文件上传

```go
import "github.com/justinwongcn/ant"

server := ant.NewHTTPServer()

// 创建文件上传处理器
uploader := &ant.FileUploader{
    FileField: "file", // 表单字段名
    DstPathFunc: func(fh *multipart.FileHeader) string {
        // 自定义上传目录
        return "./uploads/"
    },
}

// 注册上传路由
server.Handle("POST /upload", uploader.Handle())
```

#### 自定义文件名

```go
uploader := &ant.FileUploader{
    FileField: "file",
    DstPathFunc: func(fh *multipart.FileHeader) string {
        return "./uploads/"
    },
    FileNameFunc: func(originalName string) string {
        // 生成唯一文件名
        ext := filepath.Ext(originalName)
        return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
    },
}
```

#### 文件上传表单

```html
<!DOCTYPE html>
<html>
<head>
    <title>文件上传</title>
</head>
<body>
    <form action="/upload" method="post" enctype="multipart/form-data">
        <input type="file" name="file" required>
        <button type="submit">上传文件</button>
    </form>
</body>
</html>
```

### 文件下载

#### 基本文件下载

```go
// 创建文件下载处理器
downloader := &ant.FileDownloader{
    Dir: "./uploads", // 下载根目录
}

// 注册下载路由
server.Handle("GET /download", downloader.Handle())
```

使用方式：访问 `http://localhost:8080/download?file=example.txt`

#### 自定义下载处理

```go
server.Handle("GET /files/{filename}", func(ctx *ant.Context) {
    filename, _ := ctx.PathValue("filename").String()

    // 构建文件路径
    filePath := filepath.Join("./uploads", filename)

    // 检查文件是否存在
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        ctx.RespStatusCode = http.StatusNotFound
        ctx.WriteString("文件不存在")
        return
    }

    // 读取文件
    data, err := os.ReadFile(filePath)
    if err != nil {
        ctx.RespStatusCode = http.StatusInternalServerError
        ctx.WriteString("读取文件失败")
        return
    }

    // 设置响应头
    ctx.Resp.Header().Set("Content-Type", "application/octet-stream")
    ctx.Resp.Header().Set("Content-Disposition",
        fmt.Sprintf("attachment; filename=\"%s\"", filename))

    ctx.RespData = data
    ctx.RespStatusCode = http.StatusOK
})
```

### 静态文件服务

#### 基本静态文件服务

```go
// 创建静态文件处理器
staticHandler := &ant.StaticResourceHandler{
    Dir:                     "./static",
    PathPrefix:              "/static/",
    ExtensionContentTypeMap: map[string]string{
        ".html": "text/html; charset=utf-8",
        ".css":  "text/css; charset=utf-8",
        ".js":   "application/javascript; charset=utf-8",
        ".png":  "image/png",
        ".jpg":  "image/jpeg",
        ".gif":  "image/gif",
    },
}

// 注册静态文件路由
server.Handle("GET /static/{path...}", staticHandler.Handle())
```

#### 带缓存的静态文件服务

```go
// 创建带缓存的静态文件处理器
cachedHandler := &ant.StaticResourceHandlerWithCache{
    StaticResourceHandler: ant.StaticResourceHandler{
        Dir:        "./static",
        PathPrefix: "/static/",
        ExtensionContentTypeMap: map[string]string{
            ".css": "text/css; charset=utf-8",
            ".js":  "application/javascript; charset=utf-8",
        },
    },
    MaxSize: 1000, // 缓存最大条目数
}

server.Handle("GET /static/{path...}", cachedHandler.Handle())
```

## 会话管理

Ant框架提供完整的会话管理功能，支持多种存储方式和传播机制。

### 会话配置

#### 基本会话管理器

```go
import (
    "github.com/justinwongcn/ant/session"
    "github.com/justinwongcn/ant/session/cookie"
    "github.com/justinwongcn/ant/session/memory"
)

// 创建会话管理器
sessionManager := &session.Manager{
    Store: &memory.Store{
        Expiration: 30 * time.Minute, // 会话过期时间
    },
    Propagator: &cookie.Propagator{
        CookieName: "session_id",
        CookieOption: func(c *http.Cookie) {
            c.HttpOnly = true
            c.Secure = false // 开发环境设为false，生产环境设为true
            c.Path = "/"
            c.MaxAge = 1800 // 30分钟
        },
    },
    SessCtxKey: "session",
}
```

#### 会话中间件

```go
// 创建会话中间件
func SessionMiddleware(manager *session.Manager) ant.Middleware {
    return func(next ant.HandleFunc) ant.HandleFunc {
        return func(ctx *ant.Context) {
            // 尝试获取现有会话
            sess, err := manager.GetSession(*ctx)
            if err != nil {
                // 如果获取失败，创建新会话
                sessionID := generateSessionID() // 自定义会话ID生成函数
                sess, err = manager.InitSession(*ctx, sessionID)
                if err != nil {
                    ctx.WriteString("会话初始化失败")
                    return
                }
            }

            // 将会话存储到上下文中
            if ctx.UserValues == nil {
                ctx.UserValues = make(map[string]any)
            }
            ctx.UserValues["session"] = sess

            next(ctx)
        }
    }
}

// 生成会话ID的示例函数
func generateSessionID() string {
    return fmt.Sprintf("sess_%d_%s",
        time.Now().UnixNano(),
        strings.ReplaceAll(uuid.New().String(), "-", ""))
}

// 使用会话中间件
server.Use(SessionMiddleware(sessionManager))
```

### 会话操作

#### 存储和获取会话数据

```go
server.Handle("POST /login", func(ctx *ant.Context) {
    username := ctx.Req.FormValue("username")
    password := ctx.Req.FormValue("password")

    // 验证用户凭据（这里简化处理）
    if username == "admin" && password == "password" {
        // 获取会话
        sess := ctx.UserValues["session"].(session.Session)

        // 存储用户信息到会话
        sess.Set(ctx.Req.Context(), "user_id", "123")
        sess.Set(ctx.Req.Context(), "username", username)
        sess.Set(ctx.Req.Context(), "login_time", time.Now())

        ctx.WriteString("登录成功")
    } else {
        ctx.RespStatusCode = http.StatusUnauthorized
        ctx.WriteString("用户名或密码错误")
    }
})

server.Handle("GET /profile", func(ctx *ant.Context) {
    // 获取会话
    sess := ctx.UserValues["session"].(session.Session)

    // 从会话中获取用户信息
    userID, err := sess.Get(ctx.Req.Context(), "user_id")
    if err != nil {
        ctx.RespStatusCode = http.StatusUnauthorized
        ctx.WriteString("请先登录")
        return
    }

    username, _ := sess.Get(ctx.Req.Context(), "username")
    loginTime, _ := sess.Get(ctx.Req.Context(), "login_time")

    response := fmt.Sprintf("用户ID: %v, 用户名: %v, 登录时间: %v",
        userID, username, loginTime)
    ctx.WriteString(response)
})
```

#### 会话注销

```go
server.Handle("POST /logout", func(ctx *ant.Context) {
    // 删除会话
    err := sessionManager.RemoveSession(*ctx)
    if err != nil {
        ctx.WriteString("注销失败: " + err.Error())
        return
    }

    ctx.WriteString("注销成功")
})
```

#### 会话刷新

```go
server.Handle("POST /refresh", func(ctx *ant.Context) {
    // 刷新会话过期时间
    sess, err := sessionManager.RefreshSession(*ctx)
    if err != nil {
        ctx.WriteString("会话刷新失败: " + err.Error())
        return
    }

    ctx.WriteString("会话已刷新，ID: " + sess.ID())
})
```

### 自定义会话存储

#### Redis存储示例

```go
// 自定义Redis存储（示例实现）
type RedisStore struct {
    client     *redis.Client
    expiration time.Duration
}

func (r *RedisStore) Generate(ctx context.Context, id string) (session.Session, error) {
    sess := &RedisSession{
        id:     id,
        client: r.client,
        expiration: r.expiration,
    }
    return sess, nil
}

func (r *RedisStore) Get(ctx context.Context, id string) (session.Session, error) {
    // 检查会话是否存在
    exists, err := r.client.Exists(ctx, "session:"+id).Result()
    if err != nil {
        return nil, err
    }
    if exists == 0 {
        return nil, errors.New("会话不存在")
    }

    sess := &RedisSession{
        id:     id,
        client: r.client,
        expiration: r.expiration,
    }
    return sess, nil
}

func (r *RedisStore) Remove(ctx context.Context, id string) error {
    return r.client.Del(ctx, "session:"+id).Err()
}

func (r *RedisStore) Refresh(ctx context.Context, id string) error {
    return r.client.Expire(ctx, "session:"+id, r.expiration).Err()
}
```

### 会话安全

#### 安全配置

```go
// 生产环境会话配置
sessionManager := &session.Manager{
    Store: &memory.Store{
        Expiration: 30 * time.Minute,
    },
    Propagator: &cookie.Propagator{
        CookieName: "session_id",
        CookieOption: func(c *http.Cookie) {
            c.HttpOnly = true    // 防止XSS攻击
            c.Secure = true      // 仅HTTPS传输
            c.SameSite = http.SameSiteStrictMode // 防止CSRF攻击
            c.Path = "/"
            c.MaxAge = 1800
        },
    },
    SessCtxKey: "session",
}
```

#### 会话验证中间件

```go
func RequireAuthMiddleware(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        sess, ok := ctx.UserValues["session"].(session.Session)
        if !ok {
            ctx.RespStatusCode = http.StatusUnauthorized
            ctx.WriteString("未找到会话")
            return
        }

        // 检查用户是否已登录
        userID, err := sess.Get(ctx.Req.Context(), "user_id")
        if err != nil || userID == nil {
            ctx.RespStatusCode = http.StatusUnauthorized
            ctx.WriteString("请先登录")
            return
        }

        next(ctx)
    }
}

// 保护需要认证的路由
server.Handle("GET /admin", RequireAuthMiddleware(func(ctx *ant.Context) {
    ctx.WriteString("管理员页面")
}))
```

## 总结

本用户指南涵盖了Ant Web框架的主要功能和使用方法：

1. **安装和快速开始**：从环境要求到第一个Hello World应用
2. **核心概念**：DDD架构、CQRS模式和事件驱动架构
3. **路由系统**：支持Go 1.22+新语法的强大路由功能
4. **请求上下文**：完整的请求和响应处理
5. **中间件系统**：灵活的中间件链和自定义中间件
6. **模板引擎**：基于Go标准库的模板渲染
7. **文件处理**：文件上传、下载和静态资源服务
8. **会话管理**：完整的会话生命周期管理

更多详细信息和高级用法，请参考：
- [API参考文档](API_REFERENCE.md)
- [示例集合](EXAMPLES.md)
- [最佳实践指南](BEST_PRACTICES.md)
- [故障排除指南](TROUBLESHOOTING.md)
```
```
