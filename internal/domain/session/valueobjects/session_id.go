// Package valueobjects 包含会话域的值对象
package valueobjects

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/justinwongcn/ant/internal/domain/shared/errors"
)

// SessionID 表示会话标识符值对象
type SessionID struct {
	value string
}

const (
	// DefaultSessionIDLength 是会话ID的默认长度
	DefaultSessionIDLength = 32
	// MinSessionIDLength 是会话ID的最小允许长度
	MinSessionIDLength = 16
	// MaxSessionIDLength 是会话ID的最大允许长度
	MaxSessionIDLength = 128
)

// sessionIDPattern 定义会话ID的有效模式(字母数字和连字符)
var sessionIDPattern = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)

// NewSessionID 创建一个带验证的新SessionID值对象
func NewSessionID(value string) (*SessionID, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return nil, errors.NewValidationError("sessionID", "会话ID不能为空")
	}

	if len(value) < MinSessionIDLength {
		return nil, errors.NewValidationError("sessionID", fmt.Sprintf("会话ID长度不能少于%d个字符", MinSessionIDLength))
	}

	if len(value) > MaxSessionIDLength {
		return nil, errors.NewValidationError("sessionID", fmt.Sprintf("会话ID长度不能超过%d个字符", MaxSessionIDLength))
	}

	if !sessionIDPattern.MatchString(value) {
		return nil, errors.NewValidationError("sessionID", "会话ID只能包含字母、数字、连字符和下划线")
	}

	return &SessionID{value: value}, nil
}

// MustNewSessionID 创建一个新的SessionID值对象，出错时panic
func MustNewSessionID(value string) *SessionID {
	sessionID, err := NewSessionID(value)
	if err != nil {
		panic(err)
	}
	return sessionID
}

// GenerateSessionID 生成一个新的随机会话ID
func GenerateSessionID() (*SessionID, error) {
	return GenerateSessionIDWithLength(DefaultSessionIDLength)
}

// GenerateSessionIDWithLength 生成指定长度的新随机会话ID
func GenerateSessionIDWithLength(length int) (*SessionID, error) {
	if length < MinSessionIDLength {
		return nil, errors.NewValidationError("length", fmt.Sprintf("会话ID长度不能少于%d个字符", MinSessionIDLength))
	}

	if length > MaxSessionIDLength {
		return nil, errors.NewValidationError("length", fmt.Sprintf("会话ID长度不能超过%d个字符", MaxSessionIDLength))
	}

	// Generate random bytes
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return nil, errors.NewDomainErrorWithCause(
			errors.ErrCodeInvalidSessionID,
			"生成会话ID失败",
			err,
		)
	}

	// Convert to hex string
	value := hex.EncodeToString(bytes)

	// Ensure we have the exact length requested
	if len(value) > length {
		value = value[:length]
	}

	return &SessionID{value: value}, nil
}

// MustGenerateSessionID 生成一个新的随机会话ID，出错时panic
func MustGenerateSessionID() *SessionID {
	sessionID, err := GenerateSessionID()
	if err != nil {
		panic(err)
	}
	return sessionID
}

// Value 返回会话ID的字符串值
func (s *SessionID) Value() string {
	return s.value
}

// String 返回字符串表示
func (s *SessionID) String() string {
	return s.value
}

// Equals 检查两个SessionID对象是否相等
func (s *SessionID) Equals(other *SessionID) bool {
	if other == nil {
		return false
	}
	return s.value == other.value
}

// IsEmpty 如果会话ID为空则返回true
func (s *SessionID) IsEmpty() bool {
	return s.value == ""
}

// Length 返回会话ID的长度
func (s *SessionID) Length() int {
	return len(s.value)
}

// IsValid 检查会话ID是否有效
func (s *SessionID) IsValid() bool {
	if s.value == "" {
		return false
	}

	if len(s.value) < MinSessionIDLength || len(s.value) > MaxSessionIDLength {
		return false
	}

	return sessionIDPattern.MatchString(s.value)
}

// Truncate 将会话ID截断到指定长度
func (s *SessionID) Truncate(length int) (*SessionID, error) {
	if length < MinSessionIDLength {
		return nil, errors.NewValidationError("length", fmt.Sprintf("会话ID长度不能少于%d个字符", MinSessionIDLength))
	}

	if length >= len(s.value) {
		return s, nil
	}

	truncated := s.value[:length]
	return &SessionID{value: truncated}, nil
}

// HasPrefix 检查会话ID是否以给定前缀开头
func (s *SessionID) HasPrefix(prefix string) bool {
	return strings.HasPrefix(s.value, prefix)
}

// HasSuffix 检查会话ID是否以给定后缀结尾
func (s *SessionID) HasSuffix(suffix string) bool {
	return strings.HasSuffix(s.value, suffix)
}

// Contains 检查会话ID是否包含给定子字符串
func (s *SessionID) Contains(substr string) bool {
	return strings.Contains(s.value, substr)
}

// ToUpper 返回一个大写值的新SessionID
func (s *SessionID) ToUpper() *SessionID {
	return &SessionID{value: strings.ToUpper(s.value)}
}

// ToLower 返回一个小写值的新SessionID
func (s *SessionID) ToLower() *SessionID {
	return &SessionID{value: strings.ToLower(s.value)}
}

// Mask 返回用于日志记录的会话ID的掩码版本
func (s *SessionID) Mask() string {
	if len(s.value) <= 8 {
		return "****"
	}

	prefix := s.value[:4]
	suffix := s.value[len(s.value)-4:]
	return fmt.Sprintf("%s****%s", prefix, suffix)
}

// IsValidSessionIDString 检查字符串是否是有效的会话ID而不创建对象
func IsValidSessionIDString(value string) bool {
	value = strings.TrimSpace(value)

	if value == "" {
		return false
	}

	if len(value) < MinSessionIDLength || len(value) > MaxSessionIDLength {
		return false
	}

	return sessionIDPattern.MatchString(value)
}
