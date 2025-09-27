package controller

import (
	"net/http"

	Session "github.com/vrianta/agai/v1/internal/session"
	"github.com/vrianta/agai/v1/view"
)

// Routes is a map of HTTP methods to their respective controllers
type (
	Context struct {
		session *Session.Instance // Session object to handle user session

		// privte objects
		w http.ResponseWriter
		r *http.Request

		userInputs map[string]interface{}
	}

	View = func() view.Context
)
