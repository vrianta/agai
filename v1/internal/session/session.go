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
	"github.com/vrianta/agai/v1/utils"
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
	gob.Register(map[string]*Instance{})
	gob.Register(&Instance{})
	gob.Register(atomic.Int64{})
	heap.Init(&sessionHeap)
	// LoadAllSessionsFromDisk()

	switch Config.GetWebConfig().SessionStoreType {
	case "disk", "storage":
		loadAllSessionsFromDisk()
	}
}

func New(w http.ResponseWriter, r *http.Request) (*Instance, error) {
	if sessionID, err := utils.GenerateSessionID(); err != nil {
		return nil, err
	} else {
		ins := &Instance{
			ID:              sessionID,
			PostParameters:  make(HTTPPostParameters, 20),
			GetParameters:   make(HTTPGetParameters, 20),
			IsAuthenticated: false,
			ExpirationTime:  time.Now().Add(time.Second * 30),
			Data: SessionData{
				"uid": "Guest",
			},
		}
		go Store(ins)

		ins.setCookie(w, r)
		// Wake up the session handler if needed
		select {
		case cleanupTriggerChan <- struct{}{}:
		default:
		}
		return ins, nil
	}
}

func GetSessionID(r *http.Request) (string, error) {
	cookie, err := Cookies.GetCookie("sessionid", r)
	if cookie != nil {
		return cookie.Value, nil
	}
	return "", err
}

// Function to Remove Session using Session ID in
// the request, if it exists
func RemoveSession(sessionID *string) {
	if sessionID == nil {
		return
	}

	// Lock the mutex for writing
	sessionStoreMutex.Lock()
	defer sessionStoreMutex.Unlock()

	// Delete the session from the map
	delete(instances, (*sessionID))

	switch Config.GetWebConfig().SessionStoreType {
	case "disk", "storage":
		go saveAllSessionsToDisk()
	}
}

// Function to Get Session using Session ID in
// the request, if it exists
func Get(sessionID *string) (*Instance, bool) {
	if sessionID == nil {
		return nil, false
	}

	// for id, session := range instances {
	// 	fmt.Printf("Session ID: %s\n", id)
	// 	fmt.Printf("  POST: %+v\n", session.POST)
	// 	fmt.Printf("  GET: %+v\n", session.GET)
	// 	fmt.Printf("  Store: %+v\n", session.Store)
	// 	fmt.Printf("  isLoggedIn: %v\n", session.isLoggedIn)
	// 	fmt.Printf("  Expiry: %s\n", session.Expiry)
	// 	fmt.Println("------------------------")
	// }

	sessionStoreMutex.Lock()
	session, exists := instances[*sessionID]
	if !exists {
		sessionStoreMutex.Unlock()
		return nil, false
	}
	sessionStoreMutex.Unlock()

	// session.lastUsed.Store(time.Now().UnixMicro())

	go func(sessionID string) {
		// Move to front in LRU
		lruCacheMutex.Lock()
		if elem, ok := lruElementMap[sessionID]; ok {
			lruOrderList.MoveToFront(elem)
		}
		lruCacheMutex.Unlock()
	}(*sessionID)

	return session, true
}

func Store(session *Instance) {
	if session == nil {
		return
	}

	sessionStoreMutex.Lock()
	if len(instances) >= Config.GetWebConfig().MaxSessionCount {
		evictLRUSession()
	}

	instances[session.ID] = session
	sessionStoreMutex.Unlock()

	heapAccessMutex.Lock()
	heap.Push(&sessionHeap, session)
	heapAccessMutex.Unlock()

	switch Config.GetWebConfig().SessionStoreType {
	case "disk", "storage":
		go saveAllSessionsToDisk()
	}

	select {
	case lruOperationChan <- LRUCacheOperation{ID: session.ID, OperationType: "move"}:
	default:
	}
	select {
	case cleanupTriggerChan <- struct{}{}:
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

	sessionStoreMutex.RLock()
	defer sessionStoreMutex.RUnlock()

	if err := enc.Encode(instances); err != nil {
		panic(err.Error())
	}
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

	var loaded map[string]*Instance
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&loaded); err != nil {
		return fmt.Errorf("failed to decode session map: %w", err)
	}

	sessionStoreMutex.Lock()
	instances = loaded
	sessionStoreMutex.Unlock()

	fmt.Printf("[Sessions] Loaded %d sessions from sessions.data\n", len(instances))
	// for id, session := range instances {
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
	elem := lruOrderList.Back()
	if elem == nil {
		return
	}
	sessionID := elem.Value.(string)
	delete(instances, sessionID)
	delete(lruElementMap, sessionID)
	lruOrderList.Remove(elem)
}

// Function to keep checking the session expiry to remve them.
// How the function will work is it will loop through instances the current sessions and check when the expiry is time.Time
// Later it find the session which is expired and the one which has the least time to expire.
// It will Expire the expired session and the go on sleep until the least time to expire is reached.
// then it will again go on loop to to repeast the process.
func StartSessionHandler() {
	for {
		heapAccessMutex.Lock()
		if sessionHeap.Len() == 0 {
			heapAccessMutex.Unlock()
			time.Sleep(30 * time.Minute)
			continue
		}

		next := sessionHeap[0] // Peek the earliest
		heapAccessMutex.Unlock()

		now := time.Now()
		if next.ExpirationTime.Before(now) {
			RemoveSession(&next.ID)
			Log.WriteLog("Session expired: " + next.ID)
			continue
		}

		sleepDuration := time.Until(next.ExpirationTime)
		if sleepDuration < 0 {
			sleepDuration = 0
		}

		select {
		case <-time.After(sleepDuration):
		case <-cleanupTriggerChan:
		}
	}
}

// Function to handle LRU updates in a separate goroutine
// This function listens for updates on the lruOperationChan channel and processes them
// It handles two operations: "move" to move an existing session to the front of the LRU list
// and "InsertRow" to add a new session to the LRU list if it doesn't already exist.
func StartLRUHandler() {
	for update := range lruOperationChan {
		sessionStoreMutex.Lock()
		switch update.OperationType {
		case "move":
			if elem, ok := lruElementMap[update.ID]; ok {
				lruOrderList.MoveToFront(elem)
			}
		case "InsertRow":
			if _, ok := lruElementMap[update.ID]; !ok {
				elem := lruOrderList.PushFront(update.ID)
				lruElementMap[update.ID] = elem
			}
		}
		sessionStoreMutex.Unlock()
	}
}

func (sh *Instance) Login(w http.ResponseWriter, r *http.Request) {
	sh.IsAuthenticated = true
	sh.setCookie(w, r)
}

func (sh *Instance) Update(_w http.ResponseWriter, _r *http.Request) {
	sessionUpdateMutex.Lock()
	defer sessionUpdateMutex.Unlock()

	// sh.lastUsed.Store(time.Now().UnixMicro())
}

// function to clear value of POST and GET from the Session
// Make sure what ever in the store will stay for as long as the server is not stopped
// or you remove the data intentionally
func (sh *Instance) Clean() {
	sessionCleanMutex.Lock()
	defer sessionCleanMutex.Unlock()

	sh.PostParameters = make(HTTPPostParameters)
	sh.GetParameters = make(HTTPGetParameters)
}

// Sets the session cookie in the client's browser
func (sh *Instance) setCookie(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	sh.ExpirationTime = now.Add(30 * time.Minute).UTC()
	// sh.lastUsed.Store(now.UnixMicro())
	c := &http.Cookie{
		Name:     "sessionid",
		Value:    sh.ID,
		HttpOnly: Config.GetWebConfig().Https,
		Expires:  sh.ExpirationTime,
	}

	Cookies.AddCookie(c, w, r)
}

/*
 * Disable the Caching on Local Machine For certain pages to make
 */
func (s *Instance) EnableCaching() {

}
