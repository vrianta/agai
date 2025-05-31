package Session

import (
	"net/http"

	"github.com/vrianta/Server/RenderEngine"
)

type (
	SessionVars map[string]any
	PostParams  map[string]string
	GetParams   map[string]string

	Struct struct {
		ID string
		W  http.ResponseWriter
		R  *http.Request

		POST  PostParams
		GET   GetParams
		Store SessionVars

		RenderEngine RenderEngine.Struct
	}
)
