# Ant Web框架API参考文档

本文档提供Ant Web框架所有公开API的详细说明。

## 目录

- [核心接口和类型](#核心接口和类型)
- [Server接口](#server接口)
- [Context类型](#context类型)
- [中间件系统](#中间件系统)
- [模板引擎](#模板引擎)
- [文件处理](#文件处理)
- [会话管理](#会话管理)
- [DDD架构组件](#ddd架构组件)

## 核心接口和类型

### HandleFunc

```go
type HandleFunc func(ctx *Context)
```

HTTP请求处理函数类型，用于定义路由处理器和中间件。

**参数：**
- `ctx *Context`: 封装了HTTP请求上下文的Context对象

**使用示例：**
```go
var handler HandleFunc = func(ctx *ant.Context) {
    ctx.WriteString("Hello, World!")
}
```

### Middleware

```go
type Middleware func(next HandleFunc) HandleFunc
```

中间件函数类型，接收下一个处理器并返回一个新的处理器。

**参数：**
- `next HandleFunc`: 链中的下一个处理器

**返回值：**
- `HandleFunc`: 包装后的处理器

**使用示例：**
```go
func LoggingMiddleware(next ant.HandleFunc) ant.HandleFunc {
    return func(ctx *ant.Context) {
        log.Printf("请求: %s %s", ctx.Req.Method, ctx.Req.URL.Path)
        next(ctx)
    }
}
```

## Server接口

### Server

```go
type Server interface {
    http.Handler
    Handle(pattern string, handler HandleFunc)
    Run(addr string) error
}
```

HTTP服务器接口，提供路由注册和服务器启动的核心功能。

#### Handle

```go
Handle(pattern string, handler HandleFunc)
```

注册路由处理器。

**参数：**
- `pattern string`: 路由模式，支持Go 1.22+的新路由语法
  - `"GET /users"` - 基本路由
  - `"GET /users/{id}"` - 参数化路由
  - `"GET /files/{path...}"` - 通配符路由
  - `"GET /posts/{$}"` - 精确匹配
- `handler HandleFunc`: 处理该路由的处理函数

**使用示例：**
```go
server.Handle("GET /users/{id}", func(ctx *ant.Context) {
    id, _ := ctx.PathValue("id").String()
    ctx.WriteString(fmt.Sprintf("用户ID: %s", id))
})
```

#### Run

```go
Run(addr string) error
```

启动HTTP服务器。

**参数：**
- `addr string`: 监听地址，格式为 "host:port"

**返回值：**
- `error`: 服务器运行过程中发生的错误

**使用示例：**
```go
err := server.Run(":8080")
if err != nil {
    log.Fatal("服务器启动失败:", err)
}
```

### HTTPServer

```go
type HTTPServer struct {
    // 私有字段
}
```

Server接口的默认实现。

#### NewHTTPServer

```go
func NewHTTPServer() Server
```

创建新的HTTP服务器实例。

**返回值：**
- `Server`: HTTP服务器实例

**使用示例：**
```go
server := ant.NewHTTPServer()
```

#### Use

```go
func (s *HTTPServer) Use(middlewares ...Middleware)
```

添加全局中间件。

**参数：**
- `middlewares ...Middleware`: 一个或多个中间件函数

**使用示例：**
```go
server.Use(LoggingMiddleware, AuthMiddleware)
```

## Context类型

### Context

```go
type Context struct {
    Req  *http.Request       // 原始HTTP请求对象
    Resp http.ResponseWriter // HTTP响应写入器

    // 缓存解析后的URL查询参数，避免重复解析
    cacheQueryValues url.Values

    // 响应缓存数据，在最终响应时一次性写入
    RespStatusCode int    // 响应状态码
    RespData       []byte // 响应内容主体

    // 模板引擎，用于渲染HTML模板
    TemplateEngine TemplateEngine

    // 用户相关的数据，用于在请求处理过程中存储临时数据
    UserValues map[string]any
}
```

封装HTTP请求上下文，提供请求处理相关工具方法。

### 路径参数方法

#### PathValue

```go
func (c *Context) PathValue(key string) StringValue
```

获取路径参数值。

**参数：**
- `key string`: 参数名称

**返回值：**
- `StringValue`: 包装的字符串值，提供类型转换方法

**使用示例：**
```go
// 路由: GET /users/{id}
id := ctx.PathValue("id")
idStr, err := id.String()
if err != nil {
    ctx.WriteString("无效的用户ID")
    return
}
```

### 查询参数方法

#### QueryValue

```go
func (c *Context) QueryValue(key string) StringValue
```

获取查询参数值。

**参数：**
- `key string`: 参数名称

**返回值：**
- `StringValue`: 包装的字符串值

**使用示例：**
```go
// URL: /search?q=golang
keyword := ctx.QueryValue("q")
keywordStr, _ := keyword.String()
```

#### DefaultQueryValue

```go
func (c *Context) DefaultQueryValue(key, def string) StringValue
```

获取查询参数值，如果不存在则返回默认值。

**参数：**
- `key string`: 参数名称
- `def string`: 默认值

**返回值：**
- `StringValue`: 包装的字符串值

**使用示例：**
```go
page, _ := ctx.DefaultQueryValue("page", "1").String()
```

### 响应方法

#### WriteString

```go
func (c *Context) WriteString(data string) error
```

将字符串写入响应体。

**参数：**
- `data string`: 要写入的字符串数据

**返回值：**
- `error`: 写入过程中的错误

**使用示例：**
```go
ctx.WriteString("Hello, World!")
```

#### RespJSON

```go
func (c *Context) RespJSON(status int, obj any) error
```

将对象序列化为JSON并写入响应。

**参数：**
- `status int`: HTTP状态码
- `obj any`: 要序列化的对象

**返回值：**
- `error`: 序列化或写入过程中的错误

**使用示例：**
```go
user := map[string]interface{}{
    "id":   1,
    "name": "张三",
}
ctx.RespJSON(http.StatusOK, user)
```

#### RespTemplate

```go
func (c *Context) RespTemplate(tplName string, data any) error
```

渲染HTML模板并将结果写入响应。

**参数：**
- `tplName string`: 模板名称
- `data any`: 渲染数据

**返回值：**
- `error`: 渲染过程中的错误

**使用示例：**
```go
ctx.TemplateEngine = templateEngine
data := struct {
    Title string
    Content string
}{
    Title: "欢迎页面",
    Content: "这是模板渲染的内容",
}
ctx.RespTemplate("index.html", data)
```

### StringValue类型

```go
type StringValue struct {
    val string
    err error
}
```

包装字符串值，提供类型转换方法。

#### String

```go
func (s StringValue) String() (string, error)
```

获取字符串值。

**返回值：**
- `string`: 字符串值
- `error`: 转换错误

#### ToInt64

```go
func (s StringValue) ToInt64() (int64, error)
```

转换为int64类型。

**返回值：**
- `int64`: 转换后的整数值
- `error`: 转换错误

**使用示例：**
```go
id := ctx.PathValue("id")
idInt, err := id.ToInt64()
if err != nil {
    ctx.WriteString("无效的ID格式")
    return
}
```

## 中间件系统

### 内置中间件

#### 访问日志中间件

```go
package accesslog

type MiddlewareBuilder struct {
    logFunc func(accessLog string)
}

func NewBuilder() *MiddlewareBuilder
func (b *MiddlewareBuilder) LogFunc(lFn func(accessLog string)) *MiddlewareBuilder
func (b *MiddlewareBuilder) Build() ant.Middleware
func AccessLog() ant.Middleware
```

**使用示例：**
```go
import "github.com/justinwongcn/ant/middleware/accesslog"

// 使用默认配置
server.Use(accesslog.AccessLog())

// 自定义日志函数
server.Use(accesslog.NewBuilder().
    LogFunc(func(log string) {
        fmt.Println("自定义日志:", log)
    }).
    Build())
```

#### 错误处理中间件

```go
package errhandle

type MiddlewareBuilder struct {
    // 私有字段
}

func NewBuilder() *MiddlewareBuilder
func (m *MiddlewareBuilder) RegisterError(code int, handler ant.HandleFunc) *MiddlewareBuilder
func (m *MiddlewareBuilder) Build() ant.Middleware
```

**使用示例：**
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
package recovery

type MiddlewareBuilder struct {
    // 私有字段
}

func NewBuilder() *MiddlewareBuilder
func (r *MiddlewareBuilder) StatusCode(code int) *MiddlewareBuilder
func (r *MiddlewareBuilder) ErrMsg(msg string) *MiddlewareBuilder
func (r *MiddlewareBuilder) LogFunc(fn func(ctx context.Context)) *MiddlewareBuilder
func (r *MiddlewareBuilder) Build() ant.Middleware
```

**使用示例：**
```go
import "github.com/justinwongcn/ant/middleware/recovery"

server.Use(recovery.NewBuilder().
    StatusCode(http.StatusInternalServerError).
    ErrMsg("服务器内部错误").
    LogFunc(func(ctx context.Context) {
        log.Println("发生panic，已恢复")
    }).
    Build())
```

## 模板引擎

### TemplateEngine接口

```go
type TemplateEngine interface {
    Render(ctx context.Context, tplName string, data any) ([]byte, error)
}
```

模板引擎接口，提供模板渲染的核心功能。

#### Render

```go
Render(ctx context.Context, tplName string, data any) ([]byte, error)
```

渲染模板。

**参数：**
- `ctx context.Context`: 上下文对象
- `tplName string`: 模板名称
- `data any`: 渲染数据

**返回值：**
- `[]byte`: 渲染后的内容
- `error`: 渲染错误

### GoTemplateEngine

```go
type GoTemplateEngine struct {
    T *template.Template
}
```

基于Go标准库html/template的模板引擎实现。

#### LoadFromGlob

```go
func (g *GoTemplateEngine) LoadFromGlob(pattern string) error
```

从指定的glob模式加载模板。

**参数：**
- `pattern string`: glob模式字符串

**返回值：**
- `error`: 加载错误

**使用示例：**
```go
engine := &ant.GoTemplateEngine{}
err := engine.LoadFromGlob("templates/*.html")
if err != nil {
    log.Fatal("模板加载失败:", err)
}
```

#### LoadFromFiles

```go
func (g *GoTemplateEngine) LoadFromFiles(files ...string) error
```

从指定的文件列表加载模板。

**参数：**
- `files ...string`: 模板文件路径列表

**返回值：**
- `error`: 加载错误

**使用示例：**
```go
engine := &ant.GoTemplateEngine{}
err := engine.LoadFromFiles("templates/index.html", "templates/layout.html")
```

#### LoadFromFS

```go
func (g *GoTemplateEngine) LoadFromFS(fs fs.FS, paths ...string) error
```

从文件系统加载模板。

**参数：**
- `fs fs.FS`: 文件系统接口
- `paths ...string`: 模板文件路径列表

**返回值：**
- `error`: 加载错误

**使用示例：**
```go
//go:embed templates/*
var templateFS embed.FS

engine := &ant.GoTemplateEngine{}
err := engine.LoadFromFS(templateFS, "templates/*.html")
```

## 文件处理

### FileUploader

```go
type FileUploader struct {
    FileField    string
    DstPathFunc  func(fh *multipart.FileHeader) string
    FileNameFunc func(originalName string) string
}
```

文件上传处理器，提供安全的文件上传功能。

**字段说明：**
- `FileField string`: 表单中文件字段的名称
- `DstPathFunc func(*multipart.FileHeader) string`: 确定目标存储路径的函数
- `FileNameFunc func(string) string`: 生成文件名的函数，可选

#### Handle

```go
func (f *FileUploader) Handle() HandleFunc
```

返回文件上传处理函数。

**返回值：**
- `HandleFunc`: 处理上传请求的处理函数

**使用示例：**
```go
uploader := &ant.FileUploader{
    FileField: "file",
    DstPathFunc: func(fh *multipart.FileHeader) string {
        return "./uploads/"
    },
    FileNameFunc: func(originalName string) string {
        ext := filepath.Ext(originalName)
        return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
    },
}

server.Handle("POST /upload", uploader.Handle())
```

### FileDownloader

```go
type FileDownloader struct {
    Dir string
}
```

文件下载处理器，提供安全的文件下载功能。

**字段说明：**
- `Dir string`: 文件下载的根目录

#### Handle

```go
func (f *FileDownloader) Handle() HandleFunc
```

返回文件下载处理函数。

**返回值：**
- `HandleFunc`: 处理下载请求的处理函数

**使用示例：**
```go
downloader := &ant.FileDownloader{
    Dir: "./uploads",
}

server.Handle("GET /download", downloader.Handle())
// 使用方式: GET /download?file=example.txt
```

### StaticResourceHandler

```go
type StaticResourceHandler struct {
    Dir                     string
    PathPrefix              string
    ExtensionContentTypeMap map[string]string
}
```

静态资源处理器，提供静态文件服务功能。

**字段说明：**
- `Dir string`: 静态文件根目录
- `PathPrefix string`: URL路径前缀
- `ExtensionContentTypeMap map[string]string`: 文件扩展名到Content-Type的映射

#### Handle

```go
func (h *StaticResourceHandler) Handle() HandleFunc
```

返回静态资源处理函数。

**使用示例：**
```go
staticHandler := &ant.StaticResourceHandler{
    Dir:        "./static",
    PathPrefix: "/static/",
    ExtensionContentTypeMap: map[string]string{
        ".html": "text/html; charset=utf-8",
        ".css":  "text/css; charset=utf-8",
        ".js":   "application/javascript; charset=utf-8",
        ".png":  "image/png",
        ".jpg":  "image/jpeg",
    },
}

server.Handle("GET /static/{path...}", staticHandler.Handle())
```

### StaticResourceHandlerWithCache

```go
type StaticResourceHandlerWithCache struct {
    StaticResourceHandler
    MaxSize int
}
```

带缓存的静态资源处理器，继承StaticResourceHandler并添加LRU缓存。

**字段说明：**
- `StaticResourceHandler`: 嵌入的基础静态资源处理器
- `MaxSize int`: 缓存最大条目数

**使用示例：**
```go
cachedHandler := &ant.StaticResourceHandlerWithCache{
    StaticResourceHandler: ant.StaticResourceHandler{
        Dir:        "./static",
        PathPrefix: "/static/",
        ExtensionContentTypeMap: map[string]string{
            ".css": "text/css; charset=utf-8",
            ".js":  "application/javascript; charset=utf-8",
        },
    },
    MaxSize: 1000,
}

server.Handle("GET /static/{path...}", cachedHandler.Handle())
```

## 会话管理

### Session接口

```go
type Session interface {
    Get(ctx context.Context, key string) (any, error)
    Set(ctx context.Context, key string, val any) error
    ID() string
}
```

会话接口，定义会话的基本操作。

#### Get

```go
Get(ctx context.Context, key string) (any, error)
```

从会话中获取值。

**参数：**
- `ctx context.Context`: 上下文
- `key string`: 键名

**返回值：**
- `any`: 存储的值
- `error`: 获取错误

#### Set

```go
Set(ctx context.Context, key string, val any) error
```

向会话中设置值。

**参数：**
- `ctx context.Context`: 上下文
- `key string`: 键名
- `val any`: 要存储的值

**返回值：**
- `error`: 设置错误

#### ID

```go
ID() string
```

获取会话ID。

**返回值：**
- `string`: 会话ID

### Store接口

```go
type Store interface {
    Generate(ctx context.Context, id string) (Session, error)
    Get(ctx context.Context, id string) (Session, error)
    Refresh(ctx context.Context, id string) error
    Remove(ctx context.Context, id string) error
}
```

会话存储接口，定义会话的存储操作。

### Propagator接口

```go
type Propagator interface {
    Inject(id string, writer http.ResponseWriter) error
    Extract(request *http.Request) (string, error)
    Remove(writer http.ResponseWriter) error
}
```

会话传播器接口，负责会话ID在HTTP请求和响应之间的传递。

### Manager

```go
type Manager struct {
    Store
    Propagator
    SessCtxKey string
}
```

会话管理器，组合Store和Propagator接口。

#### GetSession

```go
func (m *Manager) GetSession(ctx ant.Context) (Session, error)
```

获取会话。

**参数：**
- `ctx ant.Context`: 请求上下文

**返回值：**
- `Session`: 会话实例
- `error`: 获取错误

#### InitSession

```go
func (m *Manager) InitSession(ctx ant.Context, id string) (Session, error)
```

初始化新会话。

**参数：**
- `ctx ant.Context`: 请求上下文
- `id string`: 会话ID

**返回值：**
- `Session`: 新创建的会话实例
- `error`: 初始化错误

#### RefreshSession

```go
func (m *Manager) RefreshSession(ctx ant.Context) (Session, error)
```

刷新会话过期时间。

**参数：**
- `ctx ant.Context`: 请求上下文

**返回值：**
- `Session`: 刷新后的会话实例
- `error`: 刷新错误

#### RemoveSession

```go
func (m *Manager) RemoveSession(ctx ant.Context) error
```

删除会话。

**参数：**
- `ctx ant.Context`: 请求上下文

**返回值：**
- `error`: 删除错误

### 内置实现

#### memory.Store

```go
package memory

type Store struct {
    Expiration time.Duration
}
```

基于内存的会话存储实现。

**使用示例：**
```go
import "github.com/justinwongcn/ant/session/memory"

store := &memory.Store{
    Expiration: 30 * time.Minute,
}
```

#### cookie.Propagator

```go
package cookie

type Propagator struct {
    CookieName   string
    CookieOption func(c *http.Cookie)
}
```

基于Cookie的会话传播器实现。

**使用示例：**
```go
import "github.com/justinwongcn/ant/session/cookie"

propagator := &cookie.Propagator{
    CookieName: "session_id",
    CookieOption: func(c *http.Cookie) {
        c.HttpOnly = true
        c.Secure = true
        c.Path = "/"
        c.MaxAge = 1800
    },
}
```
```
