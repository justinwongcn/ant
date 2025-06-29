// Package valueobjects 包含web服务器领域的值对象
package valueobjects

import (
	"fmt"

	"github.com/justinwongcn/ant/internal/domain/shared/errors"
)

// StatusCode 表示HTTP状态码值对象
type StatusCode struct {
	value int
}

// HTTP状态码常量
const (
	// 1xx 信息类
	StatusContinue           = 100
	StatusSwitchingProtocols = 101
	StatusProcessing         = 102

	// 2xx 成功
	StatusOK                   = 200
	StatusCreated              = 201
	StatusAccepted             = 202
	StatusNonAuthoritativeInfo = 203
	StatusNoContent            = 204
	StatusResetContent         = 205
	StatusPartialContent       = 206

	// 3xx 重定向
	StatusMultipleChoices   = 300
	StatusMovedPermanently  = 301
	StatusFound             = 302
	StatusSeeOther          = 303
	StatusNotModified       = 304
	StatusUseProxy          = 305
	StatusTemporaryRedirect = 307
	StatusPermanentRedirect = 308

	// 4xx 客户端错误
	StatusBadRequest                   = 400
	StatusUnauthorized                 = 401
	StatusPaymentRequired              = 402
	StatusForbidden                    = 403
	StatusNotFound                     = 404
	StatusMethodNotAllowed             = 405
	StatusNotAcceptable                = 406
	StatusProxyAuthRequired            = 407
	StatusRequestTimeout               = 408
	StatusConflict                     = 409
	StatusGone                         = 410
	StatusLengthRequired               = 411
	StatusPreconditionFailed           = 412
	StatusRequestEntityTooLarge        = 413
	StatusRequestURITooLong            = 414
	StatusUnsupportedMediaType         = 415
	StatusRequestedRangeNotSatisfiable = 416
	StatusExpectationFailed            = 417
	StatusTeapot                       = 418
	StatusMisdirectedRequest           = 421
	StatusUnprocessableEntity          = 422
	StatusLocked                       = 423
	StatusFailedDependency             = 424
	StatusTooEarly                     = 425
	StatusUpgradeRequired              = 426
	StatusPreconditionRequired         = 428
	StatusTooManyRequests              = 429
	StatusRequestHeaderFieldsTooLarge  = 431
	StatusUnavailableForLegalReasons   = 451

	// 5xx 服务器错误
	StatusInternalServerError           = 500
	StatusNotImplemented                = 501
	StatusBadGateway                    = 502
	StatusServiceUnavailable            = 503
	StatusGatewayTimeout                = 504
	StatusHTTPVersionNotSupported       = 505
	StatusVariantAlsoNegotiates         = 506
	StatusInsufficientStorage           = 507
	StatusLoopDetected                  = 508
	StatusNotExtended                   = 510
	StatusNetworkAuthenticationRequired = 511
)

// NewStatusCode 创建一个新的StatusCode值对象
func NewStatusCode(code int) (*StatusCode, error) {
	if !isValidStatusCode(code) {
		return nil, errors.NewValidationError("statusCode", fmt.Sprintf("无效的HTTP状态码: %d", code))
	}

	return &StatusCode{value: code}, nil
}

// MustNewStatusCode 创建一个新的StatusCode值对象，出错时panic
func MustNewStatusCode(code int) *StatusCode {
	sc, err := NewStatusCode(code)
	if err != nil {
		panic(err)
	}
	return sc
}

// Value 返回状态码的整数值
func (s *StatusCode) Value() int {
	return s.value
}

// String 返回字符串表示
func (s *StatusCode) String() string {
	return fmt.Sprintf("%d", s.value)
}

// Equals 检查两个StatusCode对象是否相等
func (s *StatusCode) Equals(other *StatusCode) bool {
	if other == nil {
		return false
	}
	return s.value == other.value
}

// IsInformational 如果状态码是信息类(1xx)则返回true
func (s *StatusCode) IsInformational() bool {
	return s.value >= 100 && s.value < 200
}

// IsSuccess 如果状态码表示成功(2xx)则返回true
func (s *StatusCode) IsSuccess() bool {
	return s.value >= 200 && s.value < 300
}

// IsRedirection 如果状态码表示重定向(3xx)则返回true
func (s *StatusCode) IsRedirection() bool {
	return s.value >= 300 && s.value < 400
}

// IsClientError 如果状态码表示客户端错误(4xx)则返回true
func (s *StatusCode) IsClientError() bool {
	return s.value >= 400 && s.value < 500
}

// IsServerError 如果状态码表示服务器错误(5xx)则返回true
func (s *StatusCode) IsServerError() bool {
	return s.value >= 500 && s.value < 600
}

// IsError 如果状态码表示错误(4xx或5xx)则返回true
func (s *StatusCode) IsError() bool {
	return s.IsClientError() || s.IsServerError()
}

// Text 返回状态码对应的状态文本
func (s *StatusCode) Text() string {
	switch s.value {
	case StatusOK:
		return "OK"
	case StatusCreated:
		return "Created"
	case StatusAccepted:
		return "Accepted"
	case StatusNoContent:
		return "No Content"
	case StatusMovedPermanently:
		return "Moved Permanently"
	case StatusFound:
		return "Found"
	case StatusSeeOther:
		return "See Other"
	case StatusNotModified:
		return "Not Modified"
	case StatusTemporaryRedirect:
		return "Temporary Redirect"
	case StatusBadRequest:
		return "Bad Request"
	case StatusUnauthorized:
		return "Unauthorized"
	case StatusForbidden:
		return "Forbidden"
	case StatusNotFound:
		return "Not Found"
	case StatusMethodNotAllowed:
		return "Method Not Allowed"
	case StatusConflict:
		return "Conflict"
	case StatusInternalServerError:
		return "Internal Server Error"
	case StatusNotImplemented:
		return "Not Implemented"
	case StatusBadGateway:
		return "Bad Gateway"
	case StatusServiceUnavailable:
		return "Service Unavailable"
	case StatusGatewayTimeout:
		return "Gateway Timeout"
	default:
		return "Unknown"
	}
}

// isValidStatusCode 检查状态码是否有效
func isValidStatusCode(code int) bool {
	return code >= 100 && code < 600
}

// 预定义的状态码实例
var (
	StatusCodeOK                  = MustNewStatusCode(StatusOK)
	StatusCodeCreated             = MustNewStatusCode(StatusCreated)
	StatusCodeAccepted            = MustNewStatusCode(StatusAccepted)
	StatusCodeNoContent           = MustNewStatusCode(StatusNoContent)
	StatusCodeMovedPermanently    = MustNewStatusCode(StatusMovedPermanently)
	StatusCodeFound               = MustNewStatusCode(StatusFound)
	StatusCodeSeeOther            = MustNewStatusCode(StatusSeeOther)
	StatusCodeNotModified         = MustNewStatusCode(StatusNotModified)
	StatusCodeTemporaryRedirect   = MustNewStatusCode(StatusTemporaryRedirect)
	StatusCodeBadRequest          = MustNewStatusCode(StatusBadRequest)
	StatusCodeUnauthorized        = MustNewStatusCode(StatusUnauthorized)
	StatusCodeForbidden           = MustNewStatusCode(StatusForbidden)
	StatusCodeNotFound            = MustNewStatusCode(StatusNotFound)
	StatusCodeMethodNotAllowed    = MustNewStatusCode(StatusMethodNotAllowed)
	StatusCodeConflict            = MustNewStatusCode(StatusConflict)
	StatusCodeInternalServerError = MustNewStatusCode(StatusInternalServerError)
	StatusCodeNotImplemented      = MustNewStatusCode(StatusNotImplemented)
	StatusCodeBadGateway          = MustNewStatusCode(StatusBadGateway)
	StatusCodeServiceUnavailable  = MustNewStatusCode(StatusServiceUnavailable)
	StatusCodeGatewayTimeout      = MustNewStatusCode(StatusGatewayTimeout)
)
