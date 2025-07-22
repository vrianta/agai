# Session management

## Chapter 1: Introduction to Sessions in AGAI

In the stateless world of HTTP, maintaining a user's identity, preferences, and interactions between requests is essential for building meaningful web applications. This is achieved through **sessions** — structures that bridge multiple HTTP requests into a single logical conversation with the user.

AGAI’s session management system implements a **custom, memory-resident, disk-persistent, LRU-aware** session store tailored for small to medium-sized Go web applications.

This guide walks you through every major concept, internal component, and lifecycle involved in the system — not as a reference, but as a walkthrough, explaining not just _what_ the code does, but _why_.

## Chapter 2: The Philosophy of Session Design

### 2.1: Why Not Use Gorilla or Third-Party Tools?

AGAI aims for minimalism and control. Instead of depending on heavy external libraries like Gorilla Sessions or Redis-backed tools, this session manager is:

*   Embedded in the app's memory.
    
*   Persistent via disk serialization (`gob`).
    
*   Concurrent-safe.
    
*   Tailored to work closely with AGAI’s request/response lifecycles.
    

This design is motivated by a balance between **performance**, **simplicity**, and **extensibility**.

## Chapter 3: Anatomy of a Session

### 3.1: The `Instance` Struct

At the heart of the session system lies the `Instance` struct. This is the representation of a single user’s state:

```go
type Instance struct { 	
    ID               string 	
    PostParameters   HTTPPostParameters 	
    GetParameters    HTTPGetParameters 	
    IsAuthenticated  bool 	
    ExpirationTime   time.Time 	
    Data             SessionData 
}
```

Think of it like a browser tab’s memory — each session:

*   Has a **unique ID** (used in cookies).
    
*   Stores GET/POST data parsed from the request.
    
*   Indicates whether the user is **authenticated**.
    
*   Contains a `Data` map — the heart of per-user storage.
    
*   Tracks **expiration** time for cleanup.

## Chapter 4: Creating and Starting Sessions

### 4.1: New Sessions

`New()` creates an empty session, setting `"uid": "Guest"` as a default identity and assigning a short expiration window.

This is a **raw** session — not yet attached to a request/response cycle.

### 4.2: Getting the Session ID from a Request

To resume a previous session, you extract the session ID from cookies:


```go 
GetSessionID(r)
```

If the cookie `"sessionid"` exists, it returns a string pointer to its value.

### 4.3: Starting a Session

Sessions are fully initialized using:

```go
StartSession(sessionID, w, r)
```

This:

1.  Binds the session to the current request.
    
2.  Creates a cookie.
    
3.  Registers the session in memory.
    
4.  Schedules it for future expiration.

## Chapter 5: The Lifecycle of a Session

### 5.1: Birth — via `CreateNewSession`

When a session starts, `CreateNewSession`:

*   Sets the session ID.
    
*   Pushes it into the expiration heap and LRU list.
    
*   Saves it to disk if needed.
    
*   Attaches a cookie.
    

### 5.2: Growth — via `Store()`

Calling `Store(sessionID, session)` is how you commit a session to memory.

Here, **thread safety is guaranteed**:

*   Locking guards against race conditions.
    
*   The session is inserted into `instances` (main map).
    
*   If memory is full, LRU eviction kicks in.
    

### 5.3: Adulthood — Active Use

During its lifetime, the session holds form data, flags, authentication state, and more. You manipulate it directly via the `Data` map or helper methods like `Login()`.

You can clear volatile form data via `Clean()` — a helpful post-processing call after handling forms.

### 5.4: Death — via Expiry or Manual Removal

When the session is too old (`ExpirationTime`) or explicitly removed via `RemoveSession(sessionID)`, it's purged from memory and optionally from disk.

## Chapter 6: Under the Hood — Memory Management

### 6.1: `instances` Map

A global in-memory map:

```go
instances map[string]*Instance
```

This is the master list of sessions.

### 6.2: LRU Management

To avoid memory leaks in long-running servers, the system uses an **LRU eviction policy**:

*   A doubly-linked list (`lruOrderList`) stores usage order.
    
*   When the session is used, it’s moved to the front.
    
*   When `instances` exceeds the configured limit, the last item is removed.
    

### 6.3: Expiration Heap

Parallel to LRU is a **min-heap**:

```go
sessionHeap []*Instance
```

Each session is pushed into this heap based on `ExpirationTime`. A background goroutine checks the earliest expiring session and sleeps until it needs to be removed.

This dual system ensures both **activity-based** and **time-based** cleanup.

* * *

## Chapter 7: Concurrency and Safety

To prevent data races, the system uses several mutexes:

*   `sessionStoreMutex`: Protects the global `instances` map.
    
*   `heapAccessMutex`: Synchronizes heap access.
    
*   `sessionCleanMutex`, `sessionUpdateMutex`: Ensure GET/POST cleanup and updates don’t overlap.
    
*   `lruCacheMutex`: Protects LRU list and map.
    

Goroutines (`StartSessionHandler`, `StartLRUHandler`) are **always safe** to run in the background indefinitely.

* * *

## Chapter 8: Persistence with Gob

Sessions are persisted using Go’s `encoding/gob`:

### 8.1: Saving Sessions

```go
saveAllSessionsToDisk()
```

*   Serializes all active sessions to `sessions.data`.
    
*   Triggered after most write operations.
    

### 8.2: Loading Sessions

```go
loadAllSessionsFromDisk()
```

*   Reads `sessions.data` at startup.
    
*   Populates `instances` with deserialized sessions.
    

All necessary types (`map[string]*Instance`, `*Instance`, `atomic.Int64`) are registered with `gob` in `init()`.

* * *

## Chapter 9: Goroutines and Channels

Two key goroutines maintain health:

### 9.1: Expiry Manager

```go
StartSessionHandler()
```

This loop:

1.  Peeks at the earliest session in the heap.
    
2.  Sleeps until it should expire.
    
3.  Removes it, logs, and repeats.
    

It’s efficient — only wakes when needed.

### 9.2: LRU Manager
```go
StartLRUHandler()
```

Listens to a channel for operations like `"move"` or `"InsertRow"`, and updates the LRU list accordingly.

* * *

## Chapter 10: How to Use the System

You don’t need to understand all internals to use sessions effectively.

Here’s a canonical handler:

```go
func HomeHandler(w http.ResponseWriter, r *http.Request) {     
    sid := session.GetSessionID(r)     
    s := session.New()     
    sid = s.StartSession(sid, w, r)     
    session.Store(sid, s)      
    if s.IsLoggedIn() {         
        fmt.Fprintln(w, "Welcome,", s.Data["uid"])     
    } else {         
        fmt.Fprintln(w, "Please log in.")     
    }      
    s.Clean() 
}
```

* * *

## Chapter 11: Extensions and Customization Ideas

You could extend this system with:

*   Redis or in-memory distributed backends.
    
*   IP-based session tracking.
    
*   Per-session flash messages.
    
*   JWT hybrid sessions (store metadata in cookies).
    
*   User-agent or fingerprint binding for extra security.
    

* * *

## Chapter 12: Closing Thoughts

The AGAI session system is simple by design but surprisingly powerful under the hood.

*   Sessions are compact and fast.
    
*   Memory is tightly managed.
    
*   Disk persistence ensures durability.
    
*   Background workers keep it clean.
    

If you're building with AGAI, this session manager gives you everything you need — no black boxes, no magic. Just clear, extensible Go code.