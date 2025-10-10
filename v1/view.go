package agai

import (
	"encoding/json"
	"log"
)

type view struct {
	Name     string // name of the view
	AsJson   bool   // indicate if we have send the value as json not to some view
	Response interface {
		Get() any
	}
}

func (c *view) ToJson() []byte {
	jsonBytes, err := json.MarshalIndent(c.Response, "", "  ")
	if err != nil {
		log.Println("Failed to marshal response to JSON:", err)
		return []byte("{}")
	}
	return jsonBytes
}
