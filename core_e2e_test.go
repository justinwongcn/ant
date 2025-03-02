package ant

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TestCoreComponentsE2E 测试核心组件（Context、Template、Server）的端到端集成
func TestCoreComponentsE2E(t *testing.T) {
	// 创建临时目录用于存放模板文件
	tmpDir, err := os.MkdirTemp("", "core-e2e-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建测试用的模板文件
	tplContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>{{.Title}}</title>
	</head>
	<body>
		<h1>{{.Title}}</h1>
		<p>{{.Content}}</p>
	</body>
	</html>
	`
	tplPath := filepath.Join(tmpDir, "test.html")
	if err := os.WriteFile(tplPath, []byte(tplContent), 0666); err != nil {
		t.Fatal(err)
	}

	// 创建并配置模板引擎
	tplEngine := &GoTemplateEngine{}
	if err := tplEngine.LoadFromFiles(tplPath); err != nil {
		t.Fatal(err)
	}

	// 创建HTTP服务器并配置模板引擎
	server := NewHTTPServer(ServerWithTemplateEngine(tplEngine))

	// 测试JSON绑定和响应
	t.Run("JSON处理", func(t *testing.T) {
		type TestData struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		// 注册处理JSON请求的路由
		server.Handle("POST /json", func(ctx *Context) {
			var data TestData
			if err := ctx.BindJSON(&data); err != nil {
				ctx.RespStatusCode = http.StatusBadRequest
				ctx.RespData = []byte(err.Error())
				return
			}
			// 增加年龄并返回
			data.Age++
			ctx.RespJSON(http.StatusOK, data)
		})

		// 创建测试请求
		inputData := TestData{Name: "张三", Age: 18}
		body, _ := json.Marshal(inputData)
		req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// 发送请求
		server.ServeHTTP(rec, req)

		// 验证响应
		if rec.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, rec.Code)
		}

		var respData TestData
		if err := json.Unmarshal(rec.Body.Bytes(), &respData); err != nil {
			t.Fatal(err)
		}

		if respData.Name != inputData.Name {
			t.Errorf("期望名字 %s, 得到 %s", inputData.Name, respData.Name)
		}
		if respData.Age != inputData.Age+1 {
			t.Errorf("期望年龄 %d, 得到 %d", inputData.Age+1, respData.Age)
		}
	})

	// 测试查询参数和路径参数
	t.Run("参数处理", func(t *testing.T) {
		// 注册处理参数的路由
		server.Handle("GET /users/{id}", func(ctx *Context) {
			id, err := ctx.PathValue("id").String()
			if err != nil {
				ctx.RespStatusCode = http.StatusBadRequest
				ctx.RespData = []byte("无效的用户ID")
				return
			}

			age, err := ctx.QueryValue("age").ToInt64()
			if err != nil {
				ctx.RespStatusCode = http.StatusBadRequest
				ctx.RespData = []byte("无效的年龄参数")
				return
			}

			response := map[string]interface{}{
				"id":  id,
				"age": age,
			}
			ctx.RespJSONOK(response)
		})

		// 创建测试请求
		req := httptest.NewRequest(http.MethodGet, "/users/123?age=25", nil)
		req.SetPathValue("id", "123")
		rec := httptest.NewRecorder()

		// 发送请求
		server.ServeHTTP(rec, req)

		// 验证响应
		if rec.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, rec.Code)
		}

		var respData map[string]interface{}
		if err := json.Unmarshal(rec.Body.Bytes(), &respData); err != nil {
			t.Fatal(err)
		}

		if respData["id"] != "123" {
			t.Errorf("期望ID %s, 得到 %v", "123", respData["id"])
		}
		if respData["age"] != float64(25) {
			t.Errorf("期望年龄 %d, 得到 %v", 25, respData["age"])
		}
	})

	// 测试模板渲染
	t.Run("模板渲染", func(t *testing.T) {
		// 注册处理模板渲染的路由
		server.Handle("GET /page", func(ctx *Context) {
			data := map[string]string{
				"Title":   "测试页面",
				"Content": "这是一个测试内容",
			}
			if err := ctx.RespTemplate("test.html", data); err != nil {
				ctx.RespStatusCode = http.StatusInternalServerError
				ctx.RespData = []byte(err.Error())
				return
			}
		})

		// 创建测试请求
		req := httptest.NewRequest(http.MethodGet, "/page", nil)
		rec := httptest.NewRecorder()

		// 发送请求
		server.ServeHTTP(rec, req)

		// 验证响应
		if rec.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, rec.Code)
		}

		if contentType := rec.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
			t.Errorf("期望Content-Type %s, 得到 %s", "text/html; charset=utf-8", contentType)
		}

		expectedContent := `<!DOCTYPE html>`
		if !bytes.Contains(rec.Body.Bytes(), []byte(expectedContent)) {
			t.Errorf("响应内容中未找到期望的HTML标记")
		}
	})

	// 测试错误处理
	t.Run("错误处理", func(t *testing.T) {
		// 测试无效的JSON绑定
		t.Run("无效的JSON", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewReader([]byte(`{"name":"张三","age":"invalid"}`)))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("期望状态码 %d, 得到 %d", http.StatusBadRequest, rec.Code)
			}
		})

		// 测试无效的查询参数
		t.Run("无效的查询参数", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/users/123?age=invalid", nil)
			req.SetPathValue("id", "123")
			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("期望状态码 %d, 得到 %d", http.StatusBadRequest, rec.Code)
			}
		})

		// 测试模板不存在
		t.Run("模板不存在", func(t *testing.T) {
			// 注册一个使用不存在模板的处理器
			server.Handle("GET /invalid-template", func(ctx *Context) {
				err := ctx.RespTemplate("nonexistent.html", nil)
				if err != nil {
					ctx.RespStatusCode = http.StatusInternalServerError
					ctx.RespData = []byte(err.Error())
				}
			})

			req := httptest.NewRequest(http.MethodGet, "/invalid-template", nil)
			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, req)

			if rec.Code != http.StatusInternalServerError {
				t.Errorf("期望状态码 %d, 得到 %d", http.StatusInternalServerError, rec.Code)
			}
		})
	})
}