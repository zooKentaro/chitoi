package config

import (
	"os"
	"path/filepath"
)

func Home() string {
	gopath := os.Getenv("GOPATH")
	return filepath.Join(gopath, "src", "github.com", "uenoryo", "chitoi")
}
