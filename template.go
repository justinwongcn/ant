package ant

import (
	"bytes"
	"context"
	"html/template"
	"io/fs"
)

// TemplateEngine 定义了模板引擎的接口
// 提供模板渲染的核心功能
type TemplateEngine interface {
	// Render 渲染页面
	// ctx: 上下文对象，用于控制渲染过程
	// tplName: 模板名称，用于指定要渲染的模板
	// data: 渲染页面所需的数据
	// 返回值:
	// - []byte: 渲染后的页面内容
	// - error: 渲染过程中发生的错误
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
}

// GoTemplateEngine 是基于Go标准库html/template的模板引擎实现
// 提供模板加载和渲染功能
type GoTemplateEngine struct {
	// T 是底层的模板对象
	// 使用单个Template实例而不是map，因为Template本身支持按名称索引子模板
	T *template.Template
}

// Render 实现了TemplateEngine接口
// 将数据渲染到指定模板中，并返回渲染结果
// 保持RespData语义，支持中间件对渲染结果进行修改
//
// ctx: 上下文对象，用于控制渲染过程
// tplName: 要渲染的模板名称
// data: 渲染所需的数据
// 返回值:
// - []byte: 渲染后的页面内容
// - error: 渲染过程中的错误
func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	res := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(res, tplName, data)
	return res.Bytes(), err
}

// LoadFromGlob 从指定的glob模式加载模板
// pattern: glob模式字符串，用于匹配模板文件
// 返回值: 加载过程中发生的错误
func (g *GoTemplateEngine) LoadFromGlob(pattern string) error {
	var err error
	g.T, err = template.ParseGlob(pattern)
	return err
}

// LoadFromFiles 从指定的文件列表加载模板
// files: 模板文件路径列表
// 返回值: 加载过程中发生的错误
func (g *GoTemplateEngine) LoadFromFiles(files ...string) error {
	var err error
	g.T, err = template.ParseFiles(files...)
	return err
}

// LoadFromFS 从文件系统加载模板
// fs: 文件系统接口，用于访问模板文件
// paths: 模板文件在文件系统中的路径列表
// 返回值: 加载过程中发生的错误
func (g *GoTemplateEngine) LoadFromFS(fs fs.FS, paths ...string) error {
	var err error
	g.T, err = template.ParseFS(fs, paths...)
	return err
}