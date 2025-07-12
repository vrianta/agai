package controller

import (
	"net/http"

	Response "github.com/vrianta/agai/v1/response"
)

// Redirects to the URI user provided
func (_c *Context) Redirect(uri string) {
	http.Redirect(_c.w, _c.r, uri, int(Response.Codes.TemporaryRedirect)) // if golang developpers worked so hard to create this why should I do it again :P
}

func (_c *Context) WithCode(uri string, _code Response.Code) {
	http.Redirect(_c.w, _c.r, uri, int(_code))
}
