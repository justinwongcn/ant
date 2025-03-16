package cookie

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewPropagator(t *testing.T) {
	// 测试默认配置
	p1 := NewPropagator()
	if p1.cookieName != "sessid" {
		t.Errorf("默认 cookieName 应为 'sessid'，实际为 '%s'", p1.cookieName)
	}

	// 测试自定义配置
	p2 := NewPropagator(WithCookieName("custom_session"))
	if p2.cookieName != "custom_session" {
		t.Errorf("自定义 cookieName 应为 'custom_session'，实际为 '%s'", p2.cookieName)
	}

	// 测试多个配置选项
	cookieOptionCalled := false
	cookieOption := func(c *http.Cookie) {
		cookieOptionCalled = true
		c.MaxAge = 3600
	}

	p3 := NewPropagator(
		WithCookieName("test_session"),
		WithCookieOption(cookieOption),
	)

	if p3.cookieName != "test_session" {
		t.Errorf("自定义 cookieName 应为 'test_session'，实际为 '%s'", p3.cookieName)
	}

	// 测试 cookieOption 是否被正确设置
	c := &http.Cookie{}
	p3.cookieOption(c)
	if !cookieOptionCalled {
		t.Error("cookieOption 未被调用")
	}
	if c.MaxAge != 3600 {
		t.Errorf("Cookie MaxAge 应为 3600，实际为 %d", c.MaxAge)
	}
}

func TestPropagatorInject(t *testing.T) {
	testCases := []struct {
		name       string
		sessionID  string
		cookieName string
		cookieOpt  func(*http.Cookie)
	}{
		{
			name:       "基本注入",
			sessionID:  "test-session-123",
			cookieName: "sessid",
			cookieOpt:  func(c *http.Cookie) {},
		},
		{
			name:       "自定义Cookie名称",
			sessionID:  "custom-session-456",
			cookieName: "custom_sessid",
			cookieOpt:  func(c *http.Cookie) {},
		},
		{
			name:       "带Cookie选项",
			sessionID:  "secure-session-789",
			cookieName: "sessid",
			cookieOpt: func(c *http.Cookie) {
				c.MaxAge = 3600
				c.Path = "/"
				c.HttpOnly = true
				c.Secure = true
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPropagator(
				WithCookieName(tc.cookieName),
				WithCookieOption(tc.cookieOpt),
			)

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 注入会话ID
			err := p.Inject(tc.sessionID, w)
			if err != nil {
				t.Fatalf("Inject() 返回错误: %v", err)
			}

			// 检查响应头中的Cookie
			cookies := w.Result().Cookies()
			if len(cookies) != 1 {
				t.Fatalf("预期有1个Cookie，实际有 %d 个", len(cookies))
			}

			cookie := cookies[0]
			if cookie.Name != tc.cookieName {
				t.Errorf("Cookie名称应为 '%s'，实际为 '%s'", tc.cookieName, cookie.Name)
			}
			if cookie.Value != tc.sessionID {
				t.Errorf("Cookie值应为 '%s'，实际为 '%s'", tc.sessionID, cookie.Value)
			}

			// 如果是带选项的测试用例，检查选项是否生效
			if tc.name == "带Cookie选项" {
				if cookie.MaxAge != 3600 {
					t.Errorf("Cookie MaxAge 应为 3600，实际为 %d", cookie.MaxAge)
				}
				if cookie.Path != "/" {
					t.Errorf("Cookie Path 应为 '/'，实际为 '%s'", cookie.Path)
				}
				if !cookie.HttpOnly {
					t.Error("Cookie HttpOnly 应为 true")
				}
				if !cookie.Secure {
					t.Error("Cookie Secure 应为 true")
				}
			}
		})
	}
}

func TestPropagatorExtract(t *testing.T) {
	testCases := []struct {
		name          string
		cookieName    string
		sessionID     string
		hasCookie     bool
		expectedError bool
	}{
		{
			name:          "成功提取",
			cookieName:    "sessid",
			sessionID:     "test-session-123",
			hasCookie:     true,
			expectedError: false,
		},
		{
			name:          "自定义Cookie名称",
			cookieName:    "custom_sessid",
			sessionID:     "custom-session-456",
			hasCookie:     true,
			expectedError: false,
		},
		{
			name:          "Cookie不存在",
			cookieName:    "sessid",
			sessionID:     "",
			hasCookie:     false,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPropagator(WithCookieName(tc.cookieName))

			// 创建请求
			req := httptest.NewRequest("GET", "http://example.com", nil)

			// 如果测试用例需要Cookie，则添加
			if tc.hasCookie {
				req.AddCookie(&http.Cookie{
					Name:  tc.cookieName,
					Value: tc.sessionID,
				})
			}

			// 提取会话ID
			id, err := p.Extract(req)

			// 检查错误
			if tc.expectedError && err == nil {
				t.Error("预期会返回错误，但没有")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("预期不会返回错误，但返回了: %v", err)
			}

			// 如果不期望错误，检查提取的ID
			if !tc.expectedError {
				if id != tc.sessionID {
					t.Errorf("提取的会话ID应为 '%s'，实际为 '%s'", tc.sessionID, id)
				}
			}
		})
	}
}

func TestPropagatorRemove(t *testing.T) {
	testCases := []struct {
		name       string
		cookieName string
	}{
		{
			name:       "默认Cookie名称",
			cookieName: "sessid",
		},
		{
			name:       "自定义Cookie名称",
			cookieName: "custom_sessid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPropagator(WithCookieName(tc.cookieName))

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 移除Cookie
			err := p.Remove(w)
			if err != nil {
				t.Fatalf("Remove() 返回错误: %v", err)
			}

			// 检查响应头中的Cookie
			cookies := w.Result().Cookies()
			if len(cookies) != 1 {
				t.Fatalf("预期有1个Cookie，实际有 %d 个", len(cookies))
			}

			cookie := cookies[0]
			if cookie.Name != tc.cookieName {
				t.Errorf("Cookie名称应为 '%s'，实际为 '%s'", tc.cookieName, cookie.Name)
			}
			if cookie.MaxAge != -1 {
				t.Errorf("Cookie MaxAge 应为 -1，实际为 %d", cookie.MaxAge)
			}
		})
	}
}

func TestWithCookieOption(t *testing.T) {
	// 测试各种Cookie选项
	testCases := []struct {
		name     string
		option   func(*http.Cookie)
		validate func(*testing.T, *http.Cookie)
	}{
		{
			name: "设置过期时间",
			option: func(c *http.Cookie) {
				expiration := time.Now().Add(24 * time.Hour)
				c.Expires = expiration
			},
			validate: func(t *testing.T, c *http.Cookie) {
				if c.Expires.IsZero() {
					t.Error("Cookie Expires 未被设置")
				}
			},
		},
		{
			name: "设置安全选项",
			option: func(c *http.Cookie) {
				c.Secure = true
				c.HttpOnly = true
				c.SameSite = http.SameSiteStrictMode
			},
			validate: func(t *testing.T, c *http.Cookie) {
				if !c.Secure {
					t.Error("Cookie Secure 应为 true")
				}
				if !c.HttpOnly {
					t.Error("Cookie HttpOnly 应为 true")
				}
				if c.SameSite != http.SameSiteStrictMode {
					t.Error("Cookie SameSite 应为 SameSiteStrictMode")
				}
			},
		},
		{
			name: "设置路径和域",
			option: func(c *http.Cookie) {
				c.Path = "/api"
				c.Domain = "example.com"
			},
			validate: func(t *testing.T, c *http.Cookie) {
				if c.Path != "/api" {
					t.Errorf("Cookie Path 应为 '/api'，实际为 '%s'", c.Path)
				}
				if c.Domain != "example.com" {
					t.Errorf("Cookie Domain 应为 'example.com'，实际为 '%s'", c.Domain)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPropagator(WithCookieOption(tc.option))

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 注入会话ID以触发Cookie选项
			err := p.Inject("test-session", w)
			if err != nil {
				t.Fatalf("Inject() 返回错误: %v", err)
			}

			// 检查响应头中的Cookie
			cookies := w.Result().Cookies()
			if len(cookies) != 1 {
				t.Fatalf("预期有1个Cookie，实际有 %d 个", len(cookies))
			}

			// 验证Cookie选项
			tc.validate(t, cookies[0])
		})
	}
}
