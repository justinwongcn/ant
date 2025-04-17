# Ant Web Framework

## 项目简介

Ant 是一个轻量级的 Go Web 框架，专注于提供简单、灵活且功能丰富的 Web 开发体验。框架采用模块化设计，支持中间件扩展，并提供了完整的测试覆盖。

## 核心特性

### 路由系统
- 支持标准的 HTTP 方法（GET、POST 等）
- 支持参数化路由和通配符（需 Go 1.22+）
  - 支持方法匹配（如 `GET /posts/{id}`）
  - 支持通配符匹配（如 `/files/{pathname...}`）
  - 支持精确匹配（如 `/posts/{$}`）
  - 智能的路由优先级规则
    - 最具体的模式优先匹配
    - 方法匹配优先于通用匹配
    - 字面量路径优先于通配符
- 灵活的路由处理器注册机制
- 自动处理 405 Method Not Allowed 响应

### 模板引擎
- 基于 Go 标准库 html/template
- 支持从文件、目录或嵌入式文件系统加载模板
- 支持条件渲染等高级特性

### 文件处理
- 文件上传：支持自定义文件名和存储路径
- 文件下载：支持安全的文件下载和类型检测
- 静态资源服务：支持缓存和资源压缩

### 会话管理
- 支持多种会话存储方式（内存存储等）
- Cookie 传播器：处理会话 ID 的存取
- 完整的会话生命周期管理

### 中间件
- 访问日志：记录请求详情
- 错误处理：统一的错误处理机制
- 恢复机制：防止服务器因 panic 而崩溃

## 项目结构

```
.
├── context.go          # 请求上下文定义
├── server.go           # HTTP 服务器核心实现
├── template.go         # 模板引擎实现
├── files.go            # 文件处理功能
├── middleware/         # 中间件实现
│   ├── accesslog/      # 访问日志中间件
│   ├── errhandle/      # 错误处理中间件
│   └── recovery/       # 恢复中间件
└── session/           # 会话管理
    ├── cookie/        # Cookie 传播器
    └── memory/        # 内存存储实现
```

## 测试工作流

本项目配置了自动化测试工作流，确保代码质量和稳定性：

### GitHub Actions

在每次推送到 `main`、`master` 或 `dev` 分支以及创建针对这些分支的拉取请求时，GitHub Actions 会自动运行测试。工作流配置文件位于 `.github/workflows/go-test.yml`。

工作流执行以下操作：
1. 检出代码
2. 设置 Go 环境
3. 安装依赖
4. 运行标准测试
5. 使用竞态检测器运行测试

### 本地 Git 钩子

项目还配置了 Git pre-commit 钩子，在本地提交前自动运行测试。这确保只有通过测试的代码才能被提交。

如果您是新克隆此仓库，可能需要手动启用 pre-commit 钩子：

```bash
chmod +x .git/hooks/pre-commit
```

## 快速开始

### 安装

```bash
go get github.com/your-username/ant
```

### 基本使用

```go
package main

import "github.com/your-username/ant"

func main() {
    server := ant.NewHTTPServer()
    
    // 注册路由
    server.Handle("GET /hello", func(ctx *ant.Context) {
        ctx.WriteString("Hello, Ant!")
    })
    
    // 使用参数化路由
    server.Handle("GET /posts/{id}", func(ctx *ant.Context) {
        id := ctx.Request.PathValue("id")
        ctx.WriteString(fmt.Sprintf("Post ID: %s", id))
    })
    
    // 使用通配符路由
    server.Handle("GET /files/{pathname...}", func(ctx *ant.Context) {
        pathname := ctx.Request.PathValue("pathname")
        ctx.WriteString(fmt.Sprintf("File path: %s", pathname))
    })
    
    // 使用中间件
    server.Use(accesslog.MiddlewareBuilder{}.Build())
    
    // 启动服务器
    server.Run(":8080")
}
```

## 开发流程

1. 编写代码和测试
2. 本地运行测试：`go test ./...`
3. 提交代码（pre-commit 钩子会自动运行测试）
4. 推送到远程仓库（GitHub Actions 会自动运行测试）

这样的工作流确保了代码质量，并在问题出现时尽早发现。

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。