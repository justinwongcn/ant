package accesslog

import (
	"encoding/json"
	"log"
	"time"

	"github.com/justinwongcn/ant"
)

// accessLog 定义访问日志的结构
type accessLog struct {
	Timestamp  string        `json:"timestamp"`
	Host       string        `json:"host"`
	HTTPMethod string        `json:"http_method"`
	Path       string        `json:"path"`
	Duration   time.Duration `json:"duration"`
}

// MiddlewareBuilder 中间件构建器
type MiddlewareBuilder struct {
	logFunc func(accessLog string)
}

// LogFunc 设置自定义日志记录函数
func (b *MiddlewareBuilder) LogFunc(lFn func(accessLog string)) *MiddlewareBuilder {
	b.logFunc = lFn
	return b
}

// NewBuilder 创建中间件构建器
func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(accessLog string) {
			log.Println(accessLog)
		},
	}
}

// Build 构建访问日志中间件
func (b *MiddlewareBuilder) Build() ant.Middleware {
	return func(next ant.HandleFunc) ant.HandleFunc {
		return func(ctx *ant.Context) {
			start := time.Now()

			// 执行下一个处理器
			next(ctx)

			// 构建访问日志
			l := accessLog{
				Timestamp:  start.Format("2006-01-02 15:04:05"),
				Host:       ctx.Req.Host,
				HTTPMethod: ctx.Req.Method,
				Path:       ctx.Req.URL.Path,
				Duration:   time.Since(start),
			}

			// 序列化并记录日志
			val, _ := json.Marshal(l)
			b.logFunc(string(val))
		}
	}
}

// AccessLog 创建默认的访问日志中间件
func AccessLog() ant.Middleware {
	return NewBuilder().Build()
}
