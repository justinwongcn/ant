# Ant Web框架示例集合

本文档提供各种使用场景的完整示例，从基础到高级，帮助您快速掌握Ant框架的使用。

## 目录

- [基础示例](#基础示例)
- [路由示例](#路由示例)
- [中间件示例](#中间件示例)
- [模板示例](#模板示例)
- [文件处理示例](#文件处理示例)
- [会话管理示例](#会话管理示例)
- [DDD架构示例](#ddd架构示例)
- [完整应用示例](#完整应用示例)

## 基础示例

### Hello World

最简单的Web服务器示例：

```go
package main

import (
    "github.com/justinwongcn/ant"
)

func main() {
    server := ant.NewHTTPServer()
    
    server.Handle("GET /", func(ctx *ant.Context) {
        ctx.WriteString("Hello, Ant Web Framework!")
    })
    
    server.Run(":8080")
}
```

### 多路由示例

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/justinwongcn/ant"
)

func main() {
    server := ant.NewHTTPServer()
    
    // 首页
    server.Handle("GET /", func(ctx *ant.Context) {
        ctx.WriteString("欢迎来到Ant框架!")
    })
    
    // 关于页面
    server.Handle("GET /about", func(ctx *ant.Context) {
        ctx.WriteString("关于我们")
    })
    
    // API端点
    server.Handle("GET /api/status", func(ctx *ant.Context) {
        ctx.RespJSON(http.StatusOK, map[string]string{
            "status": "ok",
            "message": "服务器运行正常",
        })
    })
    
    fmt.Println("服务器启动在 http://localhost:8080")
    server.Run(":8080")
}
```

## 路由示例

### 参数化路由

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/justinwongcn/ant"
)

func main() {
    server := ant.NewHTTPServer()
    
    // 用户详情页面
    server.Handle("GET /users/{id}", func(ctx *ant.Context) {
        id, err := ctx.PathValue("id").String()
        if err != nil {
            ctx.RespStatusCode = http.StatusBadRequest
            ctx.WriteString("无效的用户ID")
            return
        }
        
        ctx.WriteString(fmt.Sprintf("用户ID: %s", id))
    })
    
    // 用户的文章
    server.Handle("GET /users/{userId}/posts/{postId}", func(ctx *ant.Context) {
        userId, _ := ctx.PathValue("userId").String()
        postId, _ := ctx.PathValue("postId").String()
        
        response := fmt.Sprintf("用户 %s 的文章 %s", userId, postId)
        ctx.WriteString(response)
    })
    
    // 数字ID验证
    server.Handle("GET /products/{id}", func(ctx *ant.Context) {
        id, err := ctx.PathValue("id").ToInt64()
        if err != nil {
            ctx.RespStatusCode = http.StatusBadRequest
            ctx.WriteString("产品ID必须是数字")
            return
        }
        
        ctx.WriteString(fmt.Sprintf("产品ID: %d", id))
    })
    
    server.Run(":8080")
}
```

### 通配符路由

```go
package main

import (
    "fmt"
    "strings"
    "github.com/justinwongcn/ant"
)

func main() {
    server := ant.NewHTTPServer()
    
    // 文件路径处理
    server.Handle("GET /files/{path...}", func(ctx *ant.Context) {
        path, _ := ctx.PathValue("path").String()
        
        // 安全检查：防止目录遍历
        if strings.Contains(path, "..") {
            ctx.RespStatusCode = http.StatusBadRequest
            ctx.WriteString("非法的文件路径")
            return
        }
        
        ctx.WriteString(fmt.Sprintf("文件路径: %s", path))
    })
    
    // API版本处理
    server.Handle("GET /api/{version}/{resource...}", func(ctx *ant.Context) {
        version, _ := ctx.PathValue("version").String()
        resource, _ := ctx.PathValue("resource").String()
        
        response := fmt.Sprintf("API版本: %s, 资源: %s", version, resource)
        ctx.WriteString(response)
    })
    
    server.Run(":8080")
}
```

### 查询参数处理

```go
package main

import (
    "fmt"
    "strconv"
    "github.com/justinwongcn/ant"
)

func main() {
    server := ant.NewHTTPServer()
    
    // 搜索功能
    server.Handle("GET /search", func(ctx *ant.Context) {
        // 必需参数
        query, err := ctx.QueryValue("q").String()
        if err != nil || query == "" {
            ctx.RespStatusCode = http.StatusBadRequest
            ctx.WriteString("缺少搜索关键词")
            return
        }
        
        // 可选参数，带默认值
        pageStr, _ := ctx.DefaultQueryValue("page", "1").String()
        page, err := strconv.Atoi(pageStr)
        if err != nil || page < 1 {
            page = 1
        }
        
        limitStr, _ := ctx.DefaultQueryValue("limit", "10").String()
        limit, err := strconv.Atoi(limitStr)
        if err != nil || limit < 1 || limit > 100 {
            limit = 10
        }
        
        response := fmt.Sprintf("搜索: %s, 页码: %d, 每页: %d", query, page, limit)
        ctx.WriteString(response)
    })
    
    // 过滤器示例
    server.Handle("GET /products", func(ctx *ant.Context) {
        category, _ := ctx.DefaultQueryValue("category", "all").String()
        minPriceStr, _ := ctx.DefaultQueryValue("min_price", "0").String()
        maxPriceStr, _ := ctx.DefaultQueryValue("max_price", "999999").String()
        
        minPrice, _ := strconv.ParseFloat(minPriceStr, 64)
        maxPrice, _ := strconv.ParseFloat(maxPriceStr, 64)
        
        response := fmt.Sprintf("分类: %s, 价格范围: %.2f - %.2f", 
            category, minPrice, maxPrice)
        ctx.WriteString(response)
    })
    
    server.Run(":8080")
}
```

## 中间件示例

### 自定义日志中间件

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/justinwongcn/ant"
)

// 详细日志中间件
func DetailedLoggingMiddleware(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        start := time.Now()
        
        // 记录请求开始
        log.Printf("[开始] %s %s - %s", 
            ctx.Req.Method, 
            ctx.Req.URL.Path, 
            ctx.Req.RemoteAddr)
        
        // 执行下一个处理器
        next(ctx)
        
        // 记录请求结束
        duration := time.Since(start)
        log.Printf("[结束] %s %s - %d - %v", 
            ctx.Req.Method, 
            ctx.Req.URL.Path, 
            ctx.RespStatusCode, 
            duration)
    }
}

func main() {
    server := ant.NewHTTPServer()
    
    // 使用自定义日志中间件
    server.Use(DetailedLoggingMiddleware)
    
    server.Handle("GET /", func(ctx *ant.Context) {
        // 模拟一些处理时间
        time.Sleep(100 * time.Millisecond)
        ctx.WriteString("Hello with logging!")
    })
    
    server.Run(":8080")
}
```

### 认证中间件

```go
package main

import (
    "net/http"
    "strings"
    "github.com/justinwongcn/ant"
)

// JWT认证中间件（简化版）
func JWTAuthMiddleware(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        // 跳过公开路径
        publicPaths := []string{"/", "/login", "/register", "/health"}
        for _, path := range publicPaths {
            if ctx.Req.URL.Path == path {
                next(ctx)
                return
            }
        }
        
        // 检查Authorization头
        authHeader := ctx.Req.Header.Get("Authorization")
        if authHeader == "" {
            ctx.RespStatusCode = http.StatusUnauthorized
            ctx.RespJSON(http.StatusUnauthorized, map[string]string{
                "error": "缺少认证头",
            })
            return
        }
        
        // 验证Bearer token格式
        if !strings.HasPrefix(authHeader, "Bearer ") {
            ctx.RespStatusCode = http.StatusUnauthorized
            ctx.RespJSON(http.StatusUnauthorized, map[string]string{
                "error": "无效的认证格式",
            })
            return
        }
        
        token := strings.TrimPrefix(authHeader, "Bearer ")
        
        // 简化的token验证（实际应用中应该验证JWT）
        if token != "valid-jwt-token" {
            ctx.RespStatusCode = http.StatusUnauthorized
            ctx.RespJSON(http.StatusUnauthorized, map[string]string{
                "error": "无效的token",
            })
            return
        }
        
        // 将用户信息存储到上下文中
        if ctx.UserValues == nil {
            ctx.UserValues = make(map[string]any)
        }
        ctx.UserValues["user_id"] = "123"
        ctx.UserValues["username"] = "testuser"
        
        next(ctx)
    }
}

func main() {
    server := ant.NewHTTPServer()
    
    // 使用认证中间件
    server.Use(JWTAuthMiddleware)
    
    // 公开路径
    server.Handle("GET /", func(ctx *ant.Context) {
        ctx.WriteString("公开页面")
    })
    
    server.Handle("POST /login", func(ctx *ant.Context) {
        // 简化的登录逻辑
        ctx.RespJSON(http.StatusOK, map[string]string{
            "token": "valid-jwt-token",
            "message": "登录成功",
        })
    })
    
    // 受保护的路径
    server.Handle("GET /profile", func(ctx *ant.Context) {
        userID := ctx.UserValues["user_id"].(string)
        username := ctx.UserValues["username"].(string)
        
        ctx.RespJSON(http.StatusOK, map[string]string{
            "user_id": userID,
            "username": username,
        })
    })
    
    server.Handle("GET /admin", func(ctx *ant.Context) {
        ctx.WriteString("管理员页面")
    })
    
    server.Run(":8080")
}
```

### 限流中间件

```go
package main

import (
    "net/http"
    "sync"
    "time"
    "github.com/justinwongcn/ant"
)

// 简单的令牌桶限流器
type TokenBucket struct {
    capacity int64         // 桶容量
    tokens   int64         // 当前令牌数
    rate     int64         // 令牌生成速率（每秒）
    lastTime time.Time     // 上次更新时间
    mutex    sync.Mutex    // 互斥锁
}

func NewTokenBucket(capacity, rate int64) *TokenBucket {
    return &TokenBucket{
        capacity: capacity,
        tokens:   capacity,
        rate:     rate,
        lastTime: time.Now(),
    }
}

func (tb *TokenBucket) Allow() bool {
    tb.mutex.Lock()
    defer tb.mutex.Unlock()
    
    now := time.Now()
    elapsed := now.Sub(tb.lastTime).Seconds()
    
    // 添加新令牌
    tb.tokens += int64(elapsed * float64(tb.rate))
    if tb.tokens > tb.capacity {
        tb.tokens = tb.capacity
    }
    
    tb.lastTime = now
    
    // 检查是否有可用令牌
    if tb.tokens > 0 {
        tb.tokens--
        return true
    }
    
    return false
}

// 限流中间件
func RateLimitMiddleware(bucket *TokenBucket) ant.Middleware {
    return func(next ant.HandleFunc) ant.HandleFunc {
        return func(ctx *ant.Context) {
            if !bucket.Allow() {
                ctx.RespStatusCode = http.StatusTooManyRequests
                ctx.RespJSON(http.StatusTooManyRequests, map[string]string{
                    "error": "请求过于频繁，请稍后再试",
                })
                return
            }
            
            next(ctx)
        }
    }
}

func main() {
    server := ant.NewHTTPServer()
    
    // 创建令牌桶：容量10，每秒生成2个令牌
    bucket := NewTokenBucket(10, 2)
    
    // 使用限流中间件
    server.Use(RateLimitMiddleware(bucket))
    
    server.Handle("GET /api/data", func(ctx *ant.Context) {
        ctx.RespJSON(http.StatusOK, map[string]string{
            "message": "数据获取成功",
            "timestamp": time.Now().Format(time.RFC3339),
        })
    })
    
    server.Run(":8080")
}
```

## 模板示例

### 基本模板渲染

```go
package main

import (
    "log"
    "github.com/justinwongcn/ant"
)

func main() {
    server := ant.NewHTTPServer()

    // 创建并配置模板引擎
    templateEngine := &ant.GoTemplateEngine{}
    err := templateEngine.LoadFromGlob("templates/*.html")
    if err != nil {
        log.Fatal("模板加载失败:", err)
    }

    // 首页
    server.Handle("GET /", func(ctx *ant.Context) {
        ctx.TemplateEngine = templateEngine

        data := struct {
            Title   string
            Message string
            Items   []string
        }{
            Title:   "欢迎页面",
            Message: "这是使用Ant框架渲染的页面",
            Items:   []string{"功能1", "功能2", "功能3"},
        }

        err := ctx.RespTemplate("index.html", data)
        if err != nil {
            ctx.WriteString("模板渲染失败: " + err.Error())
        }
    })

    server.Run(":8080")
}
```

**模板文件 templates/index.html:**
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 800px; margin: 0 auto; }
        ul { list-style-type: none; padding: 0; }
        li { background: #f0f0f0; margin: 5px 0; padding: 10px; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>{{.Title}}</h1>
        <p>{{.Message}}</p>

        <h2>功能列表</h2>
        <ul>
            {{range .Items}}
            <li>{{.}}</li>
            {{end}}
        </ul>
    </div>
</body>
</html>
```

## 文件处理示例

### 文件上传下载

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"
    "github.com/justinwongcn/ant"
)

func main() {
    server := ant.NewHTTPServer()

    // 确保上传目录存在
    uploadDir := "./uploads"
    if err := os.MkdirAll(uploadDir, 0755); err != nil {
        log.Fatal("创建上传目录失败:", err)
    }

    // 文件上传处理器
    uploader := &ant.FileUploader{
        FileField: "file",
        DstPathFunc: func(fh *multipart.FileHeader) string {
            return uploadDir
        },
        FileNameFunc: func(originalName string) string {
            // 生成唯一文件名
            ext := filepath.Ext(originalName)
            name := originalName[:len(originalName)-len(ext)]
            return fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), name, ext)
        },
    }

    // 上传页面
    server.Handle("GET /upload", func(ctx *ant.Context) {
        html := `<!DOCTYPE html>
<html>
<head>
    <title>文件上传</title>
    <meta charset="UTF-8">
</head>
<body>
    <h1>文件上传</h1>
    <form action="/upload" method="post" enctype="multipart/form-data">
        <input type="file" name="file" required>
        <button type="submit">上传文件</button>
    </form>
</body>
</html>`
        ctx.Resp.Header().Set("Content-Type", "text/html; charset=utf-8")
        ctx.WriteString(html)
    })

    // 处理文件上传
    server.Handle("POST /upload", uploader.Handle())

    // 文件下载
    downloader := &ant.FileDownloader{
        Dir: uploadDir,
    }
    server.Handle("GET /download", downloader.Handle())

    server.Run(":8080")
}
```

### 静态文件服务

```go
package main

import (
    "github.com/justinwongcn/ant"
)

func main() {
    server := ant.NewHTTPServer()

    // 静态文件处理器
    staticHandler := &ant.StaticResourceHandler{
        Dir:        "./static",
        PathPrefix: "/static/",
        ExtensionContentTypeMap: map[string]string{
            ".html": "text/html; charset=utf-8",
            ".css":  "text/css; charset=utf-8",
            ".js":   "application/javascript; charset=utf-8",
            ".png":  "image/png",
            ".jpg":  "image/jpeg",
            ".gif":  "image/gif",
            ".ico":  "image/x-icon",
        },
    }

    // 注册静态文件路由
    server.Handle("GET /static/{path...}", staticHandler.Handle())

    // 首页
    server.Handle("GET /", func(ctx *ant.Context) {
        html := `<!DOCTYPE html>
<html>
<head>
    <title>静态文件示例</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <h1>静态文件服务示例</h1>
    <p>这个页面使用了静态CSS文件。</p>
    <img src="/static/images/logo.png" alt="Logo">
    <script src="/static/js/app.js"></script>
</body>
</html>`
        ctx.Resp.Header().Set("Content-Type", "text/html; charset=utf-8")
        ctx.WriteString(html)
    })

    server.Run(":8080")
}
```

## 会话管理示例

### 基本会话管理

```go
package main

import (
    "fmt"
    "net/http"
    "time"
    "github.com/justinwongcn/ant"
    "github.com/justinwongcn/ant/session"
    "github.com/justinwongcn/ant/session/cookie"
    "github.com/justinwongcn/ant/session/memory"
    "github.com/google/uuid"
)

func main() {
    server := ant.NewHTTPServer()

    // 创建会话管理器
    sessionManager := &session.Manager{
        Store: &memory.Store{
            Expiration: 30 * time.Minute,
        },
        Propagator: &cookie.Propagator{
            CookieName: "session_id",
            CookieOption: func(c *http.Cookie) {
                c.HttpOnly = true
                c.Secure = false // 开发环境
                c.Path = "/"
                c.MaxAge = 1800 // 30分钟
            },
        },
        SessCtxKey: "session",
    }

    // 会话中间件
    server.Use(func(next ant.HandleFunc) ant.HandleFunc {
        return func(ctx *ant.Context) {
            sess, err := sessionManager.GetSession(*ctx)
            if err != nil {
                // 创建新会话
                sessionID := uuid.New().String()
                sess, err = sessionManager.InitSession(*ctx, sessionID)
                if err != nil {
                    ctx.WriteString("会话初始化失败")
                    return
                }
            }

            if ctx.UserValues == nil {
                ctx.UserValues = make(map[string]any)
            }
            ctx.UserValues["session"] = sess

            next(ctx)
        }
    })

    // 登录页面
    server.Handle("GET /login", func(ctx *ant.Context) {
        html := `<!DOCTYPE html>
<html>
<head>
    <title>登录</title>
    <meta charset="UTF-8">
</head>
<body>
    <h1>用户登录</h1>
    <form action="/login" method="post">
        <p>
            <label>用户名:</label>
            <input type="text" name="username" required>
        </p>
        <p>
            <label>密码:</label>
            <input type="password" name="password" required>
        </p>
        <button type="submit">登录</button>
    </form>
</body>
</html>`
        ctx.Resp.Header().Set("Content-Type", "text/html; charset=utf-8")
        ctx.WriteString(html)
    })

    // 处理登录
    server.Handle("POST /login", func(ctx *ant.Context) {
        username := ctx.Req.FormValue("username")
        password := ctx.Req.FormValue("password")

        // 简化的认证逻辑
        if username == "admin" && password == "password" {
            sess := ctx.UserValues["session"].(session.Session)

            // 存储用户信息到会话
            sess.Set(ctx.Req.Context(), "user_id", "123")
            sess.Set(ctx.Req.Context(), "username", username)
            sess.Set(ctx.Req.Context(), "login_time", time.Now())

            ctx.WriteString("登录成功！<a href='/profile'>查看个人资料</a>")
        } else {
            ctx.RespStatusCode = http.StatusUnauthorized
            ctx.WriteString("用户名或密码错误")
        }
    })

    // 个人资料页面
    server.Handle("GET /profile", func(ctx *ant.Context) {
        sess := ctx.UserValues["session"].(session.Session)

        userID, err := sess.Get(ctx.Req.Context(), "user_id")
        if err != nil {
            ctx.WriteString("请先登录 <a href='/login'>登录</a>")
            return
        }

        username, _ := sess.Get(ctx.Req.Context(), "username")
        loginTime, _ := sess.Get(ctx.Req.Context(), "login_time")

        html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>个人资料</title>
    <meta charset="UTF-8">
</head>
<body>
    <h1>个人资料</h1>
    <p>用户ID: %v</p>
    <p>用户名: %v</p>
    <p>登录时间: %v</p>
    <p><a href="/logout">注销</a></p>
</body>
</html>`, userID, username, loginTime)

        ctx.Resp.Header().Set("Content-Type", "text/html; charset=utf-8")
        ctx.WriteString(html)
    })

    // 注销
    server.Handle("GET /logout", func(ctx *ant.Context) {
        err := sessionManager.RemoveSession(*ctx)
        if err != nil {
            ctx.WriteString("注销失败: " + err.Error())
            return
        }

        ctx.WriteString("注销成功！<a href='/login'>重新登录</a>")
    })

    server.Run(":8080")
}
```
