package server

import "net/http"

// This file into impliment the code to redirect to the desired URI

func (ss *Session) Redirect(_uri Uri, _code Response) {
	http.Redirect(ss.W, ss.R, string(_uri), int(_code))
}
