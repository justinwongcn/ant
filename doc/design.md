# Ant Web 框架设计文档

## 1. 项目概述

### 1.1 项目背景
Ant 是一个基于 Go 1.22+ 版本新特性开发的轻量级 Web 框架，旨在提供一个简单、高效且易于使用的 Web 开发解决方案。该框架借鉴了 Gin 框架的优秀设计理念，同时充分利用 Go 1.22 中新增的多路路由特性，提供更简洁的 API 和更高的性能。

### 1.2 项目目标
- 充分利用 Go 1.22 的新路由特性
- 提供简洁易用的 API
- 保持高性能和低内存占用
- 提供完善的中间件支持
- 确保良好的可扩展性

## 2. 核心功能设计

### 2.1 路由系统
基于 Go 1.22 的 http.ServeMux 实现路由功能，直接利用标准库提供的路由匹配能力：
```go
func main() {
    app := ant.Run()
    
    app.Handle("GET /users/{id}", func(c *ant.Context) {
        // 处理逻辑
    })
    
    app.Run(":8080")
}
```

#### 2.1.1 路由注册
支持 HTTP 方法和路径模式的路由注册：
```go
func main() {
    app := ant.Run()
    
    // 注册路由处理器
    app.Handle("GET /api/users/{id}", GetUser)
    app.Handle("POST /api/users", CreateUser)
    app.Handle("PUT /api/users/{id}", UpdateUser)
    
    // 注册中间件
    app.Use(Logger())
    app.Use(Recovery())
    
    app.Run(":8080")
}
```

路由特性：
- 基于 http.ServeMux：直接使用 Go 1.22 标准库的路由功能
- 支持路径参数：通过 {param} 语法定义路径参数
- 支持 HTTP 方法：自动处理不同 HTTP 方法的路由匹配
- 中间件支持：支持全局和路由级别的中间件

### 2.2 中间件系统
支持全局中间件和路由级别中间件：
```go
// 全局中间件
app.Use(Logger())

// 路由级别中间件
app.Handle("GET /admin/{path}", AdminHandler).Use(AuthMiddleware())
```

### 2.3 上下文管理
提供丰富的上下文功能：
- 请求参数解析
- 响应数据封装
- 错误处理
- 数据绑定和验证
- Protocol Buffers 消息处理

### 2.4 Protocol Buffers 支持
#### 2.4.1 消息序列化
支持 Protocol Buffers 消息的序列化和反序列化：
```go
func main() {
    app := ant.Run()
    
    app.Handle("POST /api/users", func(c *ant.Context) {
        var user proto.User
        if err := c.BindProto(&user); err != nil {
            c.Error(err)
            return
        }
        c.ProtoJSON(200, &user)
    })
}
```

#### 2.4.2 gRPC 服务集成
支持在同一应用中运行 HTTP 和 gRPC 服务：
```go
func main() {
    app := ant.Run()
    
    // 注册 gRPC 服务
    grpcServer := app.GRPCServer()
    proto.RegisterUserServiceServer(grpcServer, &UserService{})
    
    // 启动 HTTP 和 gRPC 服务
    app.RunWithGRPC(":8080", ":9090")
}
```

#### 2.4.3 服务代理
支持自动生成 HTTP 到 gRPC 的代理：
```go
func main() {
    app := ant.Run()
    
    // 自动生成 HTTP 代理
    userProxy := app.GRPCProxy(proto.RegisterUserServiceHandlerFromEndpoint)
    app.Handle("POST /api/users/{id}", userProxy)
}
```

### 2.4 错误处理
统一的错误处理机制：
- 支持自定义错误处理器
- 内置常用 HTTP 错误处理
- panic 恢复机制

## 3. 技术选型

### 3.1 核心依赖
- Go 1.22+：利用新版本的路由特性
- 标准库 net/http：底层 HTTP 服务支持
- encoding/json：JSON 序列化/反序列化
- google.golang.org/protobuf：Protocol Buffers 支持
- google.golang.org/grpc：gRPC 服务支持

### 3.2 可选依赖
- golang.org/x/crypto：用于加密相关功能
- go-playground/validator：数据验证

## 4. 领域驱动设计

### 4.1 领域模型

#### 4.1.1 聚合根
- Application：应用聚合根，管理整个应用的生命周期
- Router：路由聚合根，负责路由注册和管理
- Server：服务器聚合根，负责HTTP和gRPC服务器的管理

#### 4.1.2 实体
- Context：请求上下文实体，包含请求生命周期内的所有信息
- Middleware：中间件实体，定义中间件的行为和生命周期
- Handler：处理器实体，负责具体的请求处理逻辑

#### 4.1.3 值对象
- Route：路由值对象，包含路径模式和HTTP方法
- Response：响应值对象，封装HTTP响应数据
- Config：配置值对象，包含框架配置信息

### 4.2 领域服务
- RouterService：路由注册和匹配服务
- MiddlewareService：中间件链管理服务
- ContextService：上下文管理服务

### 4.3 应用服务
- HTTPService：HTTP服务管理
- GRPCService：gRPC服务管理
- ValidationService：数据验证服务

### 4.4 基础设施
- Logger：日志基础设施
- ErrorHandler：错误处理基础设施
- Renderer：响应渲染基础设施

### 4.2 目录结构
```
ant/
├── docs/           # 文档
├── examples/       # 示例代码
├── middleware/     # 中间件实现
├── internal/       # 内部包
│   ├── binding/    # 数据绑定
│   ├── render/     # 响应渲染
│   └── proto/      # Protocol Buffers 相关实现
├── proto/          # Protocol Buffers 定义文件
├── context.go      # 上下文实现
├── app.go          # 核心应用
├── middleware.go   # 中间件定义
├── grpc.go         # gRPC 服务支持
└── errors.go       # 错误处理
```

## 5. 性能考虑

### 5.1 路由性能
- 利用 Go 1.22 的 http.ServeMux 路由特性，确保高效的路由匹配
- 减少路由处理的中间层开销

### 5.2 内存优化
- 使用对象池复用 Context 对象
- 减少不必要的内存分配
- 优化字符串处理

## 6. 安全性设计

### 6.1 内置安全特性
- XSS 防护
- CSRF 防护
- 请求限流
- 参数验证

### 6.2 安全中间件
- 提供基本的安全中间件
- 支持自定义安全策略

## 7. 扩展性设计

### 7.1 插件系统
支持通过插件扩展框架功能：
- 自定义中间件
- 自定义渲染器
- 自定义验证器

### 7.2 接口设计
提供清晰的接口定义，便于用户扩展：
- Handler 接口
- Middleware 接口
- Renderer 接口

## 8. 测试驱动开发

### 8.1 测试策略

#### 8.1.1 单元测试
- 领域模型测试：确保每个领域对象的行为符合预期
- 领域服务测试：验证领域服务的业务逻辑
- 应用服务测试：测试应用层功能的正确性

#### 8.1.2 集成测试
- 聚合根交互测试：验证聚合根之间的协作
- 中间件链测试：确保中间件的正确执行顺序
- HTTP/gRPC服务测试：验证服务端点的功能

### 8.2 TDD工作流程

#### 8.2.1 测试先行
1. 编写失败的测试用例
2. 实现最小可工作的代码
3. 重构优化代码
4. 重复以上步骤

#### 8.2.2 测试覆盖率要求
- 领域模型：100%覆盖率
- 领域服务：95%以上覆盖率
- 应用服务：90%以上覆盖率
- 基础设施：85%以上覆盖率

#### 8.2.3 测试规范
- 测试命名规范：Test{方法名}_{场景描述}
- 使用表格驱动测试
- 每个测试只验证一个行为
- 使用mock对象隔离外部依赖

## 9. 文档规划

### 9.1 API 文档
- 详细的 API 使用说明
- 示例代码
- 最佳实践指南

### 9.2 开发文档
- 框架设计说明
- 贡献指南
- 版本更新日志

## 10. 版本规划

### 10.1 v1.0.0
- 基础路由功能
- 核心中间件
- 基本上下文功能

### 10.2 后续版本
- 更多中间件支持
- 性能优化
- 更多便捷功能
- 插件系统

## 11. 项目管理

### 11.1 版本控制
- 使用 Git 进行版本控制
- 遵循语义化版本规范

### 11.2 发布策略
- 定期发布新版本
- 及时修复安全问题
- 保持向后兼容性