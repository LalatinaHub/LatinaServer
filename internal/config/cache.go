package config

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"

	"github.com/sagernet/sing/common/json"
)

var (
	configHashCache = make(map[string]string)
	cacheMu         sync.RWMutex
)

// computeHash computes SHA256 hash of the given content.
func computeHash(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

// SaveJsonToFileWithCache writes JSON to file only if content hash has changed.
// Returns true if file was written (config changed), false if skipped (no change).
func SaveJsonToFileWithCache(filename string, content any) (bool, error) {
	b, err := json.Marshal(content)
	if err != nil {
		return false, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	newHash := computeHash(b)

	cacheMu.RLock()
	oldHash, exists := configHashCache[filename]
	cacheMu.RUnlock()

	if exists && oldHash == newHash {
		// No change, skip writing
		return false, nil
	}

	// Write to file
	if err := os.WriteFile(filename, b, 0644); err != nil {
		return false, fmt.Errorf("failed to write file %s: %w", filename, err)
	}

	// Update cache
	cacheMu.Lock()
	configHashCache[filename] = newHash
	cacheMu.Unlock()

	return true, nil
}
