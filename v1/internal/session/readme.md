# Session Package Documentation

The `internal/session` package provides a robust session management system for the Go Server Framework. It handles user authentication, session creation, secure cookie handling, and session data persistence.

## Features
- Session creation and management (in-memory or disk-based)
- Secure, random session ID generation
- LRU-based session heap for efficient cleanup
- Thread-safe access using mutexes
- Session expiry and automatic cleanup
- User login/logout helpers
- Request parsing for GET and POST data
- Session data storage via the `Store` map

## Usage
```go
import session "github.com/vrianta/Server/v1/internal/session"

sess := session.New()
sessID := sess.StartSession(nil, w, r)
if sess.IsLoggedIn() {
    // user is authenticated
}
```

## Best Practices
- Use HTTPS to secure cookies
- Regularly invalidate stale sessions
- Use environment variables to configure session storage type and limits

---
