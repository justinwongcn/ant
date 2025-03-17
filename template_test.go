package ant

import (
	"context"
	"html/template"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

// TestGoTemplateEngineRender 测试基本的模板渲染功能
func TestGoTemplateEngineRender(t *testing.T) {
	tests := []struct {
		name     string
		tpl      string
		data     any
		expected string
		wantErr  bool
	}{
		{
			name:     "基本渲染",
			tpl:      "Hello {{.Name}}",
			data:     struct{ Name string }{Name: "World"},
			expected: "Hello World",
			wantErr:  false,
		},
		{
			name:     "条件渲染",
			tpl:      "{{if .Show}}显示{{else}}隐藏{{end}}",
			data:     struct{ Show bool }{Show: true},
			expected: "显示",
			wantErr:  false,
		},
		{
			name:     "语法错误",
			tpl:      "{{.Name}", // 缺少结束标记
			data:     struct{ Name string }{Name: "World"},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := &GoTemplateEngine{}
			// 解析模板字符串
			var err error
			engine.T, err = template.New("test").Parse(tt.tpl)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("解析模板失败: %v", err)
				}
				return
			}

			// 渲染模板
			result, err := engine.Render(context.Background(), "test", tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && string(result) != tt.expected {
				t.Errorf("Render() = %v, want %v", string(result), tt.expected)
			}
		})
	}
}

// TestGoTemplateEngineLoadFromFiles 测试从文件加载模板
func TestGoTemplateEngineLoadFromFiles(t *testing.T) {
	// 创建临时模板文件
	tmpFile, err := os.CreateTemp("", "template-*.html")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入测试模板内容
	tplContent := "Hello {{.Name}}"
	if _, err := tmpFile.WriteString(tplContent); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// 测试加载和渲染
	engine := &GoTemplateEngine{}
	if err := engine.LoadFromFiles(tmpFile.Name()); err != nil {
		t.Fatalf("LoadFromFiles() error = %v", err)
	}

	// 渲染并验证结果
	// 使用文件的基本名称作为模板名称
	tplName := filepath.Base(tmpFile.Name())
	result, err := engine.Render(context.Background(), tplName, struct{ Name string }{Name: "World"})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	expected := "Hello World"
	if string(result) != expected {
		t.Errorf("Render() = %v, want %v", string(result), expected)
	}
}

// TestGoTemplateEngineLoadFromFS 测试从文件系统加载模板
func TestGoTemplateEngineLoadFromFS(t *testing.T) {
	// 创建内存文件系统
	fs := fstest.MapFS{
		"test.html": &fstest.MapFile{
			Data: []byte("Hello {{.Name}}"),
		},
	}

	// 测试加载
	engine := &GoTemplateEngine{}
	if err := engine.LoadFromFS(fs, "test.html"); err != nil {
		t.Fatalf("LoadFromFS() error = %v", err)
	}

	// 渲染并验证结果
	result, err := engine.Render(context.Background(), "test.html", struct{ Name string }{Name: "World"})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	expected := "Hello World"
	if string(result) != expected {
		t.Errorf("Render() = %v, want %v", string(result), expected)
	}
}

// TestGoTemplateEngineLoadFromGlob 测试使用glob模式加载模板
func TestGoTemplateEngineLoadFromGlob(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "templates-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建测试模板文件
	tplPath := tmpDir + "/test.html"
	if err := os.WriteFile(tplPath, []byte("Hello {{.Name}}"), 0o666); err != nil {
		t.Fatal(err)
	}

	// 测试加载
	engine := &GoTemplateEngine{}
	if err := engine.LoadFromGlob(tmpDir + "/*.html"); err != nil {
		t.Fatalf("LoadFromGlob() error = %v", err)
	}

	// 渲染并验证结果
	result, err := engine.Render(context.Background(), "test.html", struct{ Name string }{Name: "World"})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	expected := "Hello World"
	if string(result) != expected {
		t.Errorf("Render() = %v, want %v", string(result), expected)
	}
}
