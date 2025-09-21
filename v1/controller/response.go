package controller

import (
	"encoding/json"
	"log"
)

type (
	view struct {
		name     string // name of the view
		asJson   bool   // indicate if we have send the value as json not to some view
		response *Response
	}
	Response map[string]any
)

/**
 * @param - name : name of the View where you want to send the respnse
 **/
func (r *Response) ToView(name string) view {
	return view{
		name:     name,
		response: r,
	}
}

func (r *Response) AsJson() view {
	return view{
		asJson:   true,
		response: r,
	}
}

func (r Response) toJson() []byte {
	jsonBytes, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Println("Failed to marshal response to JSON:", err)
		return []byte("{}")
	}
	return jsonBytes
}
