package model_test

import (
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestServer_IsFull(t *testing.T) {
	tests := []struct {
		name     string
		server   model.Server
		expected bool
	}{
		{
			name: "server has capacity",
			server: model.Server{
				UsersCount: 10,
				UsersMax:   50,
			},
			expected: false,
		},
		{
			name: "server at max capacity",
			server: model.Server{
				UsersCount: 50,
				UsersMax:   50,
			},
			expected: true,
		},
		{
			name: "server over capacity",
			server: model.Server{
				UsersCount: 55,
				UsersMax:   50,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.server.IsFull())
		})
	}
}
