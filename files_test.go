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

// TestFileUploader 测试文件上传功能
func TestFileUploader(t *testing.T) {
	// 创建临时目录用于存储上传的文件
	tmpDir, err := os.MkdirTemp("", "upload-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name           string
		fileContent    string
		filename       string
		expectedStatus int
		expectedBody   string
		fileNameFunc   func(string) string
	}{
		{
			name:           "成功上传文件",
			fileContent:    "test content",
			filename:       "test.txt",
			expectedStatus: http.StatusOK,
			expectedBody:   "上传成功，文件大小: 12 bytes",
		},
		{
			name:           "空文件",
			fileContent:    "",
			filename:       "empty.txt",
			expectedStatus: http.StatusOK,
			expectedBody:   "上传成功，文件大小: 0 bytes",
		},
		{
			name:           "使用自定义文件名",
			fileContent:    "test content",
			filename:       "test.txt",
			expectedStatus: http.StatusOK,
			expectedBody:   "上传成功，文件大小: 12 bytes",
			fileNameFunc: func(name string) string {
				return "custom_" + name
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建multipart表单
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("file", tt.filename)
			if err != nil {
				t.Fatal(err)
			}
			_, err = io.Copy(part, strings.NewReader(tt.fileContent))
			if err != nil {
				t.Fatal(err)
			}
			writer.Close()

			// 创建请求
			req := httptest.NewRequest(http.MethodPost, "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()

			// 创建上传处理器
			uploader := &FileUploader{
				FileField: "file",
				DstPathFunc: func(fh *multipart.FileHeader) string {
					fileName := fh.Filename
					if tt.fileNameFunc != nil {
						fileName = tt.fileNameFunc(fileName)
					}
					return filepath.Join(tmpDir, fileName)
				},
			}

			// 处理请求
			ctx := &Context{
				Req:  req,
				Resp: rec,
			}
			uploader.Handle()(ctx)

			// 验证响应
			if rec.Code != tt.expectedStatus {
				t.Errorf("期望状态码 %d, 得到 %d", tt.expectedStatus, rec.Code)
			}

			if string(ctx.RespData) != tt.expectedBody {
				t.Errorf("期望响应体 %s, 得到 %s", tt.expectedBody, string(ctx.RespData))
			}

			// 验证文件是否被正确保存
			savedContent, err := os.ReadFile(filepath.Join(tmpDir, tt.filename))
			if err != nil {
				t.Fatal(err)
			}
			if string(savedContent) != tt.fileContent {
				t.Errorf("期望文件内容 %s, 得到 %s", tt.fileContent, string(savedContent))
			}
		})
	}
}

// TestFileDownloader 测试文件下载功能
func TestFileDownloader(t *testing.T) {
	// 创建临时目录和测试文件
	tmpDir, err := os.MkdirTemp("", "download-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建测试文件
	testContent := "test content"
	testFilePath := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFilePath, []byte(testContent), 0o666); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		fileName       string
		expectedStatus int
		expectedBody   string
		expectedHeader http.Header
	}{
		{
			name:           "成功下载文件",
			fileName:       "test.txt",
			expectedStatus: http.StatusOK,
			expectedBody:   testContent,
			expectedHeader: http.Header{
				"Content-Disposition": []string{`attachment; filename="test.txt"`},
				"Content-Type":        []string{"application/octet-stream"},
				"Content-Length":      []string{fmt.Sprintf("%d", len(testContent))},
			},
		},
		{
			name:           "文件不存在",
			fileName:       "nonexistent.txt",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "文件不存在",
		},
		{
			name:           "非法路径",
			fileName:       "../test.txt",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "非法的文件路径",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建请求
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/download?file=%s", tt.fileName), nil)
			rec := httptest.NewRecorder()

			// 创建下载处理器
			downloader := &FileDownloader{Dir: tmpDir}

			// 处理请求
			ctx := &Context{
				Req:  req,
				Resp: rec,
			}
			downloader.Handle()(ctx)

			// 验证响应状态码
			if ctx.RespStatusCode != tt.expectedStatus {
				t.Errorf("期望状态码 %d, 得到 %d", tt.expectedStatus, ctx.RespStatusCode)
			}

			// 对于成功的情况，验证响应头和内容
			if tt.expectedStatus == http.StatusOK {
				// 验证响应头
				for key, values := range tt.expectedHeader {
					actualValues := rec.Header()[key]
					if len(actualValues) != len(values) {
						t.Errorf("期望响应头 %s 的值数量为 %d, 得到 %d", key, len(values), len(actualValues))
						continue
					}
					for i, value := range values {
						if actualValues[i] != value {
							t.Errorf("期望响应头 %s 的值为 %s, 得到 %s", key, value, actualValues[i])
						}
					}
				}

				// 验证响应体
				if rec.Body.String() != tt.expectedBody {
					t.Errorf("期望响应体 %s, 得到 %s", tt.expectedBody, rec.Body.String())
				}
			} else {
				// 对于错误情况，验证错误消息
				if string(ctx.RespData) != tt.expectedBody {
					t.Errorf("期望错误消息 %s, 得到 %s", tt.expectedBody, string(ctx.RespData))
				}
			}
		})
	}
}

// TestStaticResourceHandler 测试静态资源处理器
func TestStaticResourceHandler(t *testing.T) {
	// 创建临时目录和测试文件
	tmpDir, err := os.MkdirTemp("", "static-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建测试文件
	testFiles := map[string]struct {
		content     string
		contentType string
	}{
		"test.html": {content: "<html>test</html>", contentType: "text/html; charset=utf-8"},
		"test.css":  {content: "body { color: red; }", contentType: "text/css; charset=utf-8"},
		"test.js":   {content: "console.log('test');", contentType: "application/javascript"},
		"test.txt":  {content: "plain text", contentType: "text/plain"},
	}

	for filename, file := range testFiles {
		if err := os.WriteFile(filepath.Join(tmpDir, filename), []byte(file.content), 0o666); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name           string
		fileName       string
		expectedStatus int
		expectedBody   string
		expectedType   string
	}{
		{
			name:           "HTML文件",
			fileName:       "test.html",
			expectedStatus: http.StatusOK,
			expectedBody:   "<html>test</html>",
			expectedType:   "text/html; charset=utf-8",
		},
		{
			name:           "CSS文件",
			fileName:       "test.css",
			expectedStatus: http.StatusOK,
			expectedBody:   "body { color: red; }",
			expectedType:   "text/css; charset=utf-8",
		},
		{
			name:           "文件不存在",
			fileName:       "nonexistent.txt",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "打开文件失败",
		},
		{
			name:           "非法路径",
			fileName:       "../test.txt",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "未指定文件名",
		},
	}

	// 创建处理器
	handler := NewStaticResourceHandler(tmpDir, "/static/", WithFileCache(1024*1024, 100))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建请求
			req := httptest.NewRequest(http.MethodGet, "/static/"+tt.fileName, nil)
			rec := httptest.NewRecorder()

			// 创建Context实例
			ctx := &Context{
				Req:  req,
				Resp: rec,
			}

			// 设置请求路径
			req.SetPathValue("file", tt.fileName)

			// 处理请求
			handler.Handle(ctx)

			// 验证响应状态码
			if ctx.RespStatusCode != tt.expectedStatus {
				t.Errorf("期望状态码 %d, 得到 %d", tt.expectedStatus, ctx.RespStatusCode)
			}

			// 对于成功的情况，验证响应头和内容
			if tt.expectedStatus == http.StatusOK {
				// 验证Content-Type
				actualType := rec.Header().Get("Content-Type")
				if actualType != tt.expectedType {
					t.Errorf("期望Content-Type %s, 得到 %s", tt.expectedType, actualType)
				}

				// 验证响应体
				if rec.Body.String() != tt.expectedBody {
					t.Errorf("期望响应体 %s, 得到 %s", tt.expectedBody, rec.Body.String())
				}

				// 验证缓存相关的响应头
				if rec.Header().Get("Cache-Control") != "public, max-age=31536000" {
					t.Error("缓存控制头设置不正确")
				}
				if rec.Header().Get("Last-Modified") == "" {
					t.Error("未设置Last-Modified头")
				}
			} else {
				// 对于错误情况，验证错误消息
				if string(ctx.RespData) != tt.expectedBody {
					t.Errorf("期望错误消息 %s, 得到 %s", tt.expectedBody, string(ctx.RespData))
				}
			}
		})
	}
}

// 在 Windows 环境下，因权限问题测试用例无法通过

func TestFileUploaderError(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name           string
		setupFunc      func(string) // 用于设置测试环境
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "创建目标目录失败",
			setupFunc: func(dir string) {
				// 创建一个只读文件，阻止创建目录
				// 使用一个不存在的路径，确保MkdirAll会失败
				os.WriteFile(filepath.Join(dir, "readonly"), []byte("test"), 0o400)
				// 修改tmpDir的权限为只读，确保无法在其中创建新目录
				os.Chmod(dir, 0o400)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "创建目录失败",
		},
		{
			name: "创建目标文件失败",
			setupFunc: func(dir string) {
				// 创建目标子目录，但设置为只读，确保无法在其中创建文件
				subdir := filepath.Join(dir, "subdir")
				os.MkdirAll(subdir, 0o755)
				os.Chmod(subdir, 0o400)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "创建文件失败",
		},
	}

	// 遍历测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建临时目录
			tmpDir, err := os.MkdirTemp("", "upload-error-test-*")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			// 设置测试环境
			if tt.setupFunc != nil {
				tt.setupFunc(tmpDir)
			}

			// 创建multipart表单
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("file", "test.txt")
			if err != nil {
				t.Fatal(err)
			}
			_, err = io.Copy(part, strings.NewReader("test content"))
			if err != nil {
				t.Fatal(err)
			}
			writer.Close()

			// 创建请求
			req := httptest.NewRequest(http.MethodPost, "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()

			// 创建上传处理器
			uploader := &FileUploader{
				FileField: "file",
				DstPathFunc: func(fh *multipart.FileHeader) string {
					return filepath.Join(tmpDir, "subdir", fh.Filename)
				},
			}

			// 处理请求
			ctx := &Context{
				Req:  req,
				Resp: rec,
			}
			uploader.Handle()(ctx)

			// 验证响应
			if ctx.RespStatusCode != tt.expectedStatus {
				t.Errorf("期望状态码 %d, 得到 %d", tt.expectedStatus, ctx.RespStatusCode)
			}

			if string(ctx.RespData) != tt.expectedBody {
				t.Errorf("期望响应体 %s, 得到 %s", tt.expectedBody, string(ctx.RespData))
			}
		})
	}
}

func TestFileUploaderMoreErrors(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func(*http.Request, *FileUploader) // 用于设置测试环境
		expectedStatus int
		expectedBody   string
		fileNameFunc   func(string) string
	}{
		{
			name: "上传失败，未找到文件",
			setupFunc: func(req *http.Request, uploader *FileUploader) {
				// 使用错误的字段名，导致找不到文件
				uploader.FileField = "wrong_field"
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "上传失败，未找到文件",
		},
		{
			name: "保存文件失败",
			setupFunc: func(req *http.Request, uploader *FileUploader) {
				// 使用一个不存在的目录路径，这在所有操作系统上都会失败
				uploader.DstPathFunc = func(fh *multipart.FileHeader) string {
					// 使用一个不可能存在的深层嵌套路径
					return filepath.Join("/nonexistent/directory/that/cannot/be/created", fh.Filename)
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "创建目录失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建临时目录
			tmpDir, err := os.MkdirTemp("", "upload-more-errors-test-*")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			// 创建multipart表单
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("file", "test.txt")
			if err != nil {
				t.Fatal(err)
			}
			_, err = io.Copy(part, strings.NewReader("test content"))
			if err != nil {
				t.Fatal(err)
			}
			writer.Close()

			// 创建请求
			req := httptest.NewRequest(http.MethodPost, "/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()

			// 创建上传处理器
			uploader := &FileUploader{
				FileField: "file",
				DstPathFunc: func(fh *multipart.FileHeader) string {
					fileName := fh.Filename
					if tt.fileNameFunc != nil {
						fileName = tt.fileNameFunc(fileName)
					}
					return filepath.Join(tmpDir, fileName)
				},
			}

			// 应用测试特定的设置
			if tt.setupFunc != nil {
				tt.setupFunc(req, uploader)
			}

			// 处理请求
			ctx := &Context{
				Req:  req,
				Resp: rec,
			}
			uploader.Handle()(ctx)

			// 验证响应
			if ctx.RespStatusCode != tt.expectedStatus {
				t.Errorf("期望状态码 %d, 得到 %d", tt.expectedStatus, ctx.RespStatusCode)
			}

			if string(ctx.RespData) != tt.expectedBody {
				t.Errorf("期望响应体 %s, 得到 %s", tt.expectedBody, string(ctx.RespData))
			}
		})
	}
}

func TestFileDownloaderMoreErrors(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func(*http.Request, *FileDownloader, string) // 用于设置测试环境
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "未指定文件名",
			setupFunc: func(req *http.Request, downloader *FileDownloader, tmpDir string) {
				// 不添加file查询参数
				*req = *httptest.NewRequest(http.MethodGet, "/download", nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "未指定文件名",
		},
		{
			name: "读取文件失败",
			setupFunc: func(req *http.Request, downloader *FileDownloader, tmpDir string) {
				// 创建一个无法读取的文件（只有写权限）
				filePath := filepath.Join(tmpDir, "unreadable.txt")
				os.WriteFile(filePath, []byte("test content"), 0o200) // 只有写权限
				*req = *httptest.NewRequest(http.MethodGet, "/download?file=unreadable.txt", nil)
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   "读取文件失败",
		},
		{
			name: "打开文件失败",
			setupFunc: func(req *http.Request, downloader *FileDownloader, tmpDir string) {
				// 创建一个目录而不是文件，尝试打开它会失败
				dirPath := filepath.Join(tmpDir, "directory")
				os.Mkdir(dirPath, 0o755)
				*req = *httptest.NewRequest(http.MethodGet, "/download?file=directory", nil)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "打开文件失败",
		},
		// 注意："发送文件失败"的情况很难在测试中模拟，因为它发生在响应已经开始发送之后
		// 这通常需要模拟网络错误或者响应写入器的错误，这在测试环境中比较复杂
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建临时目录
			tmpDir, err := os.MkdirTemp("", "download-more-errors-test-*")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			// 创建一个普通测试文件
			testFilePath := filepath.Join(tmpDir, "test.txt")
			if err := os.WriteFile(testFilePath, []byte("test content"), 0o666); err != nil {
				t.Fatal(err)
			}

			// 创建请求
			req := httptest.NewRequest(http.MethodGet, "/download?file=test.txt", nil)
			rec := httptest.NewRecorder()

			// 创建下载处理器
			downloader := &FileDownloader{Dir: tmpDir}

			// 应用测试特定的设置
			if tt.setupFunc != nil {
				tt.setupFunc(req, downloader, tmpDir)
			}

			// 处理请求
			ctx := &Context{
				Req:  req,
				Resp: rec,
			}
			downloader.Handle()(ctx)

			// 验证响应
			if ctx.RespStatusCode != tt.expectedStatus {
				t.Errorf("期望状态码 %d, 得到 %d", tt.expectedStatus, ctx.RespStatusCode)
			}

			if string(ctx.RespData) != tt.expectedBody {
				t.Errorf("期望响应体 %s, 得到 %s", tt.expectedBody, string(ctx.RespData))
			}
		})
	}
}

func TestWithMoreExtension(t *testing.T) {
	tests := []struct {
		name     string
		extMap   map[string]string
		expected map[string]string
	}{
		{
			name: "添加新的扩展名映射",
			extMap: map[string]string{
				"mp3": "audio/mpeg",
				"mp4": "video/mp4",
			},
			expected: map[string]string{
				"html": "text/html; charset=utf-8",
				"mp3":  "audio/mpeg",
				"mp4":  "video/mp4",
			},
		},
		{
			name: "覆盖已存在的扩展名映射",
			extMap: map[string]string{
				"html": "text/html",
				"txt":  "text/markdown",
			},
			expected: map[string]string{
				"html": "text/html",
				"txt":  "text/markdown",
			},
		},
		{
			name:   "添加空映射",
			extMap: map[string]string{},
			expected: map[string]string{
				"html": "text/html; charset=utf-8",
				"txt":  "text/plain",
			},
		},
		{
			name: "添加多个扩展名映射",
			extMap: map[string]string{
				"mp3":  "audio/mpeg",
				"wav":  "audio/wav",
				"mp4":  "video/mp4",
				"webm": "video/webm",
			},
			expected: map[string]string{
				"html": "text/html; charset=utf-8",
				"mp3":  "audio/mpeg",
				"wav":  "audio/wav",
				"mp4":  "video/mp4",
				"webm": "video/webm",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建一个基础的StaticResourceHandler
			h := NewStaticResourceHandler("", "")

			// 应用WithMoreExtension选项
			opt := WithMoreExtension(tt.extMap)
			opt(h)

			// 验证扩展名映射
			for ext, contentType := range tt.expected {
				actual, ok := h.extensionContentTypeMap[ext]
				if !ok {
					t.Errorf("扩展名 %s 的映射不存在", ext)
					continue
				}
				if actual != contentType {
					t.Errorf("扩展名 %s 的Content-Type期望为 %s，实际为 %s", ext, contentType, actual)
				}
			}
		})
	}
}
