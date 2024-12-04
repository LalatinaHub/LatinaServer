package config

import (
	"encoding/json"
	"os"
)

func SaveJsonToFile(filename string, content any) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := json.MarshalIndent(content, "", "\t")
	f.WriteString(string(b))

	if err != nil {
		panic(err)
	}
}
