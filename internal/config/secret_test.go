package config_test

import (
	"os"
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGetSSPassword(t *testing.T) {
	// Without env var, should return a valid non-empty random string
	pass1 := config.GetSSPassword()
	assert.NotEmpty(t, pass1)
	assert.Len(t, pass1, 32) // 16 bytes = 32 hex chars

	// Subsequent calls without env var should return cached value
	pass2 := config.GetSSPassword()
	assert.Equal(t, pass1, pass2)

	// With env var, should prioritize env var
	expectedEnvPass := "my-secure-custom-password-123"
	os.Setenv("SS_PASSWORD", expectedEnvPass)
	defer os.Unsetenv("SS_PASSWORD")

	assert.Equal(t, expectedEnvPass, config.GetSSPassword())
}

func TestGetClashSecret(t *testing.T) {
	// Clean env
	os.Unsetenv("CLASH_SECRET")
	os.Unsetenv("YACD_PASSWORD")

	assert.Empty(t, config.GetClashSecret())

	os.Setenv("CLASH_SECRET", "custom-clash-secret")
	assert.Equal(t, "custom-clash-secret", config.GetClashSecret())
	os.Unsetenv("CLASH_SECRET")

	os.Setenv("YACD_PASSWORD", "custom-yacd-pass")
	assert.Equal(t, "custom-yacd-pass", config.GetClashSecret())
	os.Unsetenv("YACD_PASSWORD")
}
