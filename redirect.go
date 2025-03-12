package server

import "net/http"

// This file into impliment the code to redirect to the desired URI

// Redirects to the URI user provided
func (ss *Session) Redirect(_uri Uri) {
	http.Redirect(ss.w, ss.r, string(_uri), int(ResponseCodes.TemporaryRedirect)) // if golang developpers worked so hard to create this why should I do it again :P
	// http.RedirectHandler()
}

func (s *Session) RedirectWithCode(_uri Uri, _code ResponseCode) {
	http.Redirect(s.w, s.r, string(_uri), int(_code))
}
