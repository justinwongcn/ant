package ant

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// Context 封装HTTP请求上下文，提供请求处理相关工具方法
// 包含原始请求/响应对象，缓存数据及响应状态信息
type Context struct {
	Req  *http.Request       // 原始HTTP请求对象
	Resp http.ResponseWriter // HTTP响应写入器

	// 缓存解析后的URL查询参数，避免重复解析
	cacheQueryValues url.Values

	// 响应缓存数据，在最终响应时一次性写入
	RespStatusCode int    // 响应状态码
	RespData       []byte // 响应内容主体

	// 模板引擎，用于渲染HTML模板
	TemplateEngine TemplateEngine

	// 用户相关的数据，用于在请求处理过程中存储临时数据
	UserValues map[string]any
}

// BindJSON 解析请求体中的JSON数据并绑定到指定结构体
// val: 需要绑定的目标结构体指针
// 返回值: 解析成功返回nil，失败返回对应错误
// 注意事项：当请求体为空时返回特定错误
func (c *Context) BindJSON(val any) error {
	if c.Req.Body == nil {
		return errors.New("web: body 为 nil")
	}
	decoder := json.NewDecoder(c.Req.Body)
	decoder.DisallowUnknownFields() // 禁止未知字段
	return decoder.Decode(val)
}

// StringValue 封装字符串值与解析错误的组合结构
// 提供类型转换方法，统一处理转换错误
type StringValue struct {
	val string
	err error
}

// String 获取原始字符串值及可能存在的错误
func (s StringValue) String() (string, error) {
	return s.val, s.err
}

// ToInt64 将字符串值转换为int64类型
// 返回值: 转换成功返回整数值，失败返回错误（包含原始错误或转换错误）
func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

// FormValue 从POST表单中获取指定键的值
// key: 表单字段名称
// 返回值: 封装后的字符串值结构，包含值或错误信息
// 注意：会自动解析表单内容，仅返回第一个值
func (c *Context) FormValue(key string) StringValue {
	if err := c.Req.ParseForm(); err != nil {
		return StringValue{err: err}
	}
	value := c.Req.FormValue(key)
	if value == "" {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: value}
}

// PostFormValue 仅从 POST 表单数据中获取指定键的值
// key: 表单字段名称
// 返回值: 封装后的字符串值结构，包含值或错误信息
func (c *Context) PostFormValue(key string) StringValue {
	if err := c.Req.ParseForm(); err != nil {
		return StringValue{err: err}
	}
	value := c.Req.PostFormValue(key)
	if value == "" {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: value}
}

// QueryValue 从 URL 查询参数中获取指定 key 的值
//
// key 参数指定要获取的查询参数名称
//
// 返回 StringValue 类型，包含查询参数的值或错误信息
func (c *Context) QueryValue(key string) StringValue {
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}

	value := c.cacheQueryValues.Get(key)
	if value == "" {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}

	return StringValue{val: value}
}

// DefaultQueryValue 从 URL 查询参数中获取指定 key 的值，如果不存在则返回默认值
//
// key 参数指定要获取的查询参数名称
// defaultValue 参数指定当参数不存在时的默认值
//
// 返回 StringValue 类型，包含查询参数的值或默认值
func (c *Context) DefaultQueryValue(key string, defaultValue string) StringValue {
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}

	value := c.cacheQueryValues.Get(key)
	if value == "" {
		return StringValue{val: defaultValue}
	}

	return StringValue{val: value}
}

// PathValue 从 URL 路径参数中获取指定 key 的值
//
// key 参数指定要获取的路径参数名称
//
// 返回 StringValue 类型，包含路径参数的值或错误信息
func (c *Context) PathValue(key string) StringValue {
	value := c.Req.PathValue(key)
	if value == "" {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}

	return StringValue{val: value}
}

// SetCookie 设置HTTP Cookie到响应中
// cookie: 需要设置的cookie对象指针
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Resp, cookie)
}

// RespJSON 将数据序列化为JSON格式响应
// code: HTTP状态码
// val: 需要序列化的数据结构
// 返回值: 序列化或写入响应时发生的错误
func (c *Context) RespJSON(code int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Resp.WriteHeader(code)
	_, err = c.Resp.Write(bs)
	return err
}

// RespJSONOK 发送状态码200的JSON成功响应
// val: 需要序列化的数据结构
// 返回值: 同RespJSON方法
func (c *Context) RespJSONOK(val any) error {
	return c.RespJSON(http.StatusOK, val)
}

// RespTemplate 渲染HTML模板并将结果写入响应
// tplName: 模板名称
// data: 渲染数据
// 返回值: 渲染过程中的错误
func (c *Context) RespTemplate(tplName string, data any) error {
	if c.TemplateEngine == nil {
		return errors.New("web: 未设置模板引擎")
	}

	// 渲染模板
	bs, err := c.TemplateEngine.Render(context.Background(), tplName, data)
	if err != nil {
		return err
	}

	// 设置状态码和响应数据
	c.RespStatusCode = http.StatusOK
	c.RespData = bs

	// 设置Content-Type
	c.Resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Resp.WriteHeader(c.RespStatusCode)

	// 直接写入响应体，确保在测试中也能正确写入
	_, err = c.Resp.Write(bs)
	return err
}

// WriteString 将字符串写入响应体
func (c *Context) WriteString(data string) error {
	c.RespData = []byte(data)
	return nil
}
