package server

import (
	"net/http"
)

// var storeMutex

func (sh *Server) routingHandler(w http.ResponseWriter, r *http.Request) {
	// WriteConsole(r.Header)
	// Log the incoming request URL
	// WriteConsole("Received request for URL: ", r.URL.Path)

	sessionID := GetSessionID(r)
	if sessionID == nil { // means no session has been established with the user
		// WriteConsole("No session found, starting a new session")
		Session := NewSession(w, r)
		sessionID = Session.StartSession()
		if sessionID != nil { // Successfuly started a New session without any error
			WriteConsolef("New session started with ID: %s \n", *sessionID)
			sh.Sessions[(*sessionID)] = *Session
			if value, ok := sh.Routes[r.URL.Path]; ok {
				WriteConsolef("Route found for URL: %s, calling handler \n", r.URL.Path)
				Session.UpdateSession(&w, r)
				Session.ParseRequest()
				value(Session)
				Session.RenderEngine.StartRender()
			} else {
				WriteConsolef("Route not found for URL: %s \n", r.URL.Path)
				http.Error(w, "404 Error : Route not found ", 404)
			}
		} else {
			// WriteConsole("Failed to start session")
			renderhandeler := NewRenderHandlerObj(w)
			renderhandeler.Render(GetResponse("RELOGIN", "Server is not have the session anymore need to relogin the session", false))
			renderhandeler.StartRender()
			return
		}
	} else {
		WriteConsolef("Session ID found: %s \n", *sessionID)

		if Session, ok := sh.Sessions[(*sessionID)]; ok { // session is already created
			if value, ok := sh.Routes[r.URL.Path]; ok {
				WriteConsolef("Route found for URL: %s, calling handler\n", r.URL.Path)
				Session.UpdateSession(&w, r)
				Session.ParseRequest()
				value(&Session)
				Session.RenderEngine.StartRender()
			} else {
				WriteConsolef("Route not found for URL: %s\n", r.URL.Path)
				http.Error(w, "404 Error : Route not found ", 404)
			}

		} else {
			WriteConsole("Session does not exist in Session, creating a new one")

			Session := NewSession(w, r)
			sessionID = Session.StartSession()
			if sessionID != nil {
				WriteConsolef("New session started with ID: %s\n", *sessionID)
				sh.Sessions[(*sessionID)] = *Session

				if value, ok := sh.Routes[r.URL.Path]; ok {
					WriteConsolef("Route found for URL: %s, calling handler\n", r.URL.Path)
					Session.ParseRequest()
					value(Session)
					Session.RenderEngine.StartRender()
				} else {
					WriteConsolef("Route not found for URL: %s\n", r.URL.Path)
					http.Error(w, "404 Error : Route not found ", 404)
				}
			} else {
				WriteConsole("Failed to start session")
				renderhandeler := NewRenderHandlerObj(w)
				renderhandeler.Render(GetResponse("RELOGIN", "Server is not have the session anymore need to relogin the session", false))
				renderhandeler.StartRender()
				return
			}
		}
	}
}
