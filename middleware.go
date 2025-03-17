package ant

// Middleware 定义中间件类型
// 中间件函数接收下一个处理器，返回一个新的处理器
type Middleware func(next HandleFunc) HandleFunc
