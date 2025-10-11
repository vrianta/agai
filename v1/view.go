package agai

import (
	"encoding/json"
	"log"
)

type view struct {
	name     string // name of the view
	asJson   bool   // indicate if we have send the value as json not to some view
	response any
}

func (c *view) ToJson() []byte {
	jsonBytes, err := json.MarshalIndent(c.response, "", "  ")
	if err != nil {
		log.Println("Failed to marshal response to JSON:", err)
		return []byte("{}")
	}
	return jsonBytes
}

// enable redirect when the user is instructed to redirte
// func (c *view)
