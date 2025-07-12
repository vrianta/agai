package session

import (
	"container/heap"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	Config "github.com/vrianta/agai/v1/config"
	Cookies "github.com/vrianta/agai/v1/cookies"
	Log "github.com/vrianta/agai/v1/log"
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
* - ProcessqueryBuilderParams: Processes and stores queryBuilder parameters.
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
* Author: Pritam Dutta
* Date: NA
 */

func init() {
	gob.Register(map[string]*Struct{})
	gob.Register(&Struct{})
	gob.Register(atomic.Int64{})
	heap.Init(&sessionHeap)
	// LoadAllSessionsFromDisk()

	switch Config.GetWebConfig().SessionStoreType {
	case "disk":
		loadAllSessionsFromDisk()
	}
}

func New() *Struct {
	return &Struct{
		POST:     make(PostParams, 20),
		GET:      make(GetParams, 20),
		LoggedIn: false,
		Expiry:   time.Now().Add(time.Second * 30),
		Store: SessionVars{
			"uid": "Guest",
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
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	// Delete the session from the map
	delete(all, (*sessionID))

	go saveAllSessionsToDisk()
}

// Function to Get Session using Session ID in
// the request, if it exists
func Get(sessionID *string) (*Struct, bool) {
	if sessionID == nil {
		return nil, false
	}

	// for id, session := range all {
	// 	fmt.Printf("Session ID: %s\n", id)
	// 	fmt.Printf("  POST: %+v\n", session.POST)
	// 	fmt.Printf("  GET: %+v\n", session.GET)
	// 	fmt.Printf("  Store: %+v\n", session.Store)
	// 	fmt.Printf("  isLoggedIn: %v\n", session.isLoggedIn)
	// 	fmt.Printf("  Expiry: %s\n", session.Expiry)
	// 	fmt.Println("------------------------")
	// }

	sessionMutex.Lock()
	session, exists := all[*sessionID]
	if !exists {
		sessionMutex.Unlock()
		return nil, false
	}
	sessionMutex.Unlock()

	// session.lastUsed.Store(time.Now().UnixMicro())

	go func(sessionID string) {
		// Move to front in LRU
		lruMutex.Lock()
		if elem, ok := lruMap[sessionID]; ok {
			lruList.MoveToFront(elem)
		}
		lruMutex.Unlock()
	}(*sessionID)

	return session, true
}

func Store(sessionID *string, session *Struct) {
	if sessionID == nil || session == nil {
		return
	}

	sessionMutex.Lock()
	if len(all) >= Config.GetWebConfig().MaxSessionCount {
		evictLRUSession()
	}
	// session.lastUsed.Store(time.Now().UnixMicro())
	all[*sessionID] = session
	sessionMutex.Unlock()

	heapMutex.Lock()
	heap.Push(&sessionHeap, session)
	heapMutex.Unlock()

	switch Config.GetWebConfig().SessionStoreType {
	case "disk":
		go saveAllSessionsToDisk()
	}

	select {
	case lruUpdateChan <- lruUpdate{SessionID: *sessionID, Op: "move"}:
	default:
	}
	select {
	case sessionWakeupChan <- struct{}{}:
	default:
	}
}

func saveAllSessionsToDisk() error {
	// Check if the file exists
	_, err := os.Stat("sessions.data")
	// fileExists := !os.IsNotExist(err)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error checking session file: %w", err)
	}

	// Open file for write (create if not exists, truncate if exists)
	f, err := os.Create("sessions.data")
	if err != nil {
		return fmt.Errorf("failed to create or open session store file: %w", err)
	}
	defer f.Close()

	// Encode the sessions
	enc := gob.NewEncoder(f)

	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

	if err := enc.Encode(all); err != nil {
		panic(err.Error())
		// fmt.Printf("failed to encode session map: %w", err)
		// return err
	}

	// if fileExists {
	// 	fmt.Println("[Sessions] sessions.data updated with latest sessions")
	// } else {
	// 	fmt.Println("[Sessions] sessions.data created and sessions saved")
	// }

	return nil
}

func loadAllSessionsFromDisk() error {
	// Check if sessions.data exists
	if _, err := os.Stat("sessions.data"); os.IsNotExist(err) {
		fmt.Println("[Sessions] sessions.data not found — skipping load")
		return nil
	} else if err != nil {
		// Some other filesystem error
		return fmt.Errorf("error checking session file: %w", err)
	}

	// Open file and decode
	f, err := os.Open("sessions.data")
	if err != nil {
		return fmt.Errorf("failed to open session store file: %w", err)
	}
	defer f.Close()

	var loaded map[string]*Struct
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&loaded); err != nil {
		return fmt.Errorf("failed to decode session map: %w", err)
	}

	sessionMutex.Lock()
	all = loaded
	sessionMutex.Unlock()

	fmt.Printf("[Sessions] Loaded %d sessions from sessions.data\n", len(all))
	// for id, session := range all {
	// 	fmt.Printf("Session ID: %s\n", id)
	// 	fmt.Printf("  POST: %+v\n", session.POST)
	// 	fmt.Printf("  GET: %+v\n", session.GET)
	// 	fmt.Printf("  Store: %+v\n", session.Store)
	// 	fmt.Printf("  isLoggedIn: %v\n", session.LoggedIn)
	// 	fmt.Printf("  Expiry: %s\n", session.Expiry)
	// 	fmt.Println("------------------------")
	// }
	return nil
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
		heapMutex.Lock()
		if sessionHeap.Len() == 0 {
			heapMutex.Unlock()
			time.Sleep(30 * time.Minute)
			continue
		}

		next := sessionHeap[0] // Peek the earliest
		heapMutex.Unlock()

		now := time.Now()
		if next.Expiry.Before(now) {
			RemoveSession(&next.ID)
			Log.WriteLog("Session expired: " + next.ID)
			continue
		}

		sleepDuration := time.Until(next.Expiry)
		if sleepDuration < 0 {
			sleepDuration = 0
		}

		select {
		case <-time.After(sleepDuration):
		case <-sessionWakeupChan:
		}
	}
}

// Function to handle LRU updates in a separate goroutine
// This function listens for updates on the lruUpdateChan channel and processes them
// It handles two operations: "move" to move an existing session to the front of the LRU list
// and "InsertRow" to add a new session to the LRU list if it doesn't already exist.
func StartLRUHandler() {
	for update := range lruUpdateChan {
		sessionMutex.Lock()
		switch update.Op {
		case "move":
			if elem, ok := lruMap[update.SessionID]; ok {
				lruList.MoveToFront(elem)
			}
		case "InsertRow":
			if _, ok := lruMap[update.SessionID]; !ok {
				elem := lruList.PushFront(update.SessionID)
				lruMap[update.SessionID] = elem
			}
		}
		sessionMutex.Unlock()
	}
}

func (sh *Struct) Login(w http.ResponseWriter, r *http.Request) {
	sh.LoggedIn = true
	sh.SetSessionCookie(&sh.ID, w, r)
}

/*
 * Checking if the user is logged in
 * @return false -> if the user is not logged in
 */
func (s *Struct) IsLoggedIn() bool {
	return s.LoggedIn
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

	// sh.lastUsed.Store(time.Now().UnixMicro())
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
	// sh.lastUsed.Store(now.UnixMicro())
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
