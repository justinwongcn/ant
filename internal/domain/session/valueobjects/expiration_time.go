// Package valueobjects 包含会话域的值对象
package valueobjects

import (
	"fmt"
	"time"

	"github.com/justinwongcn/ant/internal/domain/shared/errors"
)

// ExpirationTime 表示会话过期时间值对象
type ExpirationTime struct {
	value time.Time
}

const (
	// DefaultSessionDuration 是默认会话持续时间(30分钟)
	DefaultSessionDuration = 30 * time.Minute
	// MinSessionDuration 是最小允许的会话持续时间(1分钟)
	MinSessionDuration = 1 * time.Minute
	// MaxSessionDuration 是最大允许的会话持续时间(24小时)
	MaxSessionDuration = 24 * time.Hour
)

// NewExpirationTime 创建一个新的ExpirationTime值对象
func NewExpirationTime(value time.Time) (*ExpirationTime, error) {
	if value.IsZero() {
		return nil, errors.NewValidationError("expirationTime", "过期时间不能为零值")
	}

	now := time.Now()
	if value.Before(now) {
		return nil, errors.NewValidationError("expirationTime", "过期时间不能早于当前时间")
	}

	duration := value.Sub(now)
	if duration < MinSessionDuration {
		return nil, errors.NewValidationError("expirationTime", fmt.Sprintf("会话持续时间不能少于%v", MinSessionDuration))
	}

	if duration > MaxSessionDuration {
		return nil, errors.NewValidationError("expirationTime", fmt.Sprintf("会话持续时间不能超过%v", MaxSessionDuration))
	}

	return &ExpirationTime{value: value}, nil
}

// MustNewExpirationTime 创建一个新的ExpirationTime值对象，出错时panic
func MustNewExpirationTime(value time.Time) *ExpirationTime {
	expTime, err := NewExpirationTime(value)
	if err != nil {
		panic(err)
	}
	return expTime
}

// NewExpirationTimeFromDuration 从持续时间创建一个新的ExpirationTime
func NewExpirationTimeFromDuration(duration time.Duration) (*ExpirationTime, error) {
	if duration < MinSessionDuration {
		return nil, errors.NewValidationError("duration", fmt.Sprintf("会话持续时间不能少于%v", MinSessionDuration))
	}

	if duration > MaxSessionDuration {
		return nil, errors.NewValidationError("duration", fmt.Sprintf("会话持续时间不能超过%v", MaxSessionDuration))
	}

	expirationTime := time.Now().Add(duration)
	return &ExpirationTime{value: expirationTime}, nil
}

// MustNewExpirationTimeFromDuration 从持续时间创建一个新的ExpirationTime，出错时panic
func MustNewExpirationTimeFromDuration(duration time.Duration) *ExpirationTime {
	expTime, err := NewExpirationTimeFromDuration(duration)
	if err != nil {
		panic(err)
	}
	return expTime
}

// NewDefaultExpirationTime 用默认持续时间创建一个新的ExpirationTime
func NewDefaultExpirationTime() *ExpirationTime {
	return &ExpirationTime{value: time.Now().Add(DefaultSessionDuration)}
}

// Value 返回过期时间的time.Time值
func (e *ExpirationTime) Value() time.Time {
	return e.value
}

// String 返回字符串表示
func (e *ExpirationTime) String() string {
	return e.value.Format(time.RFC3339)
}

// Equals 检查两个ExpirationTime对象是否相等
func (e *ExpirationTime) Equals(other *ExpirationTime) bool {
	if other == nil {
		return false
	}
	return e.value.Equal(other.value)
}

// IsExpired 如果过期时间已过则返回true
func (e *ExpirationTime) IsExpired() bool {
	return time.Now().After(e.value)
}

// IsExpiredAt 如果在给定时间过期时间已过则返回true
func (e *ExpirationTime) IsExpiredAt(at time.Time) bool {
	return at.After(e.value)
}

// TimeUntilExpiration 返回距离过期的持续时间
func (e *ExpirationTime) TimeUntilExpiration() time.Duration {
	now := time.Now()
	if e.value.Before(now) {
		return 0
	}
	return e.value.Sub(now)
}

// TimeUntilExpirationAt 从给定时间返回距离过期的持续时间
func (e *ExpirationTime) TimeUntilExpirationAt(at time.Time) time.Duration {
	if e.value.Before(at) {
		return 0
	}
	return e.value.Sub(at)
}

// TimeSinceExpiration 返回过期后的持续时间(如果已过期)
func (e *ExpirationTime) TimeSinceExpiration() time.Duration {
	now := time.Now()
	if e.value.After(now) {
		return 0
	}
	return now.Sub(e.value)
}

// Extend 将过期时间延长给定的持续时间
func (e *ExpirationTime) Extend(duration time.Duration) (*ExpirationTime, error) {
	if duration < 0 {
		return nil, errors.NewValidationError("duration", "延长时间不能为负数")
	}

	newExpiration := e.value.Add(duration)
	totalDuration := time.Until(newExpiration)

	if totalDuration > MaxSessionDuration {
		return nil, errors.NewValidationError("duration", fmt.Sprintf("延长后的会话持续时间不能超过%v", MaxSessionDuration))
	}

	return &ExpirationTime{value: newExpiration}, nil
}

// Refresh 将过期时间刷新为从现在开始的默认持续时间
func (e *ExpirationTime) Refresh() *ExpirationTime {
	return NewDefaultExpirationTime()
}

// RefreshWithDuration 用给定的持续时间从现在开始刷新过期时间
func (e *ExpirationTime) RefreshWithDuration(duration time.Duration) (*ExpirationTime, error) {
	return NewExpirationTimeFromDuration(duration)
}

// Unix 返回过期时间的Unix时间戳
func (e *ExpirationTime) Unix() int64 {
	return e.value.Unix()
}

// UnixNano 返回过期时间的纳秒级Unix时间戳
func (e *ExpirationTime) UnixNano() int64 {
	return e.value.UnixNano()
}

// Format 使用给定的布局格式化过期时间
func (e *ExpirationTime) Format(layout string) string {
	return e.value.Format(layout)
}

// UTC 返回UTC时区的过期时间
func (e *ExpirationTime) UTC() *ExpirationTime {
	return &ExpirationTime{value: e.value.UTC()}
}

// Local 返回本地时区的过期时间
func (e *ExpirationTime) Local() *ExpirationTime {
	return &ExpirationTime{value: e.value.Local()}
}

// In 返回给定时区的过期时间
func (e *ExpirationTime) In(loc *time.Location) *ExpirationTime {
	return &ExpirationTime{value: e.value.In(loc)}
}

// Before 如果此过期时间早于另一个则返回true
func (e *ExpirationTime) Before(other *ExpirationTime) bool {
	if other == nil {
		return false
	}
	return e.value.Before(other.value)
}

// After 如果此过期时间晚于另一个则返回true
func (e *ExpirationTime) After(other *ExpirationTime) bool {
	if other == nil {
		return true
	}
	return e.value.After(other.value)
}

// Add 将给定持续时间加到过期时间上
func (e *ExpirationTime) Add(duration time.Duration) (*ExpirationTime, error) {
	newTime := e.value.Add(duration)
	return NewExpirationTime(newTime)
}

// Sub 从过期时间中减去给定持续时间
func (e *ExpirationTime) Sub(duration time.Duration) (*ExpirationTime, error) {
	newTime := e.value.Add(-duration)
	return NewExpirationTime(newTime)
}

// IsWithinGracePeriod 如果过期时间在宽限期内则返回true
func (e *ExpirationTime) IsWithinGracePeriod(gracePeriod time.Duration) bool {
	if gracePeriod <= 0 {
		return false
	}

	now := time.Now()
	graceTime := e.value.Add(gracePeriod)
	return now.Before(graceTime)
}

// NewExpirationTimeFromUnix 从Unix时间戳创建一个新的ExpirationTime
func NewExpirationTimeFromUnix(timestamp int64) (*ExpirationTime, error) {
	t := time.Unix(timestamp, 0)
	return NewExpirationTime(t)
}

// IsValidDuration 检查持续时间对于会话过期是否有效
func IsValidDuration(duration time.Duration) bool {
	return duration >= MinSessionDuration && duration <= MaxSessionDuration
}
