package server

import (
	"net/http"
	"os"
)

var (
	fileInfo = map[string]FileInfo{}
)

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

		if Session, ok := sh.Sessions[(*sessionID)]; ok {
			if controller, ok := sh.Routes[r.URL.Path]; ok {
				Session.UpdateSession(&w, r)
				Session.ParseRequest()
				controller(&Session)
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
}

/*
 * Handling the Requests coming for the CSS Files specially
 */
func (s *server) CSSHandlers(w http.ResponseWriter, r *http.Request) {

	_file_path := "." + r.URL.Path // path of the file

	_file_record, file_record_ok := fileInfo[_file_path]

	info, err := os.Stat(_file_path)
	if err != nil {
		WriteLog(err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/css; charset=utf-8")

	if file_record_ok && _file_record.LastModified.Compare(info.ModTime()) != 0 { // file has not been updated so alreay no not make it do extra work

		w.Write([]byte(fileInfo[_file_path].Data))
		return
	}

	_file_data := ReadFromFile(_file_path)

	fileInfo[_file_path] = FileInfo{
		Uri:          _file_path,
		LastModified: info.ModTime(),
		Data:         _file_data,
	}

	w.Write([]byte(fileInfo[_file_path].Data))

}

/*
 * Handling the Requests coming for the Js Files specially
 */
func (s *server) JsHandler(w http.ResponseWriter, r *http.Request) {
	_file_path := "." + r.URL.Path // path of the file

	_file_record, file_record_ok := fileInfo[_file_path]

	info, err := os.Stat(_file_path)
	if err != nil {
		WriteLog(err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/css; charset=utf-8")

	if file_record_ok && _file_record.LastModified.Compare(info.ModTime()) != 0 { // file has not been updated so alreay no not make it do extra work

		w.Write([]byte(fileInfo[_file_path].Data))
		return
	}

	_file_data := ReadFromFile(_file_path)

	fileInfo[_file_path] = FileInfo{
		Uri:          _file_path,
		LastModified: info.ModTime(),
		Data:         _file_data,
	}

	w.Write([]byte(fileInfo[_file_path].Data))

}
