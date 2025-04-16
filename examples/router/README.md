# 路由示例

这个示例展示了 Ant Web 框架的三种不同路由处理方式：

1. 固定路径：`GET /hello`
   - 返回："Hello, world"
   - 示例：`curl http://localhost:8080/hello`

2. 通配符路径：`GET /hello/{name}`
   - 返回："Hello, {name}"
   - 示例：`curl http://localhost:8080/hello/ant`

3. 查询参数：`GET /hello-query?name=xxx`
   - 返回："hello {name}"
   - 示例：`curl http://localhost:8080/hello-query?name=ant`

## 运行示例

```bash
# 进入示例目录
cd examples/router

# 运行示例
go run main.go
```

服务器将在 http://localhost:8080 启动，你可以使用浏览器或 curl 命令来测试不同的路由。