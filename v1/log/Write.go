package log

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	Config "github.com/vrianta/agai/v1/config"
)

// LogLevel enum
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var currentLevel = DEBUG
var useJSON = false

func SetLogLevel(level string) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		currentLevel = DEBUG
	case "INFO":
		currentLevel = INFO
	case "WARN":
		currentLevel = WARN
	case "ERROR":
		currentLevel = ERROR
	}
}

func EnableJSONOutput(enable bool) {
	useJSON = enable
}

func log(level LogLevel, label string, color string, msg string, args ...any) {
	if Config.GetBuild() {
		return
	}
	if level < currentLevel {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatted := fmt.Sprintf(msg, args...)

	if useJSON {
		logEntry := map[string]any{
			"timestamp": timestamp,
			"level":     label,
			"message":   formatted,
		}
		if data, err := json.Marshal(logEntry); err == nil {
			fmt.Println(string(data))
		}
	} else {
		fmt.Printf("%s[%s] %s: %s\033[0m\n", color, timestamp, label, formatted)
	}
}

func Success(msg string, args ...any) {
	if Config.GetBuild() {
		return
	}
	log(INFO, "[SUCCESS]", "\033[32m", msg, args...)
}

// Colored log wrappers
func Debug(msg string, args ...any) {
	if Config.GetBuild() {
		return
	}
	log(DEBUG, "[DEBUG]", "\033[36m", msg, args...)
}
func Info(msg string, args ...any) {
	if Config.GetBuild() {
		return
	}
	log(INFO, "[INFO]", "\033[32m", msg, args...)
}
func Warn(msg string, args ...any) {
	if Config.GetBuild() {
		return
	}
	log(WARN, "[WARN]", "\033[33m", msg, args...)
}
func Error(msg string, args ...any) {
	if Config.GetBuild() {
		return
	}
	log(ERROR, "[ERROR]", "\033[31m", msg, args...)
}
func Write(msg string, args ...any) {
	if Config.GetBuild() {
		return
	}
	green := "\033[32m"
	reset := "\033[0m"
	formatted := fmt.Sprintf(msg, args...)
	fmt.Printf("%s%s%s", green, formatted, reset)
}

// Legacy support
func WriteLog(messages ...any) {
	if Config.GetBuild() {
		return
	}
	fmt.Println(messages...)
}

func WriteLogf(msg string, args ...any) {
	if Config.GetBuild() {
		return
	}
	fmt.Printf(msg, args...)
}

// Standard response wrapper
func GetResponse(code string, message string, success bool) string {
	result, _ := json.Marshal(struct {
		Code    string `json:"CODE"`
		Message string `json:"MESSAGE"`
		Success bool   `json:"SUCCESS"`
	}{
		Code:    code,
		Message: message,
		Success: success,
	})
	return string(result)
}
