package utils

import (
	"os"
	"path/filepath"
)

// takes path of the file and return data as string format
func ReadFromFile(uri string) string {
	if content, err := os.ReadFile(uri); err == nil {
		return string(content)
	} else {
		return err.Error()
	}
}

func JoinPath(parts ...string) string {
	return filepath.Join(parts...)
}
