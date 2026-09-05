package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveJsonToFileWithCache(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_config.json")

	payload1 := map[string]string{"foo": "bar"}
	payload2 := map[string]string{"foo": "baz"}

	// First write should write to file
	changed, err := SaveJsonToFileWithCache(filePath, payload1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Fatal("first write should report changed=true")
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("file was not created")
	}

	// Second write with identical content should skip
	changed, err = SaveJsonToFileWithCache(filePath, payload1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Fatal("second write with identical content should report changed=false")
	}

	// Third write with new content should write
	changed, err = SaveJsonToFileWithCache(filePath, payload2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Fatal("third write with changed content should report changed=true")
	}
}
