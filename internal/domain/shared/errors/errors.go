// Package errors 定义Ant Web框架的通用领域错误
package errors

import (
	"fmt"
)

// DomainError 表示特定于领域的错误
type DomainError struct {
	Code    string
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}

// NewDomainError 创建新的领域错误
func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

// NewDomainErrorWithCause 创建带有原因的新领域错误
func NewDomainErrorWithCause(code, message string, cause error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// 通用领域错误码
const (
	// Web服务器领域错误
	ErrCodeServerAlreadyRunning = "SERVER_ALREADY_RUNNING"
	ErrCodeServerNotRunning     = "SERVER_NOT_RUNNING"
	ErrCodeInvalidRoute         = "INVALID_ROUTE"
	ErrCodeRouteAlreadyExists   = "ROUTE_ALREADY_EXISTS"
	ErrCodeInvalidMiddleware    = "INVALID_MIDDLEWARE"

	// 会话领域错误
	ErrCodeSessionNotFound    = "SESSION_NOT_FOUND"
	ErrCodeSessionExpired     = "SESSION_EXPIRED"
	ErrCodeInvalidSessionID   = "INVALID_SESSION_ID"
	ErrCodeSessionDataInvalid = "SESSION_DATA_INVALID"

	// 文件领域错误
	ErrCodeFileNotFound         = "FILE_NOT_FOUND"
	ErrCodeFileAlreadyExists    = "FILE_ALREADY_EXISTS"
	ErrCodeInvalidFilePath      = "INVALID_FILE_PATH"
	ErrCodeFilePermissionDenied = "FILE_PERMISSION_DENIED"
	ErrCodeFileSizeExceeded     = "FILE_SIZE_EXCEEDED"
	ErrCodeUnsupportedFileType  = "UNSUPPORTED_FILE_TYPE"

	// 请求处理领域错误
	ErrCodeInvalidRequest   = "INVALID_REQUEST"
	ErrCodeInvalidResponse  = "INVALID_RESPONSE"
	ErrCodeInvalidParameter = "INVALID_PARAMETER"
	ErrCodeMissingParameter = "MISSING_PARAMETER"
)

// 预定义的领域错误
var (
	// Web服务器领域
	ErrServerAlreadyRunning = NewDomainError(ErrCodeServerAlreadyRunning, "服务器已经在运行")
	ErrServerNotRunning     = NewDomainError(ErrCodeServerNotRunning, "服务器未运行")
	ErrInvalidRoute         = NewDomainError(ErrCodeInvalidRoute, "无效的路由配置")
	ErrRouteAlreadyExists   = NewDomainError(ErrCodeRouteAlreadyExists, "路由已存在")
	ErrInvalidMiddleware    = NewDomainError(ErrCodeInvalidMiddleware, "无效的中间件")

	// Session领域
	ErrSessionNotFound    = NewDomainError(ErrCodeSessionNotFound, "会话不存在")
	ErrSessionExpired     = NewDomainError(ErrCodeSessionExpired, "会话已过期")
	ErrInvalidSessionID   = NewDomainError(ErrCodeInvalidSessionID, "无效的会话ID")
	ErrSessionDataInvalid = NewDomainError(ErrCodeSessionDataInvalid, "会话数据无效")

	// 文件领域
	ErrFileNotFound         = NewDomainError(ErrCodeFileNotFound, "文件不存在")
	ErrFileAlreadyExists    = NewDomainError(ErrCodeFileAlreadyExists, "文件已存在")
	ErrInvalidFilePath      = NewDomainError(ErrCodeInvalidFilePath, "无效的文件路径")
	ErrFilePermissionDenied = NewDomainError(ErrCodeFilePermissionDenied, "文件权限被拒绝")
	ErrFileSizeExceeded     = NewDomainError(ErrCodeFileSizeExceeded, "文件大小超出限制")
	ErrUnsupportedFileType  = NewDomainError(ErrCodeUnsupportedFileType, "不支持的文件类型")

	// 请求处理领域
	ErrInvalidRequest   = NewDomainError(ErrCodeInvalidRequest, "无效的请求")
	ErrInvalidResponse  = NewDomainError(ErrCodeInvalidResponse, "无效的响应")
	ErrInvalidParameter = NewDomainError(ErrCodeInvalidParameter, "无效的参数")
	ErrMissingParameter = NewDomainError(ErrCodeMissingParameter, "缺少必需的参数")
)

// ValidationError 表示验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NewValidationError 创建一个新的验证错误
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// ValidationErrors 表示多个验证错误
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "no validation errors"
	}
	if len(e) == 1 {
		return e[0].Error()
	}
	return fmt.Sprintf("validation errors: %d errors occurred", len(e))
}

// Add 添加验证错误
func (e *ValidationErrors) Add(field, message string) {
	*e = append(*e, ValidationError{Field: field, Message: message})
}

// HasErrors 如果有验证错误则返回true
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}
