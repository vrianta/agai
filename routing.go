package server

import (
	"fmt"
	"net/http"
)

// var storeMutex

func (sh *server) routingHandler(w http.ResponseWriter, r *http.Request) {

	sessionID := GetSessionID(r)
	if sessionID == nil { // means user does not have any session with the server so creating a new clean guest session with the server
		Session := NewSession(w, r)
		sessionID = Session.StartSession()
		if sessionID != nil { // Successfuly started a New session without any error
			sh.Sessions[(*sessionID)] = *Session
			if value, ok := sh.Routes[r.URL.Path]; ok {
				Session.UpdateSession(&w, r)
				Session.ParseRequest()
				value(Session)
				Session.RenderEngine.StartRender()
			} else {
				// WriteConsolef("Route not found for URL: %s \n", r.URL.Path)
				http.Error(w, "404 Error : Route not found ", 404)
			}
		} else {
			http.Error(w, "Server Error * Failed to Create the Session for the user", 500)
		}
	} else { // User has a session ID to begin with
		// checking if the session is valid or not means it is checking if the server also has the session or not
		// if the session is valid then it will just update the session with the latest value
		fmt.Println("Session cookie found with value:", *sessionID)

		if Session, ok := sh.Sessions[(*sessionID)]; ok {
			if controller, ok := sh.Routes[r.URL.Path]; ok {
				Session.UpdateSession(&w, r)
				Session.ParseRequest()
				controller(&Session)
				WriteLog("leaving the Controller Calling", r.URL.Path)
				Session.RenderEngine.StartRender()
			} else {
				http.Error(w, "404 Error : Route not found ", 404)
			}
		} else { // server is not holding the session any more so creating a new guest session for the user
			Session := NewSession(w, r)
			sessionID = Session.StartSession()
			if sessionID != nil {
				sh.Sessions[(*sessionID)] = *Session
				if controller, ok := sh.Routes[r.URL.Path]; ok {
					Session.ParseRequest()
					controller(Session)
					Session.RenderEngine.StartRender()
				} else {
					http.Error(w, "404 Error : Route not found ", 404)
				}
			} else {
				http.Error(w, "Server Error * Failed to Create the Session for the user", 500)
			}
		}
	}

	// w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	// w.Header().Set("Pragma", "no-cache")
	// w.Header().Set("Expires", "0")
}
