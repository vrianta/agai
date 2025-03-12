package server

import "net/http"

// This file into impliment the code to redirect to the desired URI

func (ss *Session) Redirect(_uri Uri, _code ResponseCode) {
	http.Redirect(ss.w, ss.r, string(_uri), int(_code)) // if golang developpers worked so hard to create this why should I do it again :P
}
