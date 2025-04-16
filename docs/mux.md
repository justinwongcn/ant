# Go 1.18 多路路由新特性详解

## 概述

Go 1.18 版本对标准库 `net/http` 中的 `http.ServeMux` 进行了重大增强，引入了更强大的路由功能，包括 HTTP 方法匹配、通配符路径参数等特性。这些改进显著缩小了标准库与第三方路由库(如 gorilla/mux)之间的功能差距，使开发者能够在不引入外部依赖的情况下实现复杂的路由需求。

## 新增功能

### 1. HTTP 方法匹配

现在可以在路由模式中直接指定 HTTP 方法，使同一路径可以基于不同方法路由到不同的处理函数。

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /path/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "处理GET请求")
})
mux.HandleFunc("POST /path/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "处理POST请求")
})
```

### 2. 通配符路径参数

支持在路径中使用 `{name}` 形式的通配符来捕获动态路径段。

```go
mux.HandleFunc("/articles/{id}/", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id") // 获取路径参数
    fmt.Fprintf(w, "文章ID: %s", id)
})
```

### 3. 通配符扩展匹配

支持 `{name...}` 形式的通配符来匹配路径的剩余部分。

```go
mux.HandleFunc("/files/{path...}", func(w http.ResponseWriter, r *http.Request) {
    path := r.PathValue("path")
    fmt.Fprintf(w, "文件路径: %s", path)
})
```

### 4. 严格路径匹配

使用 `{$}` 可以确保路径精确匹配。

```go
mux.HandleFunc("/exact{$}", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "精确匹配路径")
})
```

## 路由命中规则

### 1. 匹配优先级

当多个模式可能匹配同一请求时，ServeMux 会按照以下规则确定哪个模式优先:

1. **具体路径优先于通配符路径**：例如 `/articles/featured` 优先于 `/articles/{id}`
2. **较长路径优先于较短路径**：例如 `/articles/{id}/details` 优先于 `/articles/{id}`
3. **显式方法优先于无方法限制**：例如 `GET /path` 优先于 `/path`
4. **严格匹配优先于非严格匹配**：例如 `/path{$}` 优先于 `/path`

### 2. 冲突处理

如果两个模式具有相同的优先级且都能匹配同一请求，ServeMux 会在注册时 panic，防止模糊匹配。

例如，以下两个模式会引发冲突:

```go
mux.HandleFunc("/articles/{id}/status/", handler1)
mux.HandleFunc("/articles/0/{action}/", handler2)
```

冲突错误信息会详细说明问题:

```
panic: pattern "/articles/0/{action}/" conflicts with pattern "/articles/{id}/status/":
/articles/0/{action}/ and /articles/{id}/status/ both match some paths, like "/articles/0/status/".
But neither is more specific than the other.
/articles/0/{action}/ matches "/articles/0/action/", but /articles/{id}/status/ doesn't.
/articles/{id}/status/ matches "/articles/id/status/", but /articles/0/{action}/ doesn't.
```

## 实用示例

### 1. REST API 实现

```go
func main() {
    mux := http.NewServeMux()
    server := NewTaskServer()
    
    mux.HandleFunc("POST /tasks/", server.createTaskHandler)
    mux.HandleFunc("GET /tasks/", server.getAllTasksHandler)
    mux.HandleFunc("DELETE /tasks/", server.deleteAllTasksHandler)
    mux.HandleFunc("GET /tasks/{id}/", server.getTaskHandler)
    mux.HandleFunc("DELETE /tasks/{id}/", server.deleteTaskHandler)
    mux.HandleFunc("GET /tags/{tag}/", server.tagHandler)
    mux.HandleFunc("GET /due/{year}/{month}/{day}/", server.dueHandler)
    
    http.ListenAndServe(":8080", mux)
}
```

### 2. 路由辅助函数封装

可以封装辅助函数简化路由注册:

```go
type Router struct {
    *http.ServeMux
}

func (r *Router) ByMethod(method string, path string, fn http.HandlerFunc) {
    r.HandleFunc(fmt.Sprintf("%s %s", method, path), fn)
}

func (r *Router) Get(path string, fn http.HandlerFunc) {
    r.ByMethod(http.MethodGet, path, fn)
}

func (r *Router) Post(path string, fn http.HandlerFunc) {
    r.ByMethod(http.MethodPost, path, fn)
}

// 使用示例
router := &Router{http.NewServeMux()}
router.Get("/api/v1/catalog", catalog.List)
router.Post("/api/v1/catalog", catalog.Create)
```

## 注意事项

1. **性能考虑**：新的路由实现比旧版本稍慢，但比大多数第三方路由库更快
2. **兼容性**：旧的路由模式仍然支持，可以逐步迁移
3. **调试**：使用 `go tool vet` 检查潜在的路由冲突
4. **测试**：建议对新路由进行充分测试，特别是涉及通配符匹配的场景

## 总结

Go 1.18 的路由增强使标准库 `http.ServeMux` 成为一个功能完备的路由解决方案，适用于大多数 Web 应用场景。这些改进减少了对外部路由库的依赖，简化了项目依赖管理，同时保持了 Go 标准库的简洁性和高性能特点。