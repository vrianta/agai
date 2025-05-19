package server

import (
	"fmt"
	"net/http"
	"time"
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

func NewSession(w http.ResponseWriter, r *http.Request) *Session {
	fmt.Println("Creating New Session")
	return &Session{
		w:    w,
		r:    r,
		POST: make(PostParams),
		GET:  make(GetParams),
		Store: SessionVars{
			"uid":        "Guest",
			"isLoggedIn": false,
		},

		RenderEngine: NewRenderHandlerObj(w),
	}
}

func GetSessionID(r *http.Request) *string {
	cookie := GetCookie("sessionid", r)
	if cookie != nil {
		return &cookie.Value
	}
	return nil
}

func (sh *Session) Login(uid string) {
	// WriteConsole("Attempting to Login")
	sh.Store["uid"] = uid
	sh.Store["isLoggedIn"] = true
	// If no valid session ID is found, create a new session
	sh.SetSessionCookie(&sh.ID)
}

func (s *Session) Logout(_redirect_uri string) {
	s.Store["uid"] = ""
	s.Store["isLoggedIn"] = false

	WriteLog("Loggingout")
	s.w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	s.w.Header().Set("Pragma", "no-cache")
	s.w.Header().Set("Expires", "0")
	s.RedirectWithCode(Uri(_redirect_uri), ResponseCodes.SeeOther)

}

/*
 * Checking if the user is logged in
 * @return false -> if the user is not logged in
 */
func (s *Session) IsLoggedIn() bool {
	if _is_loggedin, present := s.Store["isLoggedIn"]; !present {
		return false
	} else {
		// WriteLog("sending is loggedinclear")
		// WriteLog(_is_loggedin.(bool))
		return _is_loggedin.(bool)
	}

}

// StartSession attempts to retrieve or create a new session
func (s *Session) StartSession() *string {

	if sessionID := GetSessionID(s.r); sessionID != nil {
		if *sessionID == "expire" {
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

func (sh *Session) UpdateSession(_w *http.ResponseWriter, _r *http.Request) {
	sh.w = *_w
	sh.r = _r

	sh.RenderEngine.W = *_w
}

// Creates a new session and sets cookies
func (sh *Session) CreateNewSession() *string {
	// Generate a session ID
	sessionID, err := GenerateSessionID()
	if err != nil {
		return nil
	}

	sh.ID = sessionID
	sh.SetSessionCookie(&sessionID)

	return &sessionID
}

// Sets the session cookie in the client's browser
func (sh *Session) SetSessionCookie(sessionID *string) {
	c := &http.Cookie{
		Name:     "sessionid",
		Value:    *sessionID,
		HttpOnly: true,
		Expires:  time.Now().Add(30 * time.Minute).UTC(),
	}
	AddCookie(c, sh.w, sh.r)
}

func (s *Session) EndSession() {
	RemoveCookie("sessionid", s.w, s.r)
	// defer RemoveSession(s.ID)
	s = NewSession(s.w, s.r)
}

func (sh *Session) ParseRequest() {
	// Initialize queryParams once for later use
	queryParams := sh.r.URL.Query()

	sh.POST = make(PostParams)
	sh.GET = make(GetParams)

	// Check if the request method is POST
	if sh.r.Method == http.MethodPost {
		// Parse multipart form data with a 10 MB limit for file uploads
		err := sh.r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			// http.Error(sh.w, "Error parsing multipart form data", http.StatusBadRequest)
		}
		// Handle POST form data
		for key, values := range sh.r.PostForm {
			sh.ProcessPostParams(key, values)
		}
	}

	// Log handling of query parameters for non-POST methods
	for key, values := range queryParams {
		sh.ProcessQueryParams(key, values)
	}
}

// handleQueryParams processes parameters found in the URL query
func (sh *Session) ProcessQueryParams(key string, values []string) {
	var err error
	// Check for multiple values

	if len(values) > 1 {
		if sh.GET[key], err = JsonToString(values); err != nil {
			http.Error(sh.w, "Failed to convert data to JSON", http.StatusMethodNotAllowed)

		}
	} else {
		sh.GET[key] = values[0] // Store single value as a string
	}
}

// handlePostParams processes parameters found in the POST data
func (sh *Session) ProcessPostParams(key string, values []string) {
	var err error
	if len(values) > 1 {
		if sh.POST[key], err = JsonToString(values); err != nil {
			http.Error(sh.w, "Failed to convert data to JSON", http.StatusMethodNotAllowed)
		}
	} else {
		sh.POST[key] = values[0] // Store single value as a string
	}
}

/*
 * Return True if the connection established is a post connection
 */
func (ss *Session) IsPostMethod() bool {
	return ss.r.Method == http.MethodPost
}

/*
 * Return True if the connection established is a Get connection
 */
func (ss *Session) IsGetMethod() bool {
	return ss.r.Method == http.MethodGet
}

/*
 * Return True if the connection established is a DELET connection
 */
func (ss *Session) IsDeleteMethod() bool {
	return ss.r.Method == http.MethodDelete
}

/*
 * Disable the Caching on Local Machine For certain pages to make
 */
func (s *Session) EnableCaching() {

}
