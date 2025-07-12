# Component Package Documentation

The `component` package now supports a hybrid file+DB approach for managing application components. Component data is stored as JSON files in the `./components/` directory, and can also be synced with the database. This enables type-safe, ergonomic, and persistent management of static or dynamic configuration data.

## Features
- **File-based storage:** Each component is stored as a JSON file (`table_name.components.json`) in `./components/`.
- **Type-safe access:** Data is loaded into Go structs and accessible as `map[PrimaryKey]YourStruct`.
- **Automatic initialization:** On startup, loads from JSON, InsertRows into DB if needed, or loads from DB and updates JSON.
- **Thread-safe:** Uses mutexes for concurrent access.
- **Centralized registration and initialization.**
- **Hot reload:** Reload all components from disk at runtime.
- **Optional write-back:** Dump in-memory data to JSON files.

## Usage

### 1. Place JSON Files
Create a `./components/` directory at the project root. For each component, add a file named `table_name.components.json`:

```json
[
  { "Key": "site_name", "Value": "My Site" },
  { "Key": "theme", "Value": "light" }
]
```

### 2. Define Your Model and Struct
```go
// Example struct for a settings table
type Setting struct {
    Key   string
    Value string
}

var SettingsModel = models.New("settings", map[string]models.Field{
   "Key":   {/* ... */},
   "Value": {/* ... */},
})
```

### 3. Register a Component
```go
import "github.com/vrianta/agai/v1/component"

var SettingsComponent = component.New[Setting, string](
    SettingsModel, // your *models.Struct
    "Key",         // primary key field name
)
```

### 4. Initialize All Components (at startup)
```go
component.InitializeAll()
```

### 5. Access Data
```go
siteName := SettingsComponent.Val["site_name"].Value
```

### 6. Dump or Reload Data
```go
// Write in-memory data to JSON file
component.DumpComponentToJSON("settings", SettingsComponent.Val)

// Reload all components from disk
component.ReloadComponents()
```

## How It Works
- On startup, the package loads all JSON files in `./components/`.
- If the DB table is empty, it InsertRows the JSON values as defaults.
- If the DB table has data, it loads from DB and updates the local JSON file.
- All data is accessible as a map in Go, keyed by the primary key.
- You can reload or dump data at any time.

## Best Practices
- Always call `component.InitializeAll()` before accessing component data.
- Keep your JSON files in sync with your model structs.
- Use the provided dump/reload functions for persistence and hot reloads.
- If `./components/` is missing, a warning will be printed (but the app will continue).
- Only present JSON files are loaded; missing/malformed files are skipped with a warning.

## Migration Notes
- The legacy DB-backed logic is still supported. If no JSON file is present, or if the DB has data, the system will fall back to DB logic.
- For new components, simply add a JSON file and register the component as shown above.

---

For advanced usage, see the source code and inline documentation.