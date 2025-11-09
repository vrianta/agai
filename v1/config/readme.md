# Config Package Documentation

The `Config` package is the central configuration hub for your Go Server Framework project. It provides a flexible, extensible, and environment-aware way to manage all server, database, and web-related settings. This package ensures that your application can be easily configured for development, testing, and production environments.

## Features
- Centralized configuration management for server, database, and web settings
- JSON-based config files (e.g., `Config.json`, `Database.Config.json`)
- Environment variable overrides for all config values
- Type-safe Go structs for all config data
- Hot reloading support (can be extended)
- Multiple config types: server, database, web

## Directory Structure
```
Config/
├── Config.go            # Main config loader and manager
├── Database.Config.go   # Database config loader
├── Web.Config.go        # Web-specific config loader
├── type.go              # Type definitions for config structs
├── var.go               # Default values and global config variables
├── readme.md            # This documentation
```

## Main Server Configuration (`Config.go`)
- Loads and parses the main `Config.json` file from your project root.
- Defines core server settings: `Port`, `Host`, `Https`, `Build`, `StaticFolders`, `CssFolders`, `JsFolders`, `ViewFolder`, `MaxSessionCount`, `SessionStoreType`.
- Provides functions to load, validate, and access these settings.

### Example `Config.json`
```json
{
  "Port": "8080",
  "Host": "localhost",
  "Https": false,
  "Build": false,
  "StaticFolders": ["Static"],
  "CssFolders": ["Css"],
  "JsFolders": ["Js"],
  "ViewFolder": "Views",
  "MaxSessionCount": 1000,
  "SessionStoreType": "memory"
}
```

#### Supported Environment Variables for Server Config
- `SERVER_PORT`, `SERVER_HOST`, `SERVER_HTTPS`, `BUILD`, `STATIC_FOLDERS`, `CSS_FOLDERS`, `JS_FOLDERS`, `VIEW_FOLDER`, `MAX_SESSION_COUNT`, `SESSION_STORE_TYPE`

## Database Configuration (`Database.Config.go`)
- Loads and parses `Database.Config.json` from your project root.
- Defines all database connection settings: `Host`, `Port`, `User`, `Password`, `Database`, `Protocol`, `Driver`, `SSLMode`
- Supports environment variable overrides (e.g., `DB_HOST`, `DB_USER`, etc.)

### Example `Database.Config.json`
```json
{
  "Host": "localhost",
  "Port": "3306",
  "User": "root",
  "Password": "",
  "Database": "mydatabase",
  "Protocol": "tcp",
  "Driver": "mysql",
  "SSLMode": "disable"
}
```

#### Supported Environment Variables for Database Config
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_DATABASE`, `DB_PROTOCOL`, `DB_DRIVER`, `DB_SSLMODE`

## Web Configuration (`Web.Config.go`)
- Loads and parses web-specific settings (e.g., for web server, CORS, etc.)
- Can be extended for advanced web features (rate limiting, CORS, etc.)

## Type Definitions (`type.go`)
- Contains Go structs for all config types (server, database, web, etc.)
- Ensures type safety and IDE auto-completion

## Default Values and Globals (`var.go`)
- Provides default config values and global variables for use throughout the app

## Environment Variable Overrides
- All config values can be overridden by environment variables (see above for supported variables)
- This is essential for 12-factor app compatibility and cloud deployments

## Usage Example
```go
import "github.com/vrianta/agai/v1/config"

func main() {
    config := Config.Load() // Loads and parses Config.json
    fmt.Println("Server will run on:", config.Port)
}
```

## Best Practices
- Always commit example config files (e.g., `Config.example.json`) but never commit secrets
- Use environment variables for sensitive data in production
- Document all config options in your project README
- Validate config values at startup and fail fast if required values are missing

## License
MIT
