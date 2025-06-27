package Session

import (
	"net/http"
	"time"

	"github.com/vrianta/Server/Config"
	"github.com/vrianta/Server/Cookies"
	"github.com/vrianta/Server/Log"
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

func New() *Struct {
	return &Struct{
		POST:       make(PostParams, 20),
		GET:        make(GetParams, 20),
		isLoggedIn: false,
		Expiry:     time.Now().Add(time.Second * 30),
		Store: SessionVars{
			"uid":        "Guest",
			"isLoggedIn": false,
		},
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
	mutex.Lock()
	defer mutex.Unlock()
	session, exists := all[*sessionID]
	if !exists {
		return nil, false
	}
	// Move to front in LRU
	if elem, ok := lruMap[*sessionID]; ok {
		lruList.MoveToFront(elem)
	}
	session.LastUsed = time.Now()
	return session, true
}

func Store(sessionID *string, session *Struct) {
	if sessionID == nil || session == nil {
		return
	}
	mutex.Lock()
	if len(all) >= Config.GetWebConfig().MaxSessionCount {
		evictLRUSession()
	}
	session.LastUsed = time.Now()
	all[*sessionID] = session
	mutex.Unlock()

	// Send LRU update (non-blocking)
	select {
	case lruUpdateChan <- lruUpdate{SessionID: *sessionID, Op: "move"}:
	default:
		// Drop if channel is full to avoid blocking
	}

	select {
	case sessionWakeupChan <- struct{}{}:
	default:
	}
}

func evictLRUSession() {
	// Remove from the back of the list (least recently used)
	elem := lruList.Back()
	if elem == nil {
		return
	}
	sessionID := elem.Value.(string)
	delete(all, sessionID)
	delete(lruMap, sessionID)
	lruList.Remove(elem)
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
			continue // Immediately check for next expiry
		}

		var sleepDuration time.Duration
		if minExpiry != nil {
			sleepDuration = time.Until(minExpiry.Expiry)
			if sleepDuration < 0 {
				sleepDuration = 0
			}
		} else {
			sleepDuration = 30 * time.Minute
		}

		select {
		case <-time.After(sleepDuration):
			// Timer expired, loop to clean up
		case <-sessionWakeupChan:
			// New session or expiry update, re-evaluate minExpiry
		}
	}
}

// Function to handle LRU updates in a separate goroutine
// This function listens for updates on the lruUpdateChan channel and processes them
// It handles two operations: "move" to move an existing session to the front of the LRU list
// and "insert" to add a new session to the LRU list if it doesn't already exist.

func StartLRUHandler() {
	for update := range lruUpdateChan {
		mutex.Lock()
		switch update.Op {
		case "move":
			if elem, ok := lruMap[update.SessionID]; ok {
				lruList.MoveToFront(elem)
			}
		case "insert":
			if _, ok := lruMap[update.SessionID]; !ok {
				elem := lruList.PushFront(update.SessionID)
				lruMap[update.SessionID] = elem
			}
		}
		mutex.Unlock()
	}
}

func (sh *Struct) Login(w http.ResponseWriter, r *http.Request) {
	sh.isLoggedIn = true
	sh.SetSessionCookie(&sh.ID, w, r)
}

/*
 * Checking if the user is logged in
 * @return false -> if the user is not logged in
 */
func (s *Struct) IsLoggedIn() bool {
	return s.isLoggedIn
}

// StartSession attempts to retrieve or create a new session and returnt he created session ID
func (s *Struct) StartSession(sessionID *string, w http.ResponseWriter, r *http.Request) *string {

	// if sessionID := GetSessionID(s.R); sessionID != nil && (*sessionID) != s.ID {
	// 	// If the session ID doesn't match the current handler's ID, create a new session
	// 	defer RemoveSession(sessionID) // Remove the old session
	// }

	// If no valid session ID is found, create a new session
	return s.CreateNewSession(sessionID, w, r)
}

func (sh *Struct) Update(_w http.ResponseWriter, _r *http.Request) {
	updateMutex.Lock()
	defer updateMutex.Unlock()

	sh.LastUsed = time.Now()
}

// function to clear value of POST and GET from the Session
// Make sure what ever in the store will stay for as long as the server is not stopped
// or you remove the data intentionally
func (sh *Struct) Clean() {
	cleanMutex.Lock()
	defer cleanMutex.Unlock()

	sh.POST = make(PostParams, 20)
	sh.GET = make(GetParams, 20)
}

// Creates a new session and sets cookies
func (sh *Struct) CreateNewSession(sessionID *string, w http.ResponseWriter, r *http.Request) *string {
	// Generate a session ID
	if sessionID == nil {
		return nil
	}

	sh.ID = *sessionID
	sh.SetSessionCookie(sessionID, w, r)
	// Wake up the session handler if needed
	select {
	case sessionWakeupChan <- struct{}{}:
	default:
	}
	return sessionID
}

// Sets the session cookie in the client's browser
func (sh *Struct) SetSessionCookie(sessionID *string, w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	sh.Expiry = now.Add(30 * time.Minute).UTC()
	sh.LastUsed = now
	c := &http.Cookie{
		Name:     "sessionid",
		Value:    *sessionID,
		HttpOnly: Config.GetWebConfig().Https,
		Expires:  sh.Expiry,
	}

	Cookies.AddCookie(c, w, r)
}

/*
 * Disable the Caching on Local Machine For certain pages to make
 */
func (s *Struct) EnableCaching() {

}
