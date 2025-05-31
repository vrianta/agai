package RenderEngine

import "net/http"

type (
	Struct struct {
		view []byte
		W    http.ResponseWriter
	}

	RenderData map[string]interface{}
)
