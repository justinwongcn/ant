# Ant Web框架故障排除指南

本文档提供使用Ant Web框架时常见问题的解决方案和调试技巧。

## 目录

- [常见错误](#常见错误)
- [性能问题](#性能问题)
- [配置问题](#配置问题)
- [调试技巧](#调试技巧)
- [日志分析](#日志分析)
- [监控和诊断](#监控和诊断)

## 常见错误

### 路由相关问题

#### 问题：路由不匹配

**症状：**
```
404 Not Found - 路由未找到
```

**可能原因和解决方案：**

1. **路由模式错误**
```go
// ❌ 错误的路由模式
server.Handle("/users/{id}", handler) // 缺少HTTP方法

// ✅ 正确的路由模式
server.Handle("GET /users/{id}", handler)
```

2. **路由优先级问题**
```go
// ❌ 路由注册顺序错误
server.Handle("GET /users/{id}", getUserHandler)
server.Handle("GET /users/profile", getProfileHandler) // 永远不会匹配

// ✅ 正确的注册顺序
server.Handle("GET /users/profile", getProfileHandler) // 更具体的路由先注册
server.Handle("GET /users/{id}", getUserHandler)
```

3. **参数名称不匹配**
```go
// ❌ 参数名称不一致
server.Handle("GET /users/{userId}", func(ctx *ant.Context) {
    id := ctx.PathValue("id") // 应该是 "userId"
})

// ✅ 参数名称一致
server.Handle("GET /users/{userId}", func(ctx *ant.Context) {
    id := ctx.PathValue("userId")
})
```

#### 问题：路径参数解析失败

**症状：**
```go
val, err := ctx.PathValue("id").String()
// err: 参数不存在或格式错误
```

**解决方案：**
```go
// ✅ 添加参数验证
func getUserHandler(ctx *ant.Context) {
    idVal := ctx.PathValue("id")
    id, err := idVal.String()
    if err != nil {
        ctx.RespStatusCode = http.StatusBadRequest
        ctx.WriteString("无效的用户ID")
        return
    }
    
    // 进一步验证ID格式
    if id == "" {
        ctx.RespStatusCode = http.StatusBadRequest
        ctx.WriteString("用户ID不能为空")
        return
    }
    
    // 处理业务逻辑
}
```

### 中间件相关问题

#### 问题：中间件执行顺序错误

**症状：**
```
认证中间件在日志中间件之后执行，导致未认证请求也被记录
```

**解决方案：**
```go
// ✅ 正确的中间件顺序
server.Use(recovery.NewBuilder().Build())  // 1. 异常恢复（最外层）
server.Use(loggingMiddleware)              // 2. 日志记录
server.Use(authMiddleware)                 // 3. 认证检查
server.Use(validationMiddleware)           // 4. 输入验证
```

#### 问题：中间件中忘记调用next()

**症状：**
```
请求被中间件拦截，但没有继续处理
```

**解决方案：**
```go
// ❌ 忘记调用next()
func BadMiddleware(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        // 做一些处理
        log.Println("处理请求")
        // 忘记调用 next(ctx)
    }
}

// ✅ 正确调用next()
func GoodMiddleware(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        // 前置处理
        log.Println("请求开始")
        
        // 调用下一个处理器
        next(ctx)
        
        // 后置处理
        log.Println("请求结束")
    }
}
```

### 模板相关问题

#### 问题：模板渲染失败

**症状：**
```
template: no template "index.html" associated with template "template"
```

**解决方案：**
```go
// ✅ 检查模板加载和使用
func setupTemplates() *ant.GoTemplateEngine {
    engine := &ant.GoTemplateEngine{}
    
    // 确保模板文件存在
    if _, err := os.Stat("templates/index.html"); os.IsNotExist(err) {
        log.Fatal("模板文件不存在: templates/index.html")
    }
    
    // 加载模板
    if err := engine.LoadFromGlob("templates/*.html"); err != nil {
        log.Fatal("模板加载失败:", err)
    }
    
    return engine
}

// 使用模板时检查模板引擎
func pageHandler(ctx *ant.Context) {
    if ctx.TemplateEngine == nil {
        ctx.RespStatusCode = http.StatusInternalServerError
        ctx.WriteString("模板引擎未配置")
        return
    }
    
    err := ctx.RespTemplate("index.html", data)
    if err != nil {
        log.Printf("模板渲染失败: %v", err)
        ctx.RespStatusCode = http.StatusInternalServerError
        ctx.WriteString("页面渲染失败")
        return
    }
}
```

### 会话相关问题

#### 问题：会话数据丢失

**症状：**
```
用户登录后，刷新页面会话信息丢失
```

**可能原因和解决方案：**

1. **Cookie配置问题**
```go
// ❌ Cookie配置错误
propagator := &cookie.Propagator{
    CookieName: "session_id",
    CookieOption: func(c *http.Cookie) {
        c.Secure = true // 在HTTP环境下设置为true会导致Cookie不被发送
    },
}

// ✅ 正确的Cookie配置
propagator := &cookie.Propagator{
    CookieName: "session_id",
    CookieOption: func(c *http.Cookie) {
        c.HttpOnly = true
        c.Secure = false    // 开发环境设为false，生产环境设为true
        c.Path = "/"
        c.MaxAge = 1800
        c.SameSite = http.SameSiteLaxMode
    },
}
```

2. **会话存储过期**
```go
// ✅ 检查会话存储配置
store := &memory.Store{
    Expiration: 30 * time.Minute, // 确保过期时间合理
}

// 在处理器中检查会话状态
func protectedHandler(ctx *ant.Context) {
    sess, ok := ctx.UserValues["session"].(session.Session)
    if !ok {
        log.Println("会话不存在于上下文中")
        redirectToLogin(ctx)
        return
    }
    
    // 检查会话是否有效
    userID, err := sess.Get(ctx.Req.Context(), "user_id")
    if err != nil {
        log.Printf("获取会话数据失败: %v", err)
        redirectToLogin(ctx)
        return
    }
    
    if userID == nil {
        log.Println("会话中无用户数据")
        redirectToLogin(ctx)
        return
    }
}
```

## 性能问题

### 问题：响应时间过长

**诊断步骤：**

1. **添加性能监控中间件**
```go
func PerformanceMiddleware(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        start := time.Now()
        
        next(ctx)
        
        duration := time.Since(start)
        if duration > 1*time.Second {
            log.Printf("慢请求警告: %s %s 耗时 %v", 
                ctx.Req.Method, ctx.Req.URL.Path, duration)
        }
    }
}
```

2. **检查数据库查询**
```go
// ✅ 添加数据库查询监控
func (r *UserRepository) FindByID(ctx context.Context, id string) (*User, error) {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        if duration > 100*time.Millisecond {
            log.Printf("慢查询警告: FindByID 耗时 %v", duration)
        }
    }()
    
    // 数据库查询逻辑
    return r.db.QueryUser(ctx, id)
}
```

3. **检查外部API调用**
```go
// ✅ 设置超时和重试
func callExternalAPI(ctx context.Context, url string) (*Response, error) {
    client := &http.Client{
        Timeout: 5 * time.Second, // 设置超时
    }
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("外部API调用失败: %w", err)
    }
    defer resp.Body.Close()
    
    // 处理响应
    return parseResponse(resp)
}
```

### 问题：内存使用过高

**诊断和解决：**

1. **检查内存泄漏**
```go
// ✅ 使用pprof进行内存分析
import _ "net/http/pprof"

func main() {
    // 启动pprof服务器
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // 启动主服务器
    server := ant.NewHTTPServer()
    // ... 配置路由
    server.Run(":8080")
}

// 使用命令分析内存：
// go tool pprof http://localhost:6060/debug/pprof/heap
```

2. **优化大对象处理**
```go
// ❌ 一次性加载大文件到内存
func badFileHandler(ctx *ant.Context) {
    data, err := os.ReadFile("large-file.txt") // 可能导致内存问题
    if err != nil {
        ctx.WriteString("文件读取失败")
        return
    }
    ctx.RespData = data
}

// ✅ 流式处理大文件
func goodFileHandler(ctx *ant.Context) {
    file, err := os.Open("large-file.txt")
    if err != nil {
        ctx.WriteString("文件打开失败")
        return
    }
    defer file.Close()
    
    ctx.Resp.Header().Set("Content-Type", "text/plain")
    io.Copy(ctx.Resp, file) // 流式传输
}
```

## 配置问题

### 问题：环境配置错误

**解决方案：**

1. **使用配置管理**
```go
// ✅ 配置结构体
type Config struct {
    Server struct {
        Port         int           `yaml:"port" env:"SERVER_PORT" default:"8080"`
        ReadTimeout  time.Duration `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT" default:"30s"`
        WriteTimeout time.Duration `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT" default:"30s"`
    } `yaml:"server"`
    
    Database struct {
        DSN         string `yaml:"dsn" env:"DATABASE_DSN"`
        MaxOpenConn int    `yaml:"max_open_conn" env:"DATABASE_MAX_OPEN_CONN" default:"25"`
        MaxIdleConn int    `yaml:"max_idle_conn" env:"DATABASE_MAX_IDLE_CONN" default:"5"`
    } `yaml:"database"`
}

// 加载配置
func LoadConfig() (*Config, error) {
    var cfg Config
    
    // 从文件加载
    if data, err := os.ReadFile("config.yaml"); err == nil {
        yaml.Unmarshal(data, &cfg)
    }
    
    // 从环境变量覆盖
    if err := env.Parse(&cfg); err != nil {
        return nil, err
    }
    
    return &cfg, nil
}
```

2. **配置验证**
```go
// ✅ 验证配置
func (c *Config) Validate() error {
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return errors.New("无效的服务器端口")
    }
    
    if c.Database.DSN == "" {
        return errors.New("数据库DSN不能为空")
    }
    
    return nil
}
```

## 调试技巧

### 启用详细日志

```go
// ✅ 分级日志
package logger

import (
    "log/slog"
    "os"
)

func NewLogger(level string) *slog.Logger {
    var logLevel slog.Level
    switch level {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    default:
        logLevel = slog.LevelInfo
    }
    
    opts := &slog.HandlerOptions{
        Level: logLevel,
    }
    
    handler := slog.NewJSONHandler(os.Stdout, opts)
    return slog.New(handler)
}

// 在处理器中使用
func userHandler(ctx *ant.Context) {
    logger := ctx.UserValues["logger"].(*slog.Logger)
    
    logger.Debug("处理用户请求", 
        "method", ctx.Req.Method,
        "path", ctx.Req.URL.Path,
        "user_agent", ctx.Req.UserAgent())
    
    // 业务逻辑
}
```

### 请求追踪

```go
// ✅ 请求ID追踪
func RequestIDMiddleware(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        requestID := ctx.Req.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        // 设置响应头
        ctx.Resp.Header().Set("X-Request-ID", requestID)
        
        // 存储到上下文
        if ctx.UserValues == nil {
            ctx.UserValues = make(map[string]any)
        }
        ctx.UserValues["request_id"] = requestID
        
        next(ctx)
    }
}

// 在日志中包含请求ID
func LoggingMiddleware(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        requestID := ctx.UserValues["request_id"].(string)
        
        log.Printf("[%s] 请求开始: %s %s", 
            requestID, ctx.Req.Method, ctx.Req.URL.Path)
        
        start := time.Now()
        next(ctx)
        duration := time.Since(start)
        
        log.Printf("[%s] 请求结束: %d - %v", 
            requestID, ctx.RespStatusCode, duration)
    }
}
```

### 健康检查端点

```go
// ✅ 健康检查
func healthCheckHandler(ctx *ant.Context) {
    health := map[string]interface{}{
        "status":    "ok",
        "timestamp": time.Now().Unix(),
        "version":   "1.0.0",
        "checks": map[string]interface{}{
            "database": checkDatabase(),
            "cache":    checkCache(),
            "external": checkExternalServices(),
        },
    }
    
    // 如果任何检查失败，返回503
    allHealthy := true
    for _, check := range health["checks"].(map[string]interface{}) {
        if check != "ok" {
            allHealthy = false
            break
        }
    }
    
    if !allHealthy {
        ctx.RespStatusCode = http.StatusServiceUnavailable
    }
    
    ctx.RespJSON(ctx.RespStatusCode, health)
}

func checkDatabase() string {
    // 检查数据库连接
    if err := db.Ping(); err != nil {
        return "error: " + err.Error()
    }
    return "ok"
}
```

通过这些故障排除技巧，您应该能够快速定位和解决Ant Web框架应用中的常见问题。记住，良好的日志记录和监控是快速诊断问题的关键。
