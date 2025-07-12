# Utils Package Documentation

The `utils` package provides general-purpose utility functions for the Go Server Framework. These helpers are used throughout the codebase for tasks such as file I/O, environment variable access, cryptographic operations, and password hashing.

## Features
- File reading helpers (e.g., `ReadFromFile`)
- Secure random token and session ID generation (`GenerateRandomToken`, `GenerateSessionID`)
- Password hashing and verification (`HashPassword`, `CheckPassword`)
- JSON encoding helpers (`JsonToString`)
- Environment variable access (`GetEnvString`)

## Example Usage
```go
import utils "github.com/vrianta/Server/v1/utils"

data := utils.ReadFromFile("config.json")
token, _ := utils.GenerateRandomToken("user123")
hash, _ := utils.HashPassword("mypassword")
```

---
