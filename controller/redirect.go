package controller

import (
	"net/http"

	Response "github.com/vrianta/Server/response"
)

// Redirects to the URI user provided
func (_c *Struct) Redirect(uri string) {
	http.Redirect(_c.w, _c.r, uri, int(Response.Codes.TemporaryRedirect)) // if golang developpers worked so hard to create this why should I do it again :P
}

func (_c *Struct) WithCode(uri string, _code Response.Code) {
	http.Redirect(_c.w, _c.r, uri, int(_code))
}
