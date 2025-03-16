package cookie

import "net/http"

// Propagator 基于Cookie的会话传播器
// 用于在HTTP请求和响应中传递会话信息
type Propagator struct {
	// cookieName 存储会话ID的Cookie名称
	cookieName   string
	// cookieOption 用于配置Cookie属性的函数
	cookieOption func(cookie *http.Cookie)
}

// NewPropagator 创建一个新的Cookie传播器实例
// 参数:
// - opts: 可选的配置函数列表，用于自定义Propagator的行为
// 返回值:
// - *Propagator: 配置完成的Cookie传播器实例
func NewPropagator(opts ...func(*Propagator)) *Propagator {
	p := &Propagator{
		cookieName:   "sessid",
		cookieOption: func(c *http.Cookie) {},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// WithCookieName 设置Cookie名称的选项
// 参数:
// - name: 自定义的Cookie名称
// 返回值:
// - func(*Propagator): 返回一个配置函数，用于设置Cookie名称
func WithCookieName(name string) func(*Propagator) {
	return func(p *Propagator) {
		p.cookieName = name
	}
}

// WithCookieOption 设置Cookie选项的函数
// 参数:
// - fn: 用于配置http.Cookie属性的函数
// 返回值:
// - func(*Propagator): 返回一个配置函数，用于设置Cookie选项
func WithCookieOption(fn func(cookie *http.Cookie)) func(*Propagator) {
	return func(p *Propagator) {
		p.cookieOption = fn
	}
}

// Inject 将会话ID注入到HTTP响应的Cookie中
// 参数:
// - id: 要注入的会话ID
// - writer: HTTP响应写入器
// 返回值:
// - error: 注入过程中可能发生的错误
func (p *Propagator) Inject(id string, writer http.ResponseWriter) error {
	c := &http.Cookie{
		Name:  p.cookieName,
		Value: id,
	}
	p.cookieOption(c)
	http.SetCookie(writer, c)

	return nil
}

// Extract 从HTTP请求的Cookie中提取会话ID
// 参数:
// - req: HTTP请求
// 返回值:
// - string: 提取的会话ID
// - error: 提取过程中可能发生的错误，如Cookie不存在
func (p *Propagator) Extract(req *http.Request) (string, error) {
	c, err := req.Cookie(p.cookieName)
	if err != nil {
		return "", err
	}

	return c.Value, nil
}

// Remove 从HTTP响应中移除会话Cookie
// 参数:
// - writer: HTTP响应写入器
// 返回值:
// - error: 移除过程中可能发生的错误
func (p *Propagator) Remove(writer http.ResponseWriter) error {
	c := &http.Cookie{
		Name:   p.cookieName,
		MaxAge: -1,
	}
	http.SetCookie(writer, c)

	return nil
}
