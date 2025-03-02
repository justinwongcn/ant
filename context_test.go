package ant

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContextBindJSON(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		reqBody  io.Reader
		wantErr  error
		wantData TestData
	}{
		{
			name: "正常解析",
			body: `{"name":"test","age":18}`,
			wantData: TestData{
				Name: "test",
				Age:  18,
			},
		},
		{
			name:    "body 为空",
			body:    "",
			wantErr: errors.New("EOF"),
		},
		{
			name:    "body为nil",
			reqBody: nil,
			wantErr: errors.New("web: body 为 nil"),
		},
		{
			name:    "JSON 格式错误",
			body:    `{"name":"test","age":}`,
			wantErr: &json.SyntaxError{},
		},
		{
			name:    "字段类型不匹配",
			body:    `{"name":"test","age":"invalid"}`,
			wantErr: &json.UnmarshalTypeError{},
		},
		{
			name:    "未知字段",
			body:    `{"name":"test","age":18,"unknown":"field"}`,
			wantErr: errors.New("json: unknown field \"unknown\""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var bodyReader io.Reader
			if tt.reqBody != nil {
				bodyReader = tt.reqBody
			} else {
				bodyReader = bytes.NewBufferString(tt.body)
			}

			req, err := http.NewRequest(http.MethodPost, "/test", bodyReader)
			assert.NoError(t, err)

			if tt.reqBody == nil && tt.name == "body为nil" {
				req.Body = nil
			} else {
				req.Body = io.NopCloser(bodyReader)
			}

			ctx := &Context{Req: req}

			var data TestData
			err = ctx.BindJSON(&data)

			if tt.wantErr != nil {
				assert.Error(t, err)
				switch wantErr := tt.wantErr.(type) {
				case *json.SyntaxError:
					assert.ErrorAs(t, err, &wantErr)
				case *json.UnmarshalTypeError:
					assert.ErrorAs(t, err, &wantErr)
				default:
					assert.ErrorContains(t, err, tt.wantErr.Error())
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantData, data)
		})
	}
}

func TestContextFormValue(t *testing.T) {
	tests := []struct {
		name      string
		form      map[string]string
		key       string
		wantValue string
		wantErr   bool
	}{
		{
			name: "正常获取",
			form: map[string]string{
				"name": "test",
			},
			key:       "name",
			wantValue: "test",
		},
		{
			name:    "key不存在",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/test", nil)
			req.Form = make(url.Values)
			for k, v := range tt.form {
				req.Form.Add(k, v)
			}

			ctx := &Context{Req: req}
			val := ctx.FormValue(tt.key)
			actualValue, err := val.String()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantValue, actualValue)
		})
	}
}

func TestContextQueryValue(t *testing.T) {
	tests := []struct {
		name      string
		query     map[string]string
		key       string
		wantValue string
		wantErr   bool
	}{
		{
			name: "正常获取",
			query: map[string]string{
				"id": "123",
			},
			key:       "id",
			wantValue: "123",
		},
		{
			name:    "key不存在",
			key:     "notexist",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/test" // 重命名变量避免与包名冲突
			if len(tt.query) > 0 {
				vals := make(url.Values)
				for k, v := range tt.query {
					vals.Add(k, v) // 改用Add方法支持多值参数
				}
				path += "?" + vals.Encode()
			}

			req := httptest.NewRequest(http.MethodGet, path, nil)
			ctx := &Context{Req: req}
			val := ctx.QueryValue(tt.key)
			actualValue, err := val.String()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantValue, actualValue)
		})
	}
}

func TestContextPathValue(t *testing.T) {
	tests := []struct {
		name      string
		pathParam string // 路径参数名称
		paramVal  string // 路径参数值
		key       string // 要获取的key
		wantValue string // 期望值
		wantErr   bool   // 是否期望错误
	}{
		{
			name:      "正常获取",
			pathParam: "id",
			paramVal:  "123",
			key:       "id",
			wantValue: "123",
		},
		{
			name:      "key不存在",
			pathParam: "name", // 设置存在的参数
			paramVal:  "test",
			key:       "notexist", // 查询不存在的key
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 构造带路径参数的请求
			req := httptest.NewRequest("GET", "/test", nil)

			// 模拟设置路径参数（假设使用标准库的PathValue实现）
			req.SetPathValue(tt.pathParam, tt.paramVal)

			ctx := &Context{Req: req} // 初始化带有请求的上下文
			val := ctx.PathValue(tt.key)
			actualValue, err := val.String()

			if tt.wantErr {
				assert.ErrorContains(t, err, "找不到这个 key")
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantValue, actualValue)
		})
	}
}

func TestContextRespJSON(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		data     interface{}
		wantResp string
		wantErr  bool
	}{
		{
			name: "正常响应",
			code: http.StatusOK,
			data: TestData{
				Name: "test",
				Age:  18,
			},
			wantResp: `{"name":"test","age":18}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx := &Context{Resp: w}

			err := ctx.RespJSON(tt.code, tt.data)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.code, w.Code)
			assert.JSONEq(t, tt.wantResp, w.Body.String())
		})
	}
}

func TestContextRespJSONOK(t *testing.T) {
	// 定义无法序列化的测试类型
	type Unserializable struct {
		Channel chan int `json:"channel"`
	}

	tests := []struct {
		name     string
		data     interface{}
		wantCode int
		wantResp string
		wantErr  bool
	}{
		{
			name:     "正常响应",
			data:     TestData{Name: "test", Age: 20},
			wantCode: http.StatusOK,
			wantResp: `{"name":"test","age":20}`,
		},
		{
			name:    "无法序列化的数据",
			data:    Unserializable{Channel: make(chan int)},
			wantErr: true,
		},
		{
			name:    "空数据响应",
			data:    nil,
			wantErr: false, // 空值可以正常序列化为null
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx := &Context{Resp: w}

			err := ctx.RespJSONOK(tt.data)

			if tt.wantErr {
				assert.Error(t, err, "预期返回错误")
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code, "状态码必须为200")

			if tt.wantResp != "" {
				assert.JSONEq(t, tt.wantResp, w.Body.String())
			} else {
				assert.Equal(t, "null", w.Body.String(), "空值响应应为null")
			}
		})
	}
}

func TestStringValueToInt64(t *testing.T) {
	tests := []struct {
		name      string
		val       string
		err       error
		wantValue int64
		wantErr   string
	}{
		{
			name:      "正常转换",
			val:       "123",
			wantValue: 123,
		},
		{
			name:    "空字符串",
			val:     "",
			wantErr: "invalid syntax",
		},
		{
			name:    "非数字字符",
			val:     "abc",
			wantErr: "invalid syntax",
		},
		{
			name:    "超出int64范围",
			val:     "9223372036854775808", // 超过int64最大值(9223372036854775807)
			wantErr: "value out of range",
		},
		{
			name:    "携带原始错误",
			err:     errors.New("原始错误"),
			wantErr: "原始错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv := StringValue{
				val: tt.val,
				err: tt.err,
			}

			result, err := sv.ToInt64()

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantValue, result)
		})
	}
}

func TestContextSetCookie(t *testing.T) {
	tests := []struct {
		name       string
		cookie     *http.Cookie
		wantHeader string
	}{
		{
			name: "基本Cookie设置",
			cookie: &http.Cookie{
				Name:  "session",
				Value: "abc123",
			},
			wantHeader: "session=abc123",
		},
		{
			name: "带过期时间",
			cookie: &http.Cookie{
				Name:    "auth",
				Value:   "token",
				Expires: time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			},
			wantHeader: "Expires=Sun, 31 Dec 2023 00:00:00 GMT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			ctx := &Context{Resp: w}

			ctx.SetCookie(tt.cookie)

			cookies := w.Result().Cookies()
			assert.Len(t, cookies, 1, "应该设置一个Cookie")

			actual := cookies[0]
			assert.Equal(t, tt.cookie.Name, actual.Name)
			assert.Equal(t, tt.cookie.Value, actual.Value)
			assert.Equal(t, tt.cookie.Domain, actual.Domain)
			assert.Equal(t, tt.cookie.Path, actual.Path)
			assert.Equal(t, tt.cookie.Secure, actual.Secure)
			assert.Equal(t, tt.cookie.HttpOnly, actual.HttpOnly)

			setCookieHeader := w.Header().Get("Set-Cookie")
			assert.Contains(t, setCookieHeader, tt.wantHeader)
		})
	}
}

func TestContextPostFormValue(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		key      string
		wantVal  string
		wantErr  bool
		errMsg   string
	}{
		{
			name: "获取存在的表单值",
			body: "name=test",
			key: "name",
			wantVal: "test",
			wantErr: false,
		},
		{
			name: "获取不存在的表单值",
			body: "name=test",
			key: "age",
			wantErr: true,
			errMsg: "web: 找不到这个 key",
		},
		{
			name: "表单为空",
			body: "",
			key: "name",
			wantErr: true,
			errMsg: "web: 找不到这个 key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/test",
				strings.NewReader(tt.body),
			)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			ctx := &Context{
				Req: req,
			}

			val := ctx.PostFormValue(tt.key)
			str, err := val.String()

			if tt.wantErr {
				if err == nil {
					t.Fatal("期望得到错误，但是没有得到错误")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("期望得到的错误信息是 %s，但是得到的是 %s", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if str != tt.wantVal {
				t.Errorf("期望获取到的值是 %s，但是得到的是 %s", tt.wantVal, str)
			}
		})
	}
}

// TestData 用于测试的数据结构
type TestData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// MockTemplateEngine 用于测试的模拟模板引擎
type MockTemplateEngine struct {
	shouldError bool
	output      []byte
}

func (m *MockTemplateEngine) Render(_ context.Context, _ string, _ any) ([]byte, error) {
	if m.shouldError {
		return nil, errors.New("模板渲染失败")
	}
	return m.output, nil
}

func TestContextRespTemplate(t *testing.T) {
	tests := []struct {
		name        string
		engine      TemplateEngine
		expectedErr  string
		expectedResp string
	}{
		{
			name:        "模板引擎未设置",
			engine:      nil,
			expectedErr: "web: 未设置模板引擎",
		},
		{
			name: "模板渲染失败",
			engine: &MockTemplateEngine{
				shouldError: true,
			},
			expectedErr: "模板渲染失败",
		},
		{
			name: "成功渲染模板",
			engine: &MockTemplateEngine{
				output: []byte("<h1>Hello World</h1>"),
			},
			expectedResp: "<h1>Hello World</h1>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试上下文
			w := httptest.NewRecorder()
			ctx := &Context{
				Resp:           w,
				TemplateEngine: tt.engine,
			}

			// 调用被测试的方法
			err := ctx.RespTemplate("test.html", nil)

			// 验证错误
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
				return
			}

			// 验证成功情况
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResp, w.Body.String())
			
			// 验证Content-Type头部
			contentType := w.Header().Get("Content-Type")
			assert.Equal(t, "text/html; charset=utf-8", contentType)
		})
	}
}