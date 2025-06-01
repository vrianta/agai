package Session

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vrianta/Server/Cookies"
	"github.com/vrianta/Server/Log"
	"github.com/vrianta/Server/RenderEngine"
	"github.com/vrianta/Server/Utils"
)

/*
* Session Management Package
*
* This package provides a robust session management system for handling
* user authentication, session creation, and data persistence.
*
* Features:
* - Session Creation and Management
* - Secure Cookie Handling
* - User Login Management
* - Request Parsing for GET and POST methods
* - Session Expiry Control
*
* Structures:
* - Session: Manages individual user sessions with methods to parse requests,
*   handle cookies, and manage session data.
*
* Functions:
* - NewSession: Initializes a new session with default values.
* - GetSessionID: Retrieves the session ID from cookies.
* - Login: Handles user login by updating session data.
* - IsLoggedIn: Checks if the user is logged in.
* - StartSession: Starts or resumes a session.
* - UpdateSession: Updates the session with new HTTP request/response data.
* - CreateNewSession: Generates a new session ID and assigns it.
* - SetSessionCookie: Adds a session ID cookie to the HTTP response.
* - EndSession: Ends the current session by removing associated data and cookies.
* - ParseRequest: Parses HTTP request parameters for GET and POST data.
* - ProcessQueryParams: Processes and stores query parameters.
* - ProcessPostParams: Processes and stores POST form data.
*
* Usage:
* - Import the package
* - Use `NewSession()` to initialize a session in your handler functions.
* - Call `StartSession()` to begin session tracking.
* - Use `IsLoggedIn()` to verify the user's authentication status.
* - Manage session data directly using the `Store` map.
*
* Example:
*
* func MyHandler(w http.ResponseWriter, r *http.Request) {
*     session := NewSession(w, r)
*     session.StartSession()
*
*     if session.IsLoggedIn() {
*         fmt.Fprintln(w, "Welcome back, ", session.Store["uid"])
*     } else {
*         fmt.Fprintln(w, "Please log in.")
*     }
* }
*
* Best Practices:
* - Ensure that session IDs are securely generated to avoid session fixation attacks.
* - Use HTTPS to encrypt cookies for improved security.
* - Regularly invalidate stale sessions to reduce security risks.
*
* Author: [Your Name]
* Date: [Current Date]
 */

func New(w http.ResponseWriter, r *http.Request) *Struct {
	// fmt.Println("Creating New Session")
	return &Struct{
		W:    w,
		R:    r,
		POST: make(PostParams),
		GET:  make(GetParams),
		Store: SessionVars{
			"uid":        "Guest",
			"isLoggedIn": false,
		},

		RenderEngine: RenderEngine.New(w),
	}
}

func GetSessionID(r *http.Request) *string {
	cookie := Cookies.GetCookie("sessionid", r)
	if cookie != nil {
		return &cookie.Value
	}
	return nil
}

// Function to Remove Session using Session ID in
// the request, if it exists
func RemoveSession(sessionID *string) {
	if sessionID == nil {
		return
	}

	// Lock the mutex for writing
	mutex.Lock()
	defer mutex.Unlock()

	// Delete the session from the map
	delete(all, (*sessionID))
}

// Function to Get Session using Session ID in
// the request, if it exists
func Get(sessionID *string) (*Struct, bool) {
	if sessionID == nil {
		return nil, false
	}

	// Lock the mutex for reading
	mutex.RLock()
	defer mutex.RUnlock()

	// Retrieve the session from the map
	session, exists := all[(*sessionID)]
	if !exists {
		return nil, false // Session not found
	}

	return session, true
}

// Function to Set the Session using Session ID in
// the request, if it exists
func Set(sessionID *string, session *Struct) {
	if sessionID == nil || session == nil {
		return
	}

	// Lock the mutex for reading
	mutex.RLock()
	defer mutex.RUnlock()

	// Store the session in the map
	all[(*sessionID)] = session
}

// Function to keep checking the session expiry to remve them.
// How the function will work is it will loop through all the current sessions and check when the expiry is time.Time
// Later it find the session which is expired and the one which has the least time to expire.
// It will Expire the expired session and the go on sleep until the least time to expire is reached.
// then it will again go on loop to to repeast the process.
func StartSessionHandler() {
	for {
		if all == nil {
			time.Sleep(30 * time.Minute) // Sleep if no sessions are found
			continue
		}

		mutex.Lock()
		var minExpiry *Struct
		for _, sess := range all {
			if minExpiry == nil || sess.Expiry.Before(minExpiry.Expiry) {
				minExpiry = sess
			}
		}
		mutex.Unlock()

		if minExpiry != nil && minExpiry.Expiry.Before(time.Now()) {
			RemoveSession(&minExpiry.ID)
			Log.WriteLog("Session expired: " + minExpiry.ID)
		}

		if minExpiry != nil {
			time.Sleep(time.Until(minExpiry.Expiry))
		} else {
			time.Sleep(30 * time.Minute)
		}
	}
}

func (sh *Struct) Login(uid string) {
	// WriteConsole("Attempting to Login")
	sh.Store["uid"] = uid
	sh.Store["isLoggedIn"] = true
	// If no valid session ID is found, create a new session
	sh.SetSessionCookie(&sh.ID)
}

func (s *Struct) Logout() {
	s.Store["uid"] = ""
	s.Store["isLoggedIn"] = false

	Log.WriteLog("Loggingout")
	s.W.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	s.W.Header().Set("Pragma", "no-cache")
	s.W.Header().Set("Expires", "0")

}

/*
 * Checking if the user is logged in
 * @return false -> if the user is not logged in
 */
func (s *Struct) IsLoggedIn() bool {
	if _is_loggedin, present := s.Store["isLoggedIn"]; !present {
		return false
	} else {
		// WriteLog("sending is loggedinclear")
		// WriteLog(_is_loggedin.(bool))
		return _is_loggedin.(bool)
	}

}

// StartSession attempts to retrieve or create a new session
func (s *Struct) StartSession() *string {

	if sessionID := GetSessionID(s.R); sessionID != nil {
		if (*sessionID) == "expire" {
			return s.CreateNewSession()
		}
		// If the session ID doesn't match the current handler's ID, create a new session
		if (*sessionID) != s.ID {
			s.EndSession()
		}
	}

	// If no valid session ID is found, create a new session
	return s.CreateNewSession()
}

func (sh *Struct) UpdateSession(_w http.ResponseWriter, _r *http.Request) {
	sh.W = _w
	sh.R = _r

	sh.RenderEngine.W = _w
}

// Creates a new session and sets cookies
func (sh *Struct) CreateNewSession() *string {
	// Generate a session ID
	sessionID, err := Utils.GenerateSessionID()
	if err != nil {
		return nil
	}

	sh.ID = sessionID
	sh.SetSessionCookie(&sessionID)

	return &sessionID
}

// Sets the session cookie in the client's browser
func (sh *Struct) SetSessionCookie(sessionID *string) {
	sh.Expiry = time.Now().Add(30 * time.Minute).UTC()
	c := &http.Cookie{
		Name:     "sessionid",
		Value:    *sessionID,
		HttpOnly: true,
		Expires:  sh.Expiry,
	}

	Cookies.AddCookie(c, sh.W, sh.R)
}

func (s *Struct) EndSession() {
	Cookies.RemoveCookie("sessionid", s.W, s.R)
	// defer RemoveSession(s.ID)
	s = New(s.W, s.R)
}

func (sh *Struct) ParseRequest() {
	// Initialize queryParams once for later use
	// queryParams := sh.R.URL.Query()

	sh.POST = make(PostParams)
	sh.GET = make(GetParams)

	if sh.IsPostMethod() {
		contentType := sh.R.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			// Handle JSON data
			// if err := Utils.ParseJSONBody(sh.R, &sh.POST); err != nil {
			// 	http.Error(sh.W, fmt.Sprintf("Error parsing JSON body | Error - %s", err.Error()), http.StatusBadRequest)
			// 	return
			// }
			break

		case "application/x-www-form-urlencoded":
			// Handle form data (application/x-www-form-urlencoded)
			err := sh.R.ParseForm()
			if err != nil {
				http.Error(sh.W, fmt.Sprintf("Error parsing form data | Error - %s", err.Error()), http.StatusBadRequest)
				return
			}

		case "multipart/form-data":
			// Handle multipart form data (file upload)
			// Note: This case is handled separately below
			if err := sh.R.ParseMultipartForm(10 << 20); err != nil { // 10 MB
				http.Error(sh.W, fmt.Sprintf("Error parsing multipart form data | Error - %s", err.Error()), http.StatusBadRequest)
				return
			}

		default:
			break
		}

		// Log handling of query parameters for non-POST methods
		for key, values := range sh.R.PostForm {
			sh.ProcessPostParams(key, values)
		}
	}

	// Log handling of query parameters for non-POST methods
	for key, values := range sh.R.URL.Query() {
		sh.ProcessQueryParams(key, values)
	}
}

// handleQueryParams processes parameters found in the URL query
func (sh *Struct) ProcessQueryParams(key string, values []string) {
	var err error
	// Check for multiple values

	if len(values) > 1 {
		if sh.GET[key], err = Utils.JsonToString(values); err != nil {
			http.Error(sh.W, "Failed to convert data to JSON", http.StatusMethodNotAllowed)

		}
	} else {
		sh.GET[key] = values[0] // Store single value as a string
	}
}

// handlePostParams processes parameters found in the POST data
func (sh *Struct) ProcessPostParams(key string, values []string) {
	var err error
	if len(values) > 1 {
		if sh.POST[key], err = Utils.JsonToString(values); err != nil {
			http.Error(sh.W, "Failed to convert data to JSON", http.StatusMethodNotAllowed)
		}
	} else {
		sh.POST[key] = values[0] // Store single value as a string
	}
}

/*
 * Return True if the connection established is a post connection
 */
func (ss *Struct) IsPostMethod() bool {
	return ss.R.Method == http.MethodPost
}

/*
 * Return True if the connection established is a Get connection
 */
func (ss *Struct) IsGetMethod() bool {
	return ss.R.Method == http.MethodGet
}

/*
 * Return True if the connection established is a DELET connection
 */
func (ss *Struct) IsDeleteMethod() bool {
	return ss.R.Method == http.MethodDelete
}

/*
 * Disable the Caching on Local Machine For certain pages to make
 */
func (s *Struct) EnableCaching() {

}
