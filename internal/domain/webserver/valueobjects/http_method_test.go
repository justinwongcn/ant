package valueobjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPMethod(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "valid GET method",
			input:       "GET",
			expected:    "GET",
			expectError: false,
		},
		{
			name:        "valid POST method lowercase",
			input:       "post",
			expected:    "POST",
			expectError: false,
		},
		{
			name:        "valid method with spaces",
			input:       "  PUT  ",
			expected:    "PUT",
			expectError: false,
		},
		{
			name:        "empty method",
			input:       "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid method",
			input:       "INVALID",
			expected:    "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			input:       "   ",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method, err := NewHTTPMethod(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, method)
			} else {
				require.NoError(t, err)
				require.NotNil(t, method)
				assert.Equal(t, tt.expected, method.Value())
				assert.Equal(t, tt.expected, method.String())
			}
		})
	}
}

func TestMustNewHTTPMethod(t *testing.T) {
	t.Run("valid method", func(t *testing.T) {
		method := MustNewHTTPMethod("GET")
		assert.Equal(t, "GET", method.Value())
	})

	t.Run("invalid method panics", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNewHTTPMethod("INVALID")
		})
	})
}

func TestHTTPMethodEquals(t *testing.T) {
	get1 := MustNewHTTPMethod("GET")
	get2 := MustNewHTTPMethod("GET")
	post := MustNewHTTPMethod("POST")

	assert.True(t, get1.Equals(get2))
	assert.False(t, get1.Equals(post))
	assert.False(t, get1.Equals(nil))
}

func TestHTTPMethodProperties(t *testing.T) {
	tests := []struct {
		method       string
		isIdempotent bool
		isSafe       bool
		allowsBody   bool
	}{
		{"GET", true, true, false},
		{"POST", false, false, true},
		{"PUT", true, false, true},
		{"DELETE", true, false, false},
		{"PATCH", false, false, true},
		{"HEAD", true, true, false},
		{"OPTIONS", true, true, false},
		{"TRACE", true, true, false},
		{"CONNECT", false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			method := MustNewHTTPMethod(tt.method)

			assert.Equal(t, tt.isIdempotent, method.IsIdempotent(), "IsIdempotent")
			assert.Equal(t, tt.isSafe, method.IsSafe(), "IsSafe")
			assert.Equal(t, tt.allowsBody, method.AllowsBody(), "AllowsBody")
		})
	}
}

func TestPredefinedHTTPMethods(t *testing.T) {
	assert.Equal(t, "GET", HTTPMethodGET.Value())
	assert.Equal(t, "POST", HTTPMethodPOST.Value())
	assert.Equal(t, "PUT", HTTPMethodPUT.Value())
	assert.Equal(t, "DELETE", HTTPMethodDELETE.Value())
	assert.Equal(t, "PATCH", HTTPMethodPATCH.Value())
	assert.Equal(t, "HEAD", HTTPMethodHEAD.Value())
	assert.Equal(t, "OPTIONS", HTTPMethodOPTIONS.Value())
	assert.Equal(t, "TRACE", HTTPMethodTRACE.Value())
	assert.Equal(t, "CONNECT", HTTPMethodCONNECT.Value())
}

func TestGetAllHTTPMethods(t *testing.T) {
	methods := GetAllHTTPMethods()
	assert.Len(t, methods, 9)

	// Check that all expected methods are present
	expectedMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE", "CONNECT"}
	actualMethods := make([]string, len(methods))
	for i, method := range methods {
		actualMethods[i] = method.Value()
	}

	for _, expected := range expectedMethods {
		assert.Contains(t, actualMethods, expected)
	}
}

func TestIsValidHTTPMethod(t *testing.T) {
	tests := []struct {
		method string
		valid  bool
	}{
		{"GET", true},
		{"post", true},
		{"  PUT  ", true},
		{"INVALID", false},
		{"", false},
		{"   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			assert.Equal(t, tt.valid, IsValidHTTPMethod(tt.method))
		})
	}
}

func BenchmarkNewHTTPMethod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewHTTPMethod("GET")
	}
}

func BenchmarkHTTPMethodEquals(b *testing.B) {
	method1 := MustNewHTTPMethod("GET")
	method2 := MustNewHTTPMethod("GET")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		method1.Equals(method2)
	}
}
