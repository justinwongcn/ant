// Package valueobjects 包含web服务器领域的值对象
package valueobjects

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/justinwongcn/ant/internal/domain/shared/errors"
)

// URLPattern 表示URL路由模式值对象
type URLPattern struct {
	method *HTTPMethod
	path   string
	regex  *regexp.Regexp
}

// NewURLPattern 创建一个新的URLPattern值对象
func NewURLPattern(method *HTTPMethod, path string) (*URLPattern, error) {
	if method == nil {
		return nil, errors.NewValidationError("method", "HTTP方法不能为空")
	}

	path = strings.TrimSpace(path)
	if path == "" {
		return nil, errors.NewValidationError("path", "路径不能为空")
	}

	// 确保路径以/开头
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// 验证路径格式
	if err := validatePath(path); err != nil {
		return nil, err
	}

	// 编译正则表达式用于模式匹配
	regex, err := compilePathRegex(path)
	if err != nil {
		return nil, errors.NewValidationError("path", fmt.Sprintf("无效的路径模式: %v", err))
	}

	return &URLPattern{
		method: method,
		path:   path,
		regex:  regex,
	}, nil
}

// MustNewURLPattern 创建一个新的URLPattern值对象，出错时panic
func MustNewURLPattern(method *HTTPMethod, path string) *URLPattern {
	pattern, err := NewURLPattern(method, path)
	if err != nil {
		panic(err)
	}
	return pattern
}

// Method 返回HTTP方法
func (u *URLPattern) Method() *HTTPMethod {
	return u.method
}

// Path 返回路径模式
func (u *URLPattern) Path() string {
	return u.path
}

// String 返回字符串表示(方法 + 路径)
func (u *URLPattern) String() string {
	return fmt.Sprintf("%s %s", u.method.Value(), u.path)
}

// Equals 检查两个URLPattern对象是否相等
func (u *URLPattern) Equals(other *URLPattern) bool {
	if other == nil {
		return false
	}
	return u.method.Equals(other.method) && u.path == other.path
}

// Matches 检查模式是否匹配给定的方法和路径
func (u *URLPattern) Matches(method *HTTPMethod, path string) bool {
	if !u.method.Equals(method) {
		return false
	}
	return u.regex.MatchString(path)
}

// ExtractParameters 从给定路径提取路径参数
func (u *URLPattern) ExtractParameters(path string) map[string]string {
	if !u.regex.MatchString(path) {
		return nil
	}

	matches := u.regex.FindStringSubmatch(path)
	if len(matches) <= 1 {
		return make(map[string]string)
	}

	params := make(map[string]string)
	names := u.regex.SubexpNames()

	for i, match := range matches[1:] {
		if names[i+1] != "" {
			params[names[i+1]] = match
		}
	}

	return params
}

// HasParameters 如果模式包含参数则返回true
func (u *URLPattern) HasParameters() bool {
	return strings.Contains(u.path, "{") && strings.Contains(u.path, "}")
}

// IsWildcard 如果模式是通配符模式则返回true
func (u *URLPattern) IsWildcard() bool {
	return strings.Contains(u.path, "...")
}

// IsExact 如果模式是精确匹配模式则返回true
func (u *URLPattern) IsExact() bool {
	return strings.HasSuffix(u.path, "{$}")
}

// Priority 返回路由模式的优先级
// 高优先级模式应该优先匹配
func (u *URLPattern) Priority() int {
	priority := 0

	// 精确模式具有最高优先级
	if u.IsExact() {
		priority += 1000
	}

	// 字面段增加优先级
	segments := strings.Split(strings.Trim(u.path, "/"), "/")
	for _, segment := range segments {
		if !strings.Contains(segment, "{") {
			priority += 100
		}
	}

	// 参数较少则优先级更高
	paramCount := strings.Count(u.path, "{")
	priority -= paramCount * 10

	// 通配符模式优先级较低
	if u.IsWildcard() {
		priority -= 50
	}

	return priority
}

// validatePath 验证路径格式
func validatePath(path string) error {
	// 检查无效字符
	if strings.ContainsAny(path, " \t\n\r") {
		return errors.NewValidationError("path", "路径不能包含空白字符")
	}

	// 检查双斜杠
	if strings.Contains(path, "//") {
		return errors.NewValidationError("path", "路径不能包含连续的斜杠")
	}

	// 检查参数语法
	openBraces := strings.Count(path, "{")
	closeBraces := strings.Count(path, "}")
	if openBraces != closeBraces {
		return errors.NewValidationError("path", "路径参数括号不匹配")
	}

	return nil
}

// compilePathRegex 将路径模式编译为正则表达式
func compilePathRegex(path string) (*regexp.Regexp, error) {
	// 转义正则特殊字符，除了我们的参数语法
	pattern := regexp.QuoteMeta(path)

	// 将转义的参数模式替换为正则表达式组
	// {param} -> (?P<param>[^/]+)
	paramRegex := regexp.MustCompile(`\\{([^}]+)\\}`)
	pattern = paramRegex.ReplaceAllStringFunc(pattern, func(match string) string {
		// 提取参数名
		paramName := strings.Trim(match, "\\{}")

		// 处理通配符参数
		if strings.HasSuffix(paramName, "...") {
			paramName = strings.TrimSuffix(paramName, "...")
			return fmt.Sprintf("(?P<%s>.*)", paramName)
		}

		// 处理精确匹配
		if paramName == "$" {
			return "$"
		}

		// 常规参数
		return fmt.Sprintf("(?P<%s>[^/]+)", paramName)
	})

	// 锚定模式
	pattern = "^" + pattern + "$"

	return regexp.Compile(pattern)
}
