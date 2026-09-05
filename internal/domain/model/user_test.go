package model_test

import (
	"testing"
	"time"

	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
)

func TestUser_IsActive(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		user     model.User
		now      time.Time
		expected bool
	}{
		{
			name: "active user",
			user: model.User{
				Quota:      100,
				Expired:    now.Add(24 * time.Hour),
				ServerCode: "SG",
				VPN:        "trojan",
			},
			now:      now,
			expected: true,
		},
		{
			name: "quota zero",
			user: model.User{
				Quota:      0,
				Expired:    now.Add(24 * time.Hour),
				ServerCode: "SG",
				VPN:        "trojan",
			},
			now:      now,
			expected: false,
		},
		{
			name: "expired user",
			user: model.User{
				Quota:      100,
				Expired:    now.Add(-24 * time.Hour),
				ServerCode: "SG",
				VPN:        "trojan",
			},
			now:      now,
			expected: false,
		},
		{
			name: "missing server code",
			user: model.User{
				Quota:      100,
				Expired:    now.Add(24 * time.Hour),
				ServerCode: "",
				VPN:        "trojan",
			},
			now:      now,
			expected: false,
		},
		{
			name: "missing VPN type",
			user: model.User{
				Quota:      100,
				Expired:    now.Add(24 * time.Hour),
				ServerCode: "SG",
				VPN:        "",
			},
			now:      now,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.IsActive(tt.now); got != tt.expected {
				t.Errorf("IsActive() = %v, want %v", got, tt.expected)
			}
		})
	}
}
