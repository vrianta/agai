package component

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/vrianta/Server/model"
)

var (
	jsonStore        = make(map[string][]any) // map[table_name]jsonobj
	jsonStoreMu      sync.RWMutex
	componentsDir    = "./components"
	warnedMissingDir = false
)

func init() {
	LoadAllComponentsFromJSON()
}

// loadAllComponentsFromJSON loads all JSON files in ./components into jsonStore
func LoadAllComponentsFromJSON() {
	if _, err := os.Stat(componentsDir); os.IsNotExist(err) {
		if !warnedMissingDir {
			fmt.Printf("[Component] Warning: components directory '%s' does not exist.\n", componentsDir)
			warnedMissingDir = true
		}
		return
	}
	files, err := os.ReadDir(componentsDir)
	if err != nil {
		fmt.Printf("[Component] Error reading components directory: %v\n", err)
		return
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			path := filepath.Join(componentsDir, file.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("[Component] Error reading %s: %v\n", file.Name(), err)
				continue
			}
			var raw []any
			if err := json.Unmarshal(data, &raw); err != nil {
				fmt.Printf("[Component] Error unmarshaling %s: %v\n", file.Name(), err)
				continue
			}
			tableName := file.Name()[:len(file.Name())-len(".components.json")]
			jsonStore[tableName] = raw
		}
	}
}

// GetComponentJSON returns the raw JSON object for a table name
func GetComponentJSON(tableName string) (any, bool) {
	jsonStoreMu.RLock()
	defer jsonStoreMu.RUnlock()
	obj, ok := jsonStore[tableName]
	return obj, ok
}

// DumpComponentToJSON writes the in-memory component data to its JSON file
func DumpComponentToJSON(tableName string, data any) error {
	jsonStoreMu.Lock()
	defer jsonStoreMu.Unlock()
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(componentsDir, tableName+".components.json")
	return os.WriteFile(path, bytes, 0644)
}

// ReloadComponents reloads all JSON files from disk
func ReloadComponents() {
	LoadAllComponentsFromJSON()
}

// GetComponentMap returns the unmarshaled JSON as a slice of map[string]any for a table name
func GetComponentMap(tableName string) ([]map[string]any, bool) {
	jsonStoreMu.RLock()
	defer jsonStoreMu.RUnlock()
	obj, ok := jsonStore[tableName]
	if !ok {
		return nil, false
	}
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, false
	}
	var arr []map[string]any
	if err := json.Unmarshal(bytes, &arr); err != nil {
		return nil, false
	}
	return arr, true
}

// InitializeComponent syncs all JSON components with their DB tables.
// For each table in jsonStore:
//   - If the DB table is empty, insert all JSON values into the DB.
//   - If the DB table has data, load from DB, update jsonStore, and write to the JSON file.
func InitializeComponent() error {
	jsonStoreMu.Lock()
	defer jsonStoreMu.Unlock()

	if len(jsonStore) == 0 {
		fmt.Println("[Component] No components found to initialize.")
		return nil
	}
	fmt.Println("[Component] Initializing components...")
	for tableName, raw := range jsonStore {
		// Try to get the model by tableName (assumes a global registry or factory function)
		_model := getModelAndInserterByTableName(tableName)
		if _model == nil {
			panic("[ERROR] No Such Model Found for Table while creating the component: " + tableName)
		}

		// Check if DB is empty
		rows, err := _model.Get().Fetch()
		if err != nil {
			fmt.Printf("[Component] Error fetching from DB for '%s': %v\n", tableName, err)
			continue
		}
		if len(rows) == 0 {
			// Insert all JSON values into DB
			for key, value := range raw {
				fmt.Printf("[Component] Inserting into '%s': %s = %v\n", tableName, key, value)
				// Use the model's Create method to insert the data
				_model.Create().Set(key).To(value).Exec()
			}
		}
	}
	return nil
}

// getModelAndInserterByTableName returns the model and an insert function for a given table name
func getModelAndInserterByTableName(tableName string) *model.Struct {
	if m, ok := model.ModelsRegistry[tableName]; ok {
		return m
	}
	return nil
}
