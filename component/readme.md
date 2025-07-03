# Component Package Documentation

The `component` package provides a generic, type-safe, and extensible way to manage application components that are backed by database model. It is designed for settings, configuration tables, or any static/dynamic data that should be loaded, initialized, and accessed efficiently in memory.

## Features
- Generic: Works with any struct type and primary key.
- Type-safe: Access your data as `map[PrimaryKey]YourStruct`.
- Automatic initialization: Loads from DB, inserts defaults if needed.
- Thread-safe: Uses mutexes for concurrent access.
- Centralized registration and initialization of all components.

## Usage

### 1. Define Your Model and Struct
```go
// Example struct for a settings table
 type Setting struct {
     Key   string
     Value string
 }

// Register your model (using your model package)
var SettingsModel = model.New("settings", map[string]model.Field{
    "Key":   {/* ... */},
    "Value": {/* ... */},
})
```

### 2. Register a Component
```go
import "github.com/vrianta/Server/component"

var SettingsComponent = component.New[Setting, string](
    SettingsModel, // your *model.Struct
    "Key",         // primary key field name
    Setting{Key: "site_name", Value: "My Site"}, // default values (optional)
    Setting{Key: "theme", Value: "light"},
)
```

### 3. Initialize All Components (at startup)
```go
component.InitializeAll()
```

### 4. Access Data
```go
siteName := SettingsComponent.Val["site_name"].Value
```

## How It Works
- On initialization, each component checks if its table is empty. If so, it inserts the provided default values.
- All rows are loaded from the DB and stored in a map, keyed by the primary key field you specify.
- You can safely read from the `Val` map concurrently.

## Advanced
- You can register as many components as you want. All are initialized with `component.InitializeAll()`.
- The system uses reflection to map DB rows to your struct and extract the primary key.
- If you need to reload data, call `YourComponent.Initialize()` again.

## Example: Settings Component
```go
// Define your struct and model
 type Setting struct {
     Key   string
     Value string
 }
var SettingsModel = model.New("settings", map[string]model.Field{
    "Key":   {/* ... */},
    "Value": {/* ... */},
})

// Register the component
var SettingsComponent = component.New[Setting, string](SettingsModel, "Key",
    Setting{Key: "site_name", Value: "My Site"},
    Setting{Key: "theme", Value: "light"},
)

// Initialize all components at startup
tfunc main() {
    component.InitializeAll()
    // Now you can use SettingsComponent.Val["site_name"].Value
}
```

## Best Practices
- Always call `component.InitializeAll()` before accessing component data.
- Use struct field names as the primary key (case-sensitive, must match your struct).
- Use default values to ensure your tables are never empty.
- For advanced use, you can extend the `Component` struct with custom methods or hooks.

---

component.New[T any](m *model, default_values) component {
    // create new controller
    // default values are all the values for the component elements I am not sure which should be the best aproach to get the data because the model can have two or more fields
    create the struct of component which will store the 
    {
        model,
        val T
        default_values
    }

    // store the components in a storage
    return component
}

initialised() {
    loop through all the components and check if the model is initialsed and then 

    check if the table has no content then we should create elements with default value

    else we should get the data and store it in val using reflect 
}

the T probably map[primarykey]struct{elements}