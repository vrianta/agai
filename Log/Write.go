package Log

import (
	"encoding/json"
	"fmt"
)

func WriteLog(massages ...any) {
	fmt.Println(massages...)
}

func WriteLogf(massage string, args ...any) {
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
