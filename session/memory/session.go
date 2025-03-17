package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/justinwongcn/ant/session"
	cache "github.com/patrickmn/go-cache"
)

// Store 内存会话存储实现
// 利用内存缓存来管理会话的存储和过期时间
type Store struct {
	// c 内存缓存实例，用于管理会话数据和过期时间
	c *cache.Cache
	// expiration 会话的过期时间
	expiration time.Duration
}

// NewStore 创建一个 Store 的实例
// expiration: 会话的过期时间
// 返回值: 创建的 Store 实例
// 实际上，这里也可以考虑使用 Option 设计模式，允许用户控制过期检查的间隔
func NewStore(expiration time.Duration) *Store {
	return &Store{
		c:          cache.New(expiration, expiration),
		expiration: expiration,
	}
}

// memorySession 内存会话实例
// 实现了 session.Session 接口
type memorySession struct {
	// id 会话的唯一标识符
	id string
	// data 存储会话数据的映射
	data map[string]any
	// expiration 会话的过期时间
	expiration time.Duration
	// mu 保护 data 的互斥锁
	mu sync.Mutex
}

// Get 获取会话中的数据
// ctx: 上下文（当前未使用）
// key: 数据的键
// 返回值:
// - 获取到的数据
// - 如果键不存在则返回错误
func (m *memorySession) Get(_ context.Context, key string) (any, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	val, ok := m.data[key]
	if !ok {
		return "", errors.New("找不到这个 key")
	}

	return val, nil
}

// Set 设置会话中的数据
// ctx: 上下文（当前未使用）
// key: 数据的键
// value: 要存储的数据
// 返回值: 设置过程中的错误
func (m *memorySession) Set(_ context.Context, key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}

// ID 获取会话ID
// 返回值: 会话的唯一标识符
func (m *memorySession) ID() string {
	return m.id
}

// Generate 生成一个新的会话
// ctx: 上下文（当前未使用）
// id: 会话ID
// 返回值:
// - 生成的会话实例
// - 可能发生的错误
func (m *Store) Generate(_ context.Context, id string) (session.Session, error) {
	sess := &memorySession{
		id:         id,
		data:       make(map[string]any),
		expiration: m.expiration,
	}

	m.c.Set(sess.ID(), sess, m.expiration)

	return sess, nil
}

// Refresh 刷新会话
// ctx: 上下文
// id: 会话ID
// 返回值: 刷新过程中的错误
func (m *Store) Refresh(ctx context.Context, id string) error {
	sess, err := m.Get(ctx, id)
	if err != nil {
		return err
	}

	m.c.Set(sess.ID(), sess, m.expiration)

	return nil
}

// Remove 删除会话
// ctx: 上下文（当前未使用）
// id: 要删除的会话ID
// 返回值: 删除过程中的错误
func (m *Store) Remove(_ context.Context, id string) error {
	m.c.Delete(id)
	return nil
}

// Get 获取会话
// ctx: 上下文（当前未使用）
// id: 会话ID
// 返回值:
// - 获取到的会话实例
// - 如果会话不存在则返回错误
func (m *Store) Get(_ context.Context, id string) (session.Session, error) {
	sess, ok := m.c.Get(id)
	if !ok {
		return nil, errors.New("session not found")
	}

	return sess.(*memorySession), nil
}
