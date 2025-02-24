package ant

import (
	"net/http"
)

type Context struct {
	Req    *http.Request
	Resp   http.ResponseWriter
	Params func(string) string // Go 1.22 的路径参数获取函数
}

// PathParam 获取路径参数
func (c *Context) PathParam(name string) string {
	return c.Params(name)
}

// Query 获取查询参数
func (c *Context) Query(name string) string {
	return c.Req.URL.Query().Get(name)
}

