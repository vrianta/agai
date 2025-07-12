# Log Package Documentation

The `log` package provides logging utilities for the Go Server Framework. It is used for error reporting, debugging, and informational output throughout the application.

## Features
- Simple logging functions (`WriteLog`, `WriteLogf`)
- Conditional logging based on build mode
- JSON response formatting for API errors

## Example Usage
```go
import log "github.com/vrianta/Server/v1/log"

log.WriteLog("Server started")
log.WriteLogf("User %s logged in", username)
```

---
