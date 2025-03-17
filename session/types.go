package session

import (
	"context"
	"net/http"
)

// Session 表示一个会话实例
type Session interface {
	// Get 获取会话中的数据
	Get(ctx context.Context, key string) (any, error)
	// Set 设置会话中的数据
	Set(ctx context.Context, key string, value any) error
	// ID 获取会话ID
	ID() (id string)
}

// Store 定义会话存储接口
type Store interface {
	// Generate 生成一个新的会话
	// ctx: 上下文
	// id: 会话ID
	// 返回值:
	// - 生成的会话实例
	// - 可能发生的错误
	Generate(ctx context.Context, id string) (Session, error)

	// Refresh 刷新会话
	// ctx: 上下文
	// id: 会话ID
	// 返回值: 刷新过程中的错误
	Refresh(ctx context.Context, id string) error

	// Remove 删除会话
	// ctx: 上下文
	// id: 要删除的会话ID
	// 返回值: 删除过程中的错误
	Remove(ctx context.Context, id string) error

	// Get 获取会话
	// ctx: 上下文
	// id: 会话ID
	// 返回值:
	// - 获取到的会话实例
	// - 可能发生的错误
	Get(ctx context.Context, id string) (Session, error)
}

// Propagator 定义会话传播器接口
type Propagator interface {
	// Inject 将会话ID注入到HTTP响应中
	// id: 会话ID
	// writer: HTTP响应写入器
	// 返回值: 注入过程中的错误
	Inject(id string, writer http.ResponseWriter) error

	// Extract 从HTTP请求中提取会话ID
	// req: HTTP请求
	// 返回值:
	// - 提取的会话ID
	// - 可能发生的错误
	Extract(req *http.Request) (id string, err error)

	// Remove 从HTTP响应中移除会话ID
	// writer: HTTP响应写入器
	// 返回值: 移除过程中的错误
	Remove(writer http.ResponseWriter) error
}
