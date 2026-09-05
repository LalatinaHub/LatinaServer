package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveJsonToFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "util_save_test.json")

	content := map[string]interface{}{
		"server": "example.com",
		"port":   8080,
	}

	err := config.SaveJsonToFile(filePath, content)
	require.NoError(t, err)

	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "example.com")
	assert.Contains(t, string(data), "8080")
}

func TestReadCaddyConfig(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("non-existent file returns error", func(t *testing.T) {
		_, err := config.ReadCaddyConfig(filepath.Join(tempDir, "missing.json"))
		assert.Error(t, err)
	})

	t.Run("invalid json returns error", func(t *testing.T) {
		badFile := filepath.Join(tempDir, "bad.json")
		err := os.WriteFile(badFile, []byte("{invalid-json"), 0644)
		require.NoError(t, err)

		_, err = config.ReadCaddyConfig(badFile)
		assert.Error(t, err)
	})

	t.Run("valid caddy json parses", func(t *testing.T) {
		validFile := filepath.Join(tempDir, "valid_caddy.json")
		err := os.WriteFile(validFile, []byte(`{"admin":{"disabled":true}}`), 0644)
		require.NoError(t, err)

		cfg, err := config.ReadCaddyConfig(validFile)
		require.NoError(t, err)
		assert.NotNil(t, cfg)
	})
}

func TestReadSingConfig(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("non-existent file returns error", func(t *testing.T) {
		_, err := config.ReadSingConfig(filepath.Join(tempDir, "missing_sing.json"))
		assert.Error(t, err)
	})

	t.Run("invalid json returns error", func(t *testing.T) {
		badFile := filepath.Join(tempDir, "bad_sing.json")
		err := os.WriteFile(badFile, []byte("not a json"), 0644)
		require.NoError(t, err)

		_, err = config.ReadSingConfig(badFile)
		assert.Error(t, err)
	})

	t.Run("valid sing-box json parses", func(t *testing.T) {
		validFile := filepath.Join(tempDir, "valid_sing.json")
		err := os.WriteFile(validFile, []byte(`{"log":{"level":"info"}}`), 0644)
		require.NoError(t, err)

		options, err := config.ReadSingConfig(validFile)
		require.NoError(t, err)
		assert.NotNil(t, options.Log)
	})
}
