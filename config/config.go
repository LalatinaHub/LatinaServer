package config

import (
	"os"

	"github.com/sagernet/sing/common/json"
)

func SaveJsonToFile(filename string, content any) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := json.Marshal(content)
	f.WriteString(string(b))

	if err != nil {
		panic(err)
	}
}
