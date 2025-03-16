// 使用 session 包名，但将测试放在单独的文件中
package session

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justinwongcn/ant"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 自定义错误
var errSessionNotFoundForTest = errors.New("session not found for test")

// 创建一个增强的模拟传播器，用于集成测试
type mockPropagatorForIntegration struct {
	*mockPropagator
	req *http.Request // 用于模拟请求
}

func newMockPropagatorForIntegration() *mockPropagatorForIntegration {
	return &mockPropagatorForIntegration{
		mockPropagator: newMockPropagator(),
	}
}

// Extract 重写Extract方法，使用req字段
func (m *mockPropagatorForIntegration) Extract(req *http.Request) (string, error) {
	// 如果设置了特定的请求，使用它
	if m.req != nil && req == m.req {
		// 如果会话存在，返回会话ID
		for id, exists := range m.sessions {
			if exists {
				return id, nil
			}
		}
		return "", errSessionNotFoundForTest
	}
	return m.mockPropagator.Extract(req)
}

// Remove 重写Remove方法，确保会话被正确删除
func (m *mockPropagatorForIntegration) Remove(writer http.ResponseWriter) error {
	// 删除所有会话
	for id := range m.sessions {
		m.sessions[id] = false
	}
	return nil
}

// TestSessionIntegration 测试会话管理器与存储和传播器的集成
func TestSessionIntegration(t *testing.T) {
	// 创建模拟的存储和传播器
	store := newMockStore()
	propagator := newMockPropagatorForIntegration()

	// 创建会话管理器
	manager := &Manager{
		Store:      store,
		Propagator: propagator,
		SessCtxKey: "session",
	}

	// 测试完整的会话生命周期
	t.Run("完整会话生命周期", func(t *testing.T) {
		// 创建HTTP请求和响应
		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		resp := httptest.NewRecorder()

		// 创建上下文
		ctx := ant.Context{
			Req:        req,
			Resp:       resp,
			UserValues: make(map[string]any),
		}

		// 1. 初始化会话
		sess, err := manager.InitSession(ctx, "test-session-id")
		require.NoError(t, err)
		require.NotNil(t, sess)
		assert.Equal(t, "test-session-id", sess.ID())

		// 验证会话已设置
		assert.True(t, propagator.sessions["test-session-id"])

		// 2. 设置会话数据
		err = sess.Set(context.Background(), "username", "testuser")
		require.NoError(t, err)

		// 3. 获取会话
		// 创建新的请求，模拟带上Cookie
		req2 := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		propagator.req = req2 // 设置请求，使Extract能够返回会话ID
		resp2 := httptest.NewRecorder()

		ctx2 := ant.Context{
			Req:        req2,
			Resp:       resp2,
			UserValues: make(map[string]any),
		}

		sess2, err := manager.GetSession(ctx2)
		require.NoError(t, err)
		require.NotNil(t, sess2)
		assert.Equal(t, "test-session-id", sess2.ID())

		// 4. 获取会话数据
		val, err := sess2.Get(context.Background(), "username")
		require.NoError(t, err)
		assert.Equal(t, "testuser", val)

		// 5. 刷新会话
		sess3, err := manager.RefreshSession(ctx2)
		require.NoError(t, err)
		require.NotNil(t, sess3)
		assert.Equal(t, "test-session-id", sess3.ID())

		// 验证会话已刷新
		assert.True(t, propagator.sessions["test-session-id"])

		// 6. 删除会话
		err = manager.RemoveSession(ctx2)
		require.NoError(t, err)

		// 验证会话已删除 - 在删除后，会话应该不存在或为false
		assert.False(t, propagator.sessions["test-session-id"])

		// 7. 尝试获取已删除的会话
		// 创建新的请求，不带Cookie
		req3 := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		propagator.req = req3 // 设置请求，使Extract返回错误
		resp3 := httptest.NewRecorder()

		ctx3 := ant.Context{
			Req:        req3,
			Resp:       resp3,
			UserValues: make(map[string]any),
		}

		_, err = manager.GetSession(ctx3)
		assert.Error(t, err) // 应该返回错误，因为会话已删除
	})

	// 测试会话数据持久性
	t.Run("会话数据持久性", func(t *testing.T) {
		// 重置模拟对象
		store = newMockStore()
		propagator = newMockPropagatorForIntegration()
		manager = &Manager{
			Store:      store,
			Propagator: propagator,
			SessCtxKey: "session",
		}

		// 创建HTTP请求和响应
		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		resp := httptest.NewRecorder()

		// 创建上下文
		ctx := ant.Context{
			Req:        req,
			Resp:       resp,
			UserValues: make(map[string]any),
		}

		// 初始化会话
		sess, err := manager.InitSession(ctx, "persist-session-id")
		require.NoError(t, err)

		// 设置多个会话数据
		err = sess.Set(context.Background(), "user_id", 12345)
		require.NoError(t, err)
		err = sess.Set(context.Background(), "is_admin", true)
		require.NoError(t, err)
		err = sess.Set(context.Background(), "preferences", map[string]string{"theme": "dark"})
		require.NoError(t, err)

		// 创建新的请求，模拟带上Cookie
		req2 := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		propagator.req = req2 // 设置请求，使Extract能够返回会话ID
		propagator.sessions["persist-session-id"] = true
		resp2 := httptest.NewRecorder()

		ctx2 := ant.Context{
			Req:        req2,
			Resp:       resp2,
			UserValues: make(map[string]any),
		}

		// 获取会话
		sess2, err := manager.GetSession(ctx2)
		require.NoError(t, err)

		// 验证所有数据都正确持久化
		val1, err := sess2.Get(context.Background(), "user_id")
		require.NoError(t, err)
		assert.Equal(t, 12345, val1)

		val2, err := sess2.Get(context.Background(), "is_admin")
		require.NoError(t, err)
		assert.Equal(t, true, val2)

		val3, err := sess2.Get(context.Background(), "preferences")
		require.NoError(t, err)
		prefs, ok := val3.(map[string]string)
		assert.True(t, ok)
		assert.Equal(t, "dark", prefs["theme"])
	})

	// 测试错误处理
	t.Run("错误处理", func(t *testing.T) {
		// 重置模拟对象
		store = newMockStore()
		propagator = newMockPropagatorForIntegration()
		manager = &Manager{
			Store:      store,
			Propagator: propagator,
			SessCtxKey: "session",
		}

		// 创建HTTP请求和响应
		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		resp := httptest.NewRecorder()

		// 创建上下文
		ctx := ant.Context{
			Req:        req,
			Resp:       resp,
			UserValues: make(map[string]any),
		}

		// 设置模拟错误
		propagator.mockPropagator.extractErr = true
		store.getErr = true

		// 尝试获取不存在的会话
		_, err := manager.GetSession(ctx)
		assert.Error(t, err)

		// 尝试刷新不存在的会话
		_, err = manager.RefreshSession(ctx)
		assert.Error(t, err)

		// 尝试删除不存在的会话
		err = manager.RemoveSession(ctx)
		assert.Error(t, err)

		// 重置错误标志
		propagator.mockPropagator.extractErr = false
		store.getErr = false

		// 尝试获取会话中不存在的键
		sess, err := manager.InitSession(ctx, "error-session-id")
		require.NoError(t, err)

		_, err = sess.Get(context.Background(), "non_existent_key")
		assert.Error(t, err)
	})
}

// TestSessionWithMiddleware 测试会话管理器与中间件的集成
func TestSessionWithMiddleware(t *testing.T) {
	// 创建模拟的存储和传播器
	store := newMockStore()
	propagator := newMockPropagatorForIntegration()

	// 创建会话管理器
	manager := &Manager{
		Store:      store,
		Propagator: propagator,
		SessCtxKey: "session",
	}

	// 创建一个简单的处理函数，用于测试会话中间件
	handler := func(ctx ant.Context) error {
		// 获取会话
		sess, err := manager.GetSession(ctx)
		if err != nil {
			// 如果没有会话，创建一个新的
			sess, err = manager.InitSession(ctx, "middleware-session-id")
			if err != nil {
				return err
			}
		}

		// 设置会话数据
		err = sess.Set(ctx.Req.Context(), "visited", true)
		if err != nil {
			return err
		}

		// 返回会话ID
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte(sess.ID())

		// 写入响应
		_, err = ctx.Resp.Write(ctx.RespData)
		if err != nil {
			return err
		}

		return nil
	}

	// 创建HTTP请求和响应
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	resp := httptest.NewRecorder()

	// 创建上下文
	ctx := ant.Context{
		Req:        req,
		Resp:       resp,
		UserValues: make(map[string]any),
	}

	// 调用处理函数
	err := handler(ctx)
	require.NoError(t, err)

	// 验证响应
	assert.Equal(t, "middleware-session-id", resp.Body.String())

	// 验证会话已设置
	assert.True(t, propagator.sessions["middleware-session-id"])

	// 创建第二个请求，模拟带上Cookie
	req2 := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	propagator.req = req2 // 设置请求，使Extract能够返回会话ID
	resp2 := httptest.NewRecorder()

	ctx2 := ant.Context{
		Req:        req2,
		Resp:       resp2,
		UserValues: make(map[string]any),
	}

	// 再次调用处理函数
	err = handler(ctx2)
	require.NoError(t, err)

	// 验证响应
	assert.Equal(t, "middleware-session-id", resp2.Body.String())

	// 验证会话数据
	sess, err := manager.GetSession(ctx2)
	require.NoError(t, err)

	val, err := sess.Get(context.Background(), "visited")
	require.NoError(t, err)
	assert.Equal(t, true, val)
}
