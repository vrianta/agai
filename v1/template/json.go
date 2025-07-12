package template

import (
	"encoding/json"
	"log"
)

func (r Response) AsJson() []byte {
	jsonBytes, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Println("Failed to marshal response to JSON:", err)
		return []byte("{}")
	}
	return jsonBytes
}
