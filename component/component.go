package component

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	config "github.com/vrianta/Server/config"
	"github.com/vrianta/Server/model"
)

// Joson pattern will be
/*
{
 "primary_key": {
	// table components including the primary key
 }
}
*/

// map[string]map[string]any -> "[component_key/field_key value] => { "tableheading" : "value" } "
type component map[string]map[string]any

// [table_name](all the components)
type storage map[string]component

var (
	jsonStore        storage // store all the tables
	jsonStoreMu      sync.RWMutex
	componentsDir    = "./components"
	warnedMissingDir = false
)

// func init() {
// 	loadAllComponentsFromJSON()
// }

// ReloadComponents reloads all JSON files from disk
func ReloadComponents() {
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

	file_count := len(files)
	jsonStore = make(storage, file_count) // makign the map with fixed size less GC load

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			path := filepath.Join(componentsDir, file.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("[Component] Error reading %s: %v\n", file.Name(), err)
				continue
			}
			var raw component
			if err := json.Unmarshal(data, &raw); err != nil {
				if err.Error() == "json: cannot unmarshal array into Go value of type component.component" {
					panic("Make sure you component structure is properly structured for reference \n {\n \"Value of the PrimaryKey\":{\n  elemet: \"data\"\n}\n \n}\n and make sure you also have to include the primary key in the json object")
				}
				panic("[Component] Error unmarshaling " + file.Name() + " " + err.Error() + "\n")
			}
			// tableName := file.Name()[:len(file.Name())-len(".components.json")+1]
			tableName := strings.TrimSuffix(file.Name(), ".component.json")
			jsonStore[tableName] = raw
		} else {
			fmt.Println("[Warning] Non json file found in components - FileName: ", file.Name(), " Extension: ", filepath.Ext(file.Name()))
		}
	}

	initializeComponent()
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
	path := filepath.Join(componentsDir, tableName+".component.json")
	return os.WriteFile(path, bytes, 0644)
}

// GetComponentMap returns the unmarshaled JSON as a slice of map[string]any for a table name
func GetComponentMap(tableName string) (component, bool) {
	jsonStoreMu.RLock()
	defer jsonStoreMu.RUnlock()
	obj, ok := jsonStore[tableName]
	if !ok {
		return nil, false
	}
	// bytes, err := json.Marshal(obj)
	// if err != nil {
	// 	return nil, false
	// }
	// var arr []map[string]any
	// if err := json.Unmarshal(bytes, &arr); err != nil {
	// 	return nil, false
	// }
	// return arr, true

	return obj, true
}

// InitializeComponent syncs all JSON components with their DB tables.
// For each table in jsonStore:
//   - If the DB table is empty, insert all JSON values into the DB.
//   - If the DB table has data, load from DB, update jsonStore, and write to the JSON file.
func initializeComponent() error {

	if len(jsonStore) == 0 {
		fmt.Println("[Component] No components found to initialize.")
		return nil
	}
	fmt.Println("[Component] Initializing components...")
	for tableName, localList := range jsonStore {
		// Get the model for this table
		tableModel := getModelAndInserterByTableName(tableName)
		if tableModel == nil {
			panic("[ERROR] No Such Model Found for Table while creating the component: " + tableName)
		}

		// Get all rows from the database
		dbList, err := tableModel.Get().Fetch()
		if err != nil {
			fmt.Printf("[Component] Error fetching from DB for '%s': %v\n", tableName, err)
			continue
		}
		if len(dbList) == 0 {
			// If the database is empty, add everything from the local list
			for _, localItem := range localList {
				addRow := tableModel.Create()
				for key, value := range localItem {
					fmt.Printf("[Component] Inserting into '%s': %s = %v\n", tableName, key, value)
					addRow.Set(key).To(value)
				}
				if err := addRow.Exec(); err != nil {
					panic("[ERROR] - failed on component creation " + err.Error())
				}
			}
		} else {
			if config.GetBuild() {
				// localItem_key -> primary key value
				for localItem_key, localItem := range localList {
					if component, err := tableModel.Get().Where(dbList[0].GetPrimary().Name).Is(localItem_key).First(); err != nil {
						// means error on DB Connection or Failed to get the item
						panic(err.Error())
					} else if component == nil {
						// need to create this component in the DB with default value
						tableModel.Insert(localItem)
					} else {
						// no need to do anything
					}
				}
				// delete items from DB if it is present in the DB but not in local
				for _, dbModel := range dbList {
					// dbMap := dbModel.ToMap()
					if _, ok := localList[dbModel.GetPrimary().GetVal()]; !ok {
						// removed from the local file
						dbModel.Delete().Where(dbModel.GetPrimary().Name).Is(dbModel.GetPrimary().GetVal())
					}
				}
				// if it is build then we have to sync with the Database
				updatedLocalList := make(component, len(localList))
				for _, dbModel := range dbList {
					dbMap := dbModel.ToMap()
					updatedLocalList[fmt.Sprint(dbModel.GetPrimary().GetVal())] = dbMap
				}
				localList = updatedLocalList
				// At the end, update the localList in local storage
				if err := DumpComponentToJSON(tableName, updatedLocalList); err != nil {
					fmt.Printf("[Component] Error updating JSON file for '%s': %v\n", tableName, err)
				}
			} else {
				// sync the files with database
				// update the database according to the local file if the component is not present in the DB create it in DB
				// if the item is present in db but not in local component then delete it from db
				// basicaly sync it with DB

				// 1. Remove DB items not present in local
				for _, dbModel := range dbList {
					dbMap := dbModel.ToMap()
					found := false
					for _, localItem := range localList {
						match := true
						for key, value := range localItem {
							if dbVal, ok := dbMap[key]; !ok || dbVal != value {
								match = false
								break
							}
						}
						if match {
							found = true
							break
						}
					}
					if !found {
						// Not found in local, so delete from DB
						deleteQuery := dbModel.Delete()
						ifFirst := false
						for key, val := range dbMap {
							if ifFirst {
								deleteQuery.And()
							}
							deleteQuery.Where(key).Is(val)
							if !ifFirst {
								ifFirst = true
							}
						}
						if err := deleteQuery.Exec(); err != nil {
							fmt.Println("[Component] Error deleting DB item not in local:", err)
						}
					}
				}

				// 2. Add missing local items to DB
				for _, localItem := range localList {
					found := false
					for _, dbModel := range dbList {
						dbMap := dbModel.ToMap()
						match := true
						for key, value := range localItem {
							if dbVal, ok := dbMap[key]; !ok || dbVal != value {
								match = false
								break
							}
						}
						if match {
							found = true
							break
						}
					}
					if !found {
						// Not found in DB, so insert into DB
						addRow := tableModel.Create()
						for key, value := range localItem {
							addRow.Set(key).To(value)
						}
						if err := addRow.Exec(); err != nil {
							fmt.Println("[Component] Error inserting missing local item into DB:", err)
						}
					}
				}

			}

		}
	}
	return nil
}

// getModelAndInserterByTableName returns the model and an insert function for a given table name
func getModelAndInserterByTableName(tableName string) *model.Struct {
	if m, ok := model.ModelsRegistry[tableName]; ok {
		if !m.PrimaryKeyExists() {
			panic("[ERROR-" + tableName + "] Model need to have Primary Key if that is being used for components")
		}
		return m
	}

	fmt.Println(model.ModelsRegistry)
	panic("[ERROR] No Model Found with the Table Name: " + tableName)
}
