package utils

import "os"

// takes path of the file and return data as string format
func ReadFromFile(uri string) string {
	if content, err := os.ReadFile(uri); err == nil {
		return string(content)
	} else {
		return err.Error()
	}
}
