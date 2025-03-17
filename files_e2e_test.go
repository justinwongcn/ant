package ant

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFileHandlingE2E 测试文件处理的端到端流程
func TestFileHandlingE2E(t *testing.T) {
	// 创建临时目录用于测试
	tmpDir, err := os.MkdirTemp("", "file-handling-e2e-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建HTTP服务器
	server := NewHTTPServer()

	// 配置文件上传处理器
	uploader := &FileUploader{
		FileField: "file",
		DstPathFunc: func(fh *multipart.FileHeader) string {
			return filepath.Join(tmpDir, fh.Filename)
		},
	}
	server.Handle("POST /upload", uploader.Handle())

	// 配置文件下载处理器
	downloader := &FileDownloader{Dir: tmpDir}
	server.Handle("GET /download", downloader.Handle())

	// 配置静态资源处理器
	static := NewStaticResourceHandler(tmpDir, "/static/", WithFileCache(1024*1024, 100))
	server.Handle("GET /static/{file}", static.Handle)

	// 测试文件上传
	t.Run("文件上传和下载流程", func(t *testing.T) {
		// 准备上传文件内容
		fileContent := "这是一个测试文件内容"
		fileName := "test.txt"

		// 创建multipart表单
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", fileName)
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(part, strings.NewReader(fileContent))
		if err != nil {
			t.Fatal(err)
		}
		writer.Close()

		// 发送上传请求
		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)

		// 验证上传响应
		if rec.Code != http.StatusOK {
			t.Errorf("文件上传失败，状态码：%d", rec.Code)
		}

		// 验证文件是否被正确保存
		savedContent, err := os.ReadFile(filepath.Join(tmpDir, fileName))
		if err != nil {
			t.Fatal(err)
		}
		if string(savedContent) != fileContent {
			t.Error("保存的文件内容与上传的内容不匹配")
		}

		// 测试文件下载
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/download?file=%s", fileName), nil)
		rec = httptest.NewRecorder()
		server.ServeHTTP(rec, req)

		// 验证下载响应
		if rec.Code != http.StatusOK {
			t.Errorf("文件下载失败，状态码：%d", rec.Code)
		}
		if rec.Body.String() != fileContent {
			t.Error("下载的文件内容与原始内容不匹配")
		}
	})

	// 测试静态资源服务
	t.Run("静态资源服务", func(t *testing.T) {
		// 创建测试HTML文件
		htmlContent := "<html><body>测试页面</body></html>"
		htmlFile := filepath.Join(tmpDir, "test.html")
		if err := os.WriteFile(htmlFile, []byte(htmlContent), 0o666); err != nil {
			t.Fatal(err)
		}

		// 发送静态资源请求
		req := httptest.NewRequest(http.MethodGet, "/static/test.html", nil)
		req.SetPathValue("file", "test.html")
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)

		// 验证响应
		if rec.Code != http.StatusOK {
			t.Errorf("静态资源请求失败，状态码：%d", rec.Code)
		}
		if rec.Body.String() != htmlContent {
			t.Error("静态资源内容与原始内容不匹配")
		}
		if rec.Header().Get("Content-Type") != "text/html; charset=utf-8" {
			t.Error("Content-Type不正确")
		}
	})

	// 测试错误处理
	t.Run("错误处理", func(t *testing.T) {
		// 测试下载不存在的文件
		req := httptest.NewRequest(http.MethodGet, "/download?file=nonexistent.txt", nil)
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("期望状态码404，得到：%d", rec.Code)
		}

		// 测试上传时未提供文件
		req = httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(""))
		req.Header.Set("Content-Type", "multipart/form-data")
		rec = httptest.NewRecorder()
		server.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("期望状态码400，得到：%d", rec.Code)
		}

		// 测试访问不存在的静态资源
		req = httptest.NewRequest(http.MethodGet, "/static/nonexistent.html", nil)
		req.SetPathValue("file", "nonexistent.html")
		rec = httptest.NewRecorder()
		server.ServeHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("期望状态码500，得到：%d", rec.Code)
		}
	})
}
