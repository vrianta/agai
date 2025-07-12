package log

import (
	"encoding/json"
	"fmt"

	Config "github.com/vrianta/agai/v1/internal/config"
)

func WriteLog(massages ...any) {
	if Config.GetBuild() {
		return
	}
	fmt.Println(massages...)

}

func WriteLogf(massage string, args ...any) {
	if Config.GetBuild() {
		return
	}
	fmt.Printf(massage, args...)
}

func GetResponse(_ErrorCode string, _Message string, _Success bool) string {
	result, _ := json.Marshal(struct {
		Code    string `json:"CODE"`
		Message string `json:"MESSAGE"`
		Success bool   `json:"SUCCESS"`
	}{
		Code:    _ErrorCode,
		Message: _Message,
		Success: _Success,
	})
	return string(result)
}
