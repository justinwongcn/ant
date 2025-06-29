# Ant Web Framework - DDD架构版本

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![Architecture](https://img.shields.io/badge/Architecture-DDD-green.svg)](https://en.wikipedia.org/wiki/Domain-driven_design)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 项目简介

Ant 是一个基于领域驱动设计(DDD)架构的现代化Go Web框架，提供清晰的分层架构、强类型安全和高性能。框架采用严格的四层架构设计，支持CQRS模式和事件驱动架构。

## 🏗️ DDD架构概览

本项目采用严格的DDD四层架构：

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

## 🎯 DDD核心特性

### 领域驱动设计
- **聚合根**: Web服务器作为聚合根管理路由和中间件
- **实体**: 路由和中间件作为独立实体
- **值对象**: HTTP方法、URL模式等不可变值对象
- **领域事件**: 服务器启动、路由注册等事件
- **仓储模式**: 抽象数据访问层

### CQRS模式
- **命令**: 修改系统状态的操作
- **查询**: 读取系统状态的操作
- **分离关注点**: 读写操作完全分离

### 事件驱动架构
- **领域事件**: 业务状态变化的通知
- **事件总线**: 解耦的事件发布订阅机制
- **异步处理**: 非阻塞的事件处理

## 🚀 快速开始

### 安装依赖

```bash
go mod tidy
```

### 运行示例

```bash
# 运行HTTP服务器示例
cd examples/http-server
go run main.go
```

### 运行测试

```bash
# 运行所有测试
go test ./... -v

# 运行端到端测试
go test ./test/e2e -v

# 运行性能测试
go test ./test/performance -bench=. -benchmem
```

## 🔧 API示例

### 创建Web服务器

```go
// 创建服务器
createReq := &dto.CreateWebServerRequest{
    Name:    "my-server",
    Address: ":8080",
}

response, err := appService.CreateServer(ctx, createReq)
```

### 注册路由

```go
// 注册路由
routeReq := &dto.RegisterRouteRequest{
    ServerID:    serverID,
    Method:      "GET",
    Path:        "/api/users/{id}",
    HandlerName: "user-handler",
    Name:        "get-user",
}

response, err := appService.RegisterRoute(ctx, routeReq)
```

## 📊 性能指标

基于基准测试的性能数据：

- **创建服务器**: ~1,405 ns/op, 682 B/op, 12 allocs/op
- **注册路由**: ~146,682 ns/op, 113,873 B/op, 1,404 allocs/op
- **查找路由**: 高效的路由匹配算法
- **并发安全**: 支持高并发读写操作

## 传统特性

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

## 📁 DDD项目结构

```
ant/
├── internal/                          # 内部代码
│   ├── domain/                        # 领域层
│   │   ├── shared/                    # 共享领域概念
│   │   │   ├── events/               # 领域事件
│   │   │   └── errors/               # 领域错误
│   │   └── webserver/                # Web服务器领域
│   │       ├── aggregates/           # 聚合根
│   │       ├── entities/             # 实体
│   │       ├── valueobjects/         # 值对象
│   │       ├── repositories/         # 仓储接口
│   │       └── services/             # 领域服务
│   ├── application/                   # 应用层
│   │   ├── commands/                 # 命令对象
│   │   ├── queries/                  # 查询对象
│   │   ├── dto/                      # 数据传输对象
│   │   ├── handlers/                 # 命令/查询处理器
│   │   └── services/                 # 应用服务
│   ├── infrastructure/               # 基础设施层
│   │   ├── repositories/             # 仓储实现
│   │   │   └── memory/              # 内存仓储
│   │   ├── events/                   # 事件总线实现
│   │   └── handlers/                 # 处理器注册表
│   └── interfaces/                   # 接口层
│       └── http/                     # HTTP接口
│           ├── server/               # HTTP服务器
│           ├── handlers/             # HTTP处理器
│           └── middleware/           # HTTP中间件
├── examples/                         # 示例代码
│   └── http-server/                  # HTTP服务器示例
├── test/                            # 测试代码
│   ├── e2e/                         # 端到端测试
│   └── performance/                 # 性能测试
└── docs/                           # 文档

## 传统项目结构（已重构为DDD）

```
.
├── context.go          # 请求上下文定义（已迁移到接口层）
├── server.go           # HTTP 服务器核心实现（已迁移到接口层）
├── template.go         # 模板引擎实现（已迁移到基础设施层）
├── files.go            # 文件处理功能（已迁移到基础设施层）
├── middleware/         # 中间件实现（已迁移到接口层）
│   ├── accesslog/      # 访问日志中间件
│   ├── errhandle/      # 错误处理中间件
│   └── recovery/       # 恢复中间件
└── session/           # 会话管理（已迁移到基础设施层）
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