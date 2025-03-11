package session

import (
	"github.com/justinwongcn/ant"
)

// Manager 会话管理器
// 组合了 Store 和 Propagator 接口，用于管理会话的完整生命周期
// Store: 负责会话的存储和检索
// Propagator: 负责会话ID在HTTP请求和响应之间的传递
// SessCtxKey: 用于在上下文中存储会话的键名
type Manager struct {
	Store
	Propagator
	SessCtxKey string
}

// GetSession 获取会话
// ctx: 上下文，包含请求和响应信息
// 返回值:
// - 获取到的会话实例
// - 可能发生的错误
// 首先尝试从上下文中获取会话，如果不存在则从请求中提取会话ID并获取会话
func (m *Manager) GetSession(ctx ant.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}

	val, ok := ctx.UserValues[m.SessCtxKey]
	if ok {
		return val.(Session), nil
	}

	id, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}

	sess, err := m.Get(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	ctx.UserValues[m.SessCtxKey] = sess

	return sess, nil
}

// InitSession 初始化会话
// ctx: 上下文，包含请求和响应信息
// id: 新会话的唯一标识符
// 返回值:
// - 创建的会话实例
// - 可能发生的错误
// 生成新的会话并将其注入到HTTP响应中
func (m *Manager) InitSession(ctx ant.Context, id string) (Session, error) {
	sess, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}

	if err = m.Inject(id, ctx.Resp); err != nil {
		return nil, err
	}
	return sess, nil
}

// RefreshSession 刷新会话
// ctx: 上下文，包含请求和响应信息
// 返回值:
// - 刷新后的会话实例
// - 可能发生的错误
// 刷新会话的过期时间并重新注入到HTTP响应中
func (m *Manager) RefreshSession(ctx ant.Context) (Session, error) {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return nil, err
	}

	// 刷新存储的过期时间
	err = m.Refresh(ctx.Req.Context(), sess.ID())
	if err != nil {
		return nil, err
	}

	// 重新注入到HTTP响应中
	if err = m.Inject(sess.ID(), ctx.Resp); err != nil {
		return nil, err
	}
	return sess, nil
}

// RemoveSession 删除会话
// ctx: 上下文，包含请求和响应信息
// 返回值: 删除过程中的错误
// 从存储中删除会话并从HTTP响应中移除会话ID
func (m *Manager) RemoveSession(ctx ant.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}

	// 删除会话
	err = m.Store.Remove(ctx.Req.Context(), sess.ID())
	if err != nil {
		return err
	}

	// 从HTTP响应中移除会话ID
	return m.Propagator.Remove(ctx.Resp)
}
