package Redirect

import (
	"net/http"

	"github.com/vrianta/Server/Response"
	"github.com/vrianta/Server/Session"
)

// This file into impliment the code to redirect to the desired URI

func New(_Uri string, session *Session.Struct) *response {
	return &response{
		uri:     _Uri,
		session: session,
	}
}

// Redirects to the URI user provided
func (r *response) Redirect() {
	http.Redirect(r.session.W, r.session.R, r.uri, int(Response.Codes.TemporaryRedirect)) // if golang developpers worked so hard to create this why should I do it again :P
	// http.RedirectHandler()
}

func (r *response) WithCode(_code Response.Code) {
	http.Redirect(r.session.W, r.session.R, r.uri, int(_code))
}
