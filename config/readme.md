# Config Package Documentation

The `Config` package is the central configuration hub for your Go Server Framework project. It provides a flexible, extensible, and environment-aware way to manage all server, database, and web-related settings. This package ensures that your application can be easily configured for development, testing, and production environments.

---

## Features
- **Centralized Configuration Management**: All server, database, and web settings are managed in one place.
- **JSON-Based Config Files**: Uses human-readable JSON files for configuration.
- **Environment Variable Overrides**: Supports overriding config values with environment variables for CI/CD and production. All config values can be set via environment variables, which take precedence over file values. See the list of supported variables below.
- **Type-Safe Access**: Strongly-typed Go structs for all config data.
- **Hot Reloading (optional)**: Can be extended to support config reloads without restarting the server.
- **Multiple Config Types**: Separate config files for server, database, and web settings.

---

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

---

## 1. Main Server Configuration (`Config.go`)

- Loads and parses the main `Config.json` file from your project root.
- Defines the core server settings, such as:
  - `Port`: HTTP/HTTPS port
  - `Host`: Hostname or IP
  - `Https`: Enable/disable HTTPS
  - `Build`: Build mode (affects template caching, migration, etc.)
  - `StaticFolders`, `CssFolders`, `JsFolders`, `ViewFolder`: Folder locations
  - `MaxSessionCount`: Session limits
  - `SessionStoreType`: Session backend (memory, redis, etc.)
- Provides functions to load, validate, and access these settings throughout your application.

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
- `SERVER_PORT`
- `SERVER_HOST`
- `SERVER_HTTPS`
- `BUILD`
- `STATIC_FOLDERS` (comma-separated)
- `CSS_FOLDERS` (comma-separated)
- `JS_FOLDERS` (comma-separated)
- `VIEW_FOLDER`
- `MAX_SESSION_COUNT`
- `SESSION_STORE_TYPE`

Environment variables take precedence over values in `Config.json`.

---

## 2. Database Configuration (`Database.Config.go`)

- Loads and parses `Database.Config.json` from your project root.
- Defines all database connection settings:
  - `Host`, `Port`, `User`, `Password`, `Database`, `Protocol`, `Driver`, `SSLMode`
- Supports environment variable overrides (e.g., `DB_HOST`, `DB_USER`, etc.).
- Used by the database handler to establish and manage DB connections.

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
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_DATABASE`
- `DB_PROTOCOL`
- `DB_DRIVER`
- `DB_SSLMODE`

Environment variables take precedence over values in `Database.Config.json`.

---

## 3. Web Configuration (`Web.Config.go`)

- Loads and parses web-specific settings (e.g., for web server, CORS, etc.).
- Can be extended for advanced web features (rate limiting, CORS, etc.).

---

## 4. Type Definitions (`type.go`)

- Contains Go structs for all config types (server, database, web, etc.).
- Ensures type safety and IDE auto-completion.
- Example:
  ```go
  type ServerConfig struct {
      Port             string
      Host             string
      Https            bool
      Build            bool
      StaticFolders    []string
      CssFolders       []string
      JsFolders        []string
      ViewFolder       string
      MaxSessionCount  int
      SessionStoreType string
  }
  ```

---

## 5. Default Values and Globals (`var.go`)

- Provides default config values and global variables for use throughout the app.
- Useful for fallback values and for sharing config state between packages.

---

## 6. Environment Variable Overrides

- All config values can be overridden by environment variables (see above for supported variables).
- This is essential for 12-factor app compatibility and cloud deployments.
- The config loader checks for environment variables before falling back to JSON file values.

---

## 7. Usage Example

```go
import "github.com/vrianta/Server/Config"

func main() {
    config := Config.Load() // Loads and parses Config.json
    fmt.Println("Server will run on:", config.Port)
}
```

---

## 8. Best Practices
- Always commit example config files (e.g., `Config.example.json`) but never commit secrets.
- Use environment variables for sensitive data in production.
- Document all config options in your project README.
- Validate config values at startup and fail fast if required values are missing.

---

## 9. Extending the Config Package
- Add new FieldTypes to the config structs in `type.go` as needed.
- Update the loader functions to parse new FieldTypes.
- Add new config files (e.g., `Email.Config.json`) for additional features.

---

## 10. Troubleshooting
- If config values are not loading, check for typos in file names and field names.
- Use environment variables to debug and override values quickly.
- Log the loaded config at startup for verification.

---

## 11. License
MIT
