package view

import (
	"encoding/json"
	"log"
)

type Context struct {
	Name     string // name of the view
	AsJson   bool   // indicate if we have send the value as json not to some view
	Response interface {
		Get() any
	}
}

func (c *Context) ToJson() []byte {
	jsonBytes, err := json.MarshalIndent(c.Response, "", "  ")
	if err != nil {
		log.Println("Failed to marshal response to JSON:", err)
		return []byte("{}")
	}
	return jsonBytes
}
