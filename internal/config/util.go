package config

import (
	"os"

	appErrors "github.com/LalatinaHub/LatinaServer/pkg/errors"
	"github.com/sagernet/sing/common/json"
)

// SaveJsonToFile serializes content to formatted JSON and writes to disk.

func SaveJsonToFile(filename string, content any) error {
	f, err := os.Create(filename)
	if err != nil {
		return appErrors.NewConfigError("failed to create config file", err)
	}
	defer f.Close()

	b, err := json.Marshal(content)
	if err != nil {
		return appErrors.NewConfigError("failed to marshal json", err)
	}

	if _, err := f.Write(b); err != nil {
		return appErrors.NewConfigError("failed to write json to file", err)
	}

	return nil
}


