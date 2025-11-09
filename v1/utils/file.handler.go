package utils

import (
	"os"
	"path/filepath"
)

// takes path of the file and return data as string format
func ReadFromFile(uri string) []byte {
	if content, err := os.ReadFile(uri); err == nil {
		return content
	} else {
		return nil
	}
}

func JoinPath(parts ...string) string {
	return filepath.Join(parts...)
}
