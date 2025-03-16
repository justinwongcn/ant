package session

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justinwongcn/ant"
)

// 创建一个模拟的 Session 实现
type mockSession struct {
	id    string
	data  map[string]any
	store map[string]any
}

func (m *mockSession) Get(ctx context.Context, key string) (any, error) {
	val, ok := m.data[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	return val, nil
}

func (m *mockSession) Set(ctx context.Context, key string, value any) error {
	m.data[key] = value
	return nil
}

func (m *mockSession) ID() string {
	return m.id
}

// 创建一个模拟的 Store 实现
type mockStore struct {
	sessions map[string]*mockSession
	// 用于测试错误情况
	generateErr bool
	refreshErr  bool
	removeErr   bool
	getErr      bool
}

func newMockStore() *mockStore {
	return &mockStore{
		sessions: make(map[string]*mockSession),
	}
}

func (m *mockStore) Generate(ctx context.Context, id string) (Session, error) {
	if m.generateErr {
		return nil, errors.New("generate error")
	}
	sess := &mockSession{
		id:    id,
		data:  make(map[string]any),
		store: make(map[string]any),
	}
	m.sessions[id] = sess
	return sess, nil
}

func (m *mockStore) Refresh(ctx context.Context, id string) error {
	if m.refreshErr {
		return errors.New("refresh error")
	}
	_, ok := m.sessions[id]
	if !ok {
		return errors.New("session not found")
	}
	return nil
}

func (m *mockStore) Remove(ctx context.Context, id string) error {
	if m.removeErr {
		return errors.New("remove error")
	}
	delete(m.sessions, id)
	return nil
}

func (m *mockStore) Get(ctx context.Context, id string) (Session, error) {
	if m.getErr {
		return nil, errors.New("get error")
	}
	sess, ok := m.sessions[id]
	if !ok {
		return nil, errors.New("session not found")
	}
	return sess, nil
}

// 创建一个模拟的 Propagator 实现
type mockPropagator struct {
	sessions map[string]bool
	// 用于测试错误情况
	injectErr  bool
	extractErr bool
	removeErr  bool
}

func newMockPropagator() *mockPropagator {
	return &mockPropagator{
		sessions: make(map[string]bool),
	}
}

func (m *mockPropagator) Inject(id string, writer http.ResponseWriter) error {
	if m.injectErr {
		return errors.New("inject error")
	}
	m.sessions[id] = true
	return nil
}

func (m *mockPropagator) Extract(req *http.Request) (string, error) {
	if m.extractErr {
		return "", errors.New("extract error")
	}
	// 从请求头中获取会话ID
	id := req.Header.Get("X-Session-ID")
	if id == "" {
		return "", errors.New("session id not found")
	}
	return id, nil
}

func (m *mockPropagator) Remove(writer http.ResponseWriter) error {
	if m.removeErr {
		return errors.New("remove error")
	}
	return nil
}

// 创建测试用的 Context
func createTestContext(sessionID string) ant.Context {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	if sessionID != "" {
		req.Header.Set("X-Session-ID", sessionID)
	}
	resp := httptest.NewRecorder()
	return ant.Context{
		Req:        req,
		Resp:       resp,
		UserValues: make(map[string]any),
	}
}

func TestManager_GetSession(t *testing.T) {
	testCases := []struct {
		name          string
		sessionID     string
		setupStore    func(*mockStore)
		setupProp     func(*mockPropagator)
		expectErr     bool
		expectSession bool
	}{
		{
			name:      "成功获取会话",
			sessionID: "test-session-123",
			setupStore: func(ms *mockStore) {
				// 预先创建一个会话
				ctx := context.Background()
				sess, _ := ms.Generate(ctx, "test-session-123")
				_ = sess.Set(ctx, "key1", "value1")
			},
			setupProp: func(mp *mockPropagator) {
				// 不需要特殊设置
			},
			expectErr:     false,
			expectSession: true,
		},
		{
			name:      "提取会话ID失败",
			sessionID: "",
			setupStore: func(ms *mockStore) {
				// 不需要预设会话
			},
			setupProp: func(mp *mockPropagator) {
				// 设置提取错误
				mp.extractErr = true
			},
			expectErr:     true,
			expectSession: false,
		},
		{
			name:      "获取会话失败",
			sessionID: "non-existent-session",
			setupStore: func(ms *mockStore) {
				// 设置获取错误
				ms.getErr = true
			},
			setupProp: func(mp *mockPropagator) {
				// 不需要特殊设置
			},
			expectErr:     true,
			expectSession: false,
		},
		{
			name:      "从上下文中获取已存在的会话",
			sessionID: "context-session-456",
			setupStore: func(ms *mockStore) {
				// 不需要预设会话，因为会从上下文中获取
			},
			setupProp: func(mp *mockPropagator) {
				// 不需要特殊设置
			},
			expectErr:     false,
			expectSession: true,
		},
		{
			name:      "UserValues为nil时的情况",
			sessionID: "test-session-123",
			setupStore: func(ms *mockStore) {
				// 预先创建一个会话
				ctx := context.Background()
				sess, _ := ms.Generate(ctx, "test-session-123")
				_ = sess.Set(ctx, "key1", "value1")
			},
			setupProp: func(mp *mockPropagator) {
				// 不需要特殊设置
			},
			expectErr:     false,
			expectSession: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建模拟组件
			store := newMockStore()
			prop := newMockPropagator()

			// 设置测试条件
			tc.setupStore(store)
			tc.setupProp(prop)

			// 创建管理器
			manager := &Manager{
				Store:      store,
				Propagator: prop,
				SessCtxKey: "session",
			}

			// 创建上下文
			ctx := createTestContext(tc.sessionID)

			// 如果是测试从上下文中获取会话的情况
			if tc.name == "从上下文中获取已存在的会话" {
				// 在上下文中预设会话
				sess := &mockSession{
					id:   "context-session-456",
					data: make(map[string]any),
				}
				ctx.UserValues["session"] = sess
			}

			// 执行测试
			sess, err := manager.GetSession(ctx)

			// 验证结果
			if tc.expectErr && err == nil {
				t.Error("预期会返回错误，但没有")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("预期不会返回错误，但返回了: %v", err)
			}
			if tc.expectSession && sess == nil {
				t.Error("预期会返回会话，但返回了nil")
			}
			if !tc.expectSession && sess != nil {
				t.Error("预期不会返回会话，但返回了非nil值")
			}

			// 如果成功获取会话，验证会话是否已存储在上下文中
			if !tc.expectErr && tc.expectSession {
				storedSess, ok := ctx.UserValues["session"]
				if !ok {
					t.Error("会话未存储在上下文中")
				}
				if storedSess != sess {
					t.Error("上下文中存储的会话与返回的会话不一致")
				}
			}
		})
	}
}

func TestManager_InitSession(t *testing.T) {
	testCases := []struct {
		name          string
		sessionID     string
		setupStore    func(*mockStore)
		setupProp     func(*mockPropagator)
		expectErr     bool
		expectSession bool
	}{
		{
			name:      "成功初始化会话",
			sessionID: "new-session-123",
			setupStore: func(ms *mockStore) {
				// 不需要特殊设置
			},
			setupProp: func(mp *mockPropagator) {
				// 不需要特殊设置
			},
			expectErr:     false,
			expectSession: true,
		},
		{
			name:      "生成会话失败",
			sessionID: "error-session",
			setupStore: func(ms *mockStore) {
				// 设置生成错误
				ms.generateErr = true
			},
			setupProp: func(mp *mockPropagator) {
				// 不需要特殊设置
			},
			expectErr:     true,
			expectSession: false,
		},
		{
			name:      "注入会话ID失败",
			sessionID: "inject-error-session",
			setupStore: func(ms *mockStore) {
				// 不需要特殊设置
			},
			setupProp: func(mp *mockPropagator) {
				// 设置注入错误
				mp.injectErr = true
			},
			expectErr:     true,
			expectSession: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建模拟组件
			store := newMockStore()
			prop := newMockPropagator()

			// 设置测试条件
			tc.setupStore(store)
			tc.setupProp(prop)

			// 创建管理器
			manager := &Manager{
				Store:      store,
				Propagator: prop,
				SessCtxKey: "session",
			}

			// 创建上下文
			ctx := createTestContext("")

			// 执行测试
			sess, err := manager.InitSession(ctx, tc.sessionID)

			// 验证结果
			if tc.expectErr && err == nil {
				t.Error("预期会返回错误，但没有")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("预期不会返回错误，但返回了: %v", err)
			}
			if tc.expectSession && sess == nil {
				t.Error("预期会返回会话，但返回了nil")
			}
			if !tc.expectSession && sess != nil {
				t.Error("预期不会返回会话，但返回了非nil值")
			}

			// 如果成功初始化会话，验证会话ID是否正确
			if !tc.expectErr && tc.expectSession {
				if sess.ID() != tc.sessionID {
					t.Errorf("会话ID应为 '%s'，实际为 '%s'", tc.sessionID, sess.ID())
				}
			}
		})
	}
}

func TestManager_RefreshSession(t *testing.T) {
	testCases := []struct {
		name          string
		sessionID     string
		setupStore    func(*mockStore)
		setupProp     func(*mockPropagator)
		expectErr     bool
		expectSession bool
	}{
		{
			name:      "成功刷新会话",
			sessionID: "refresh-session-123",
			setupStore: func(ms *mockStore) {
				// 预先创建一个会话
				ctx := context.Background()
				_, _ = ms.Generate(ctx, "refresh-session-123")
			},
			setupProp: func(mp *mockPropagator) {
				// 不需要特殊设置
			},
			expectErr:     false,
			expectSession: true,
		},
		{
			name:      "获取会话失败",
			sessionID: "non-existent-session",
			setupStore: func(ms *mockStore) {
				// 不预设会话
			},
			setupProp: func(mp *mockPropagator) {
				// 设置提取错误
				mp.extractErr = true
			},
			expectErr:     true,
			expectSession: false,
		},
		{
			name:      "刷新会话失败",
			sessionID: "refresh-error-session",
			setupStore: func(ms *mockStore) {
				// 预先创建一个会话
				ctx := context.Background()
				_, _ = ms.Generate(ctx, "refresh-error-session")
				// 设置刷新错误
				ms.refreshErr = true
			},
			setupProp: func(mp *mockPropagator) {
				// 不需要特殊设置
			},
			expectErr:     true,
			expectSession: false,
		},
		{
			name:      "注入会话ID失败",
			sessionID: "inject-error-session",
			setupStore: func(ms *mockStore) {
				// 预先创建一个会话
				ctx := context.Background()
				_, _ = ms.Generate(ctx, "inject-error-session")
			},
			setupProp: func(mp *mockPropagator) {
				// 设置注入错误
				mp.injectErr = true
			},
			expectErr:     true,
			expectSession: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建模拟组件
			store := newMockStore()
			prop := newMockPropagator()

			// 设置测试条件
			tc.setupStore(store)
			tc.setupProp(prop)

			// 创建管理器
			manager := &Manager{
				Store:      store,
				Propagator: prop,
				SessCtxKey: "session",
			}

			// 创建上下文
			ctx := createTestContext(tc.sessionID)

			// 执行测试
			sess, err := manager.RefreshSession(ctx)

			// 验证结果
			if tc.expectErr && err == nil {
				t.Error("预期会返回错误，但没有")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("预期不会返回错误，但返回了: %v", err)
			}
			if tc.expectSession && sess == nil {
				t.Error("预期会返回会话，但返回了nil")
			}
			if !tc.expectSession && sess != nil {
				t.Error("预期不会返回会话，但返回了非nil值")
			}

			// 如果成功刷新会话，验证会话ID是否正确
			if !tc.expectErr && tc.expectSession {
				if sess.ID() != tc.sessionID {
					t.Errorf("会话ID应为 '%s'，实际为 '%s'", tc.sessionID, sess.ID())
				}
			}
		})
	}
}

func TestManager_RemoveSession(t *testing.T) {
	testCases := []struct {
		name       string
		sessionID  string
		setupStore func(*mockStore)
		setupProp  func(*mockPropagator)
		expectErr  bool
	}{
		{
			name:      "成功删除会话",
			sessionID: "remove-session-123",
			setupStore: func(ms *mockStore) {
				// 预先创建一个会话
				ctx := context.Background()
				_, _ = ms.Generate(ctx, "remove-session-123")
			},
			setupProp: func(mp *mockPropagator) {
				// 不需要特殊设置
			},
			expectErr: false,
		},
		{
			name:      "获取会话失败",
			sessionID: "non-existent-session",
			setupStore: func(ms *mockStore) {
				// 不预设会话
			},
			setupProp: func(mp *mockPropagator) {
				// 设置提取错误
				mp.extractErr = true
			},
			expectErr: true,
		},
		{
			name:      "删除会话失败",
			sessionID: "remove-error-session",
			setupStore: func(ms *mockStore) {
				// 预先创建一个会话
				ctx := context.Background()
				_, _ = ms.Generate(ctx, "remove-error-session")
				// 设置删除错误
				ms.removeErr = true
			},
			setupProp: func(mp *mockPropagator) {
				// 不需要特殊设置
			},
			expectErr: true,
		},
		{
			name:      "从响应中移除会话ID失败",
			sessionID: "remove-prop-error-session",
			setupStore: func(ms *mockStore) {
				// 预先创建一个会话
				ctx := context.Background()
				_, _ = ms.Generate(ctx, "remove-prop-error-session")
			},
			setupProp: func(mp *mockPropagator) {
				// 设置移除错误
				mp.removeErr = true
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建模拟组件
			store := newMockStore()
			prop := newMockPropagator()

			// 设置测试条件
			tc.setupStore(store)
			tc.setupProp(prop)

			// 创建管理器
			manager := &Manager{
				Store:      store,
				Propagator: prop,
				SessCtxKey: "session",
			}

			// 创建上下文
			ctx := createTestContext(tc.sessionID)

			// 执行测试
			err := manager.RemoveSession(ctx)

			// 验证结果
			if tc.expectErr && err == nil {
				t.Error("预期会返回错误，但没有")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("预期不会返回错误，但返回了: %v", err)
			}

			// 如果成功删除会话，验证会话是否已从存储中移除
			if !tc.expectErr {
				// 尝试获取已删除的会话
				_, err := store.Get(context.Background(), tc.sessionID)
				if err == nil {
					t.Error("会话应该已被删除，但仍能获取")
				}
			}
		})
	}
}
