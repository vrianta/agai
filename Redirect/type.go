package Redirect

import "github.com/vrianta/Server/Session"

type (
	response struct {
		uri     string
		session *Session.Struct
	}
)
