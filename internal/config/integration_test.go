//go:build integration

package config_test

import (
	"os"
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLiveConfigGeneration(t *testing.T) {
	if os.Getenv("TURSO_DATABASE_URL") == "" {
		t.Skip("TURSO_DATABASE_URL not set; skipping live config generation test")
	}

	err := config.GenerateConfigsParallel()
	assert.NoError(t, err)
}
