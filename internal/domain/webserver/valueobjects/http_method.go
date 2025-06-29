// Package valueobjects 包含web服务器领域的值对象
package valueobjects

import (
	"fmt"
	"strings"

	"github.com/justinwongcn/ant/internal/domain/shared/errors"
)

// HTTPMethod 表示HTTP方法值对象
type HTTPMethod struct {
	value string
}

// 有效的HTTP方法
const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	TRACE   = "TRACE"
	CONNECT = "CONNECT"
)

var validMethods = map[string]bool{
	GET:     true,
	POST:    true,
	PUT:     true,
	DELETE:  true,
	PATCH:   true,
	HEAD:    true,
	OPTIONS: true,
	TRACE:   true,
	CONNECT: true,
}

// NewHTTPMethod 创建一个新的HTTPMethod值对象
func NewHTTPMethod(method string) (*HTTPMethod, error) {
	method = strings.ToUpper(strings.TrimSpace(method))

	if method == "" {
		return nil, errors.NewValidationError("method", "HTTP方法不能为空")
	}

	if !validMethods[method] {
		return nil, errors.NewValidationError("method", fmt.Sprintf("无效的HTTP方法: %s", method))
	}

	return &HTTPMethod{value: method}, nil
}

// MustNewHTTPMethod 创建一个新的HTTPMethod值对象，出错时panic
func MustNewHTTPMethod(method string) *HTTPMethod {
	m, err := NewHTTPMethod(method)
	if err != nil {
		panic(err)
	}
	return m
}

// Value 返回HTTP方法的字符串值
func (h *HTTPMethod) Value() string {
	return h.value
}

// String 返回字符串表示
func (h *HTTPMethod) String() string {
	return h.value
}

// Equals 检查两个HTTPMethod对象是否相等
func (h *HTTPMethod) Equals(other *HTTPMethod) bool {
	if other == nil {
		return false
	}
	return h.value == other.value
}

// IsIdempotent 如果HTTP方法是幂等的则返回true
func (h *HTTPMethod) IsIdempotent() bool {
	switch h.value {
	case GET, PUT, DELETE, HEAD, OPTIONS, TRACE:
		return true
	default:
		return false
	}
}

// IsSafe 如果HTTP方法是安全的(只读)则返回true
func (h *HTTPMethod) IsSafe() bool {
	switch h.value {
	case GET, HEAD, OPTIONS, TRACE:
		return true
	default:
		return false
	}
}

// AllowsBody 如果HTTP方法允许请求体则返回true
func (h *HTTPMethod) AllowsBody() bool {
	switch h.value {
	case POST, PUT, PATCH:
		return true
	default:
		return false
	}
}

// 预定义的HTTP方法实例
var (
	HTTPMethodGET     = MustNewHTTPMethod(GET)
	HTTPMethodPOST    = MustNewHTTPMethod(POST)
	HTTPMethodPUT     = MustNewHTTPMethod(PUT)
	HTTPMethodDELETE  = MustNewHTTPMethod(DELETE)
	HTTPMethodPATCH   = MustNewHTTPMethod(PATCH)
	HTTPMethodHEAD    = MustNewHTTPMethod(HEAD)
	HTTPMethodOPTIONS = MustNewHTTPMethod(OPTIONS)
	HTTPMethodTRACE   = MustNewHTTPMethod(TRACE)
	HTTPMethodCONNECT = MustNewHTTPMethod(CONNECT)
)

// GetAllHTTPMethods 返回所有有效的HTTP方法
func GetAllHTTPMethods() []*HTTPMethod {
	return []*HTTPMethod{
		HTTPMethodGET,
		HTTPMethodPOST,
		HTTPMethodPUT,
		HTTPMethodDELETE,
		HTTPMethodPATCH,
		HTTPMethodHEAD,
		HTTPMethodOPTIONS,
		HTTPMethodTRACE,
		HTTPMethodCONNECT,
	}
}

// IsValidHTTPMethod 检查字符串是否是有效的HTTP方法
func IsValidHTTPMethod(method string) bool {
	method = strings.ToUpper(strings.TrimSpace(method))
	return validMethods[method]
}
