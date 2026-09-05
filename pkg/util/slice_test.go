package util_test

import (
	"testing"

	"github.com/LalatinaHub/LatinaServer/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestRemoveDuplicate(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no duplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with duplicates",
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.RemoveDuplicate(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCheckProxyIP_Invalid(t *testing.T) {
	// Empty string
	info, err := util.CheckProxyIP("")
	assert.NoError(t, err)
	assert.False(t, info.ProxyIP)

	// Invalid format (no port)
	info, err = util.CheckProxyIP("invalid-format")
	assert.Error(t, err)
	assert.False(t, info.ProxyIP)
}
