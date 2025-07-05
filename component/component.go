package component

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vrianta/Server/config"
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

func Init() {
	if config.SyncComponentsEnabled {
		loadComponents()
		syncComponent()
		loadComponents()
	} else {
		RefreshComponentFromDB()
	}
}

// ReloadComponents reloads all JSON files from disk
func ReloadComponents() {
	loadComponents()
}

// loadAllComponentsFromJSON loads all JSON files in ./components into jsonStore
func loadComponents() {

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
			var raw components
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
			// fmt.Println("[Warning] Non json file found in components - FileName: ", file.Name(), " Extension: ", filepath.Ext(file.Name()))
		}
	}

}

// InitializeComponent syncs all JSON components with their DB tables.
// For each table in jsonStore:
//   - If the DB table is empty, insert all JSON values into the DB.
//   - If the DB table has data, load from DB, update jsonStore, and write to the JSON file.
func syncComponent() error {
	if len(jsonStore) == 0 {
		fmt.Println("[Component] No components found to initialize.")
		return nil
	}

	for tableName := range jsonStore {
		localList := jsonStore[tableName]
		fmt.Println("[Component] Initializing ", tableName, " components...")
		tableModel := getModelAndInserterByTableName(tableName)
		if tableModel == nil {
			panic("[ERROR] No Such Model Found for Table while creating the component: " + tableName)
		}

		dbResults, err := tableModel.Get().Fetch()
		if err != nil {
			fmt.Printf("[Component] Error fetching from DB for '%s': %v\n", tableName, err)
			continue
		}
		// dbResults.PrintAsTable()

		if len(dbResults) == 0 {
			// DB is empty, insert all local components
			for _, localItem := range localList {
				addRow := tableModel.Create()
				for key, value := range *localItem {
					fmt.Printf("[Component] Inserting into '%s': %s = %v\n", tableName, key, value)
					addRow.Set(key).To(value)
				}
				if err := addRow.Exec(); err != nil {
					panic("[ERROR] - failed on component creation " + err.Error())
				}
			}
			continue
		}

		fmt.Println("[INFO-Component] Syncing Component: ", tableName)

		// Add new components from local file
		for localItemKey, localItem := range localList {
			if _, ok := dbResults[localItemKey]; !ok {
				fmt.Println("[INFO-COMPONENT] Inserting new Component ", localItemKey, " in DB table: ", tableName)
				dbResults[localItemKey] = model.Result(*localItem)
				tableModel.Insert(*localItem)
			}
		}

		// Remove DB entries not present in local
		for pk := range dbResults {
			if _, ok := localList[fmt.Sprint(pk)]; !ok {
				fmt.Printf("[INFO-COMPONENT] Deleting %s From %s\n", pk, tableName)
				println(localList)
				if err := tableModel.Delete().Where(tableModel.GetPrimaryKey().Name).Is(pk).Exec(); err != nil {
					panic("[ERROR] - Failed to Delete " + fmt.Sprint(pk) + " in table " + tableName)
				}
			}
		}

		// Update local JSON file with DB state
		updatedLocalList := make(components, len(dbResults))
		for pk, dbItem := range dbResults {
			comp := component(dbItem)
			i := fmt.Sprint(pk)
			updatedLocalList[i] = &comp
		}
		if err := dumpComponentToJSON(tableName, updatedLocalList); err != nil {
			fmt.Printf("[Component] Error updating JSON file for '%s': %v\n", tableName, err)
		}

	}
	return nil
}

func RefreshComponentFromDB() {
	jsonStoreMu.Lock()
	defer jsonStoreMu.Unlock()

	for tableName, tableModel := range model.ModelsRegistry {
		fmt.Println("[Component] Refreshing from DB: ", tableName)

		if !tableModel.PrimaryKeyExists() {
			fmt.Printf("[Component] Skipping table %s: no primary key found.\n", tableName)
			continue
		}

		dbResults, err := tableModel.Get().Fetch()
		if err != nil {
			fmt.Printf("[Component] Failed to fetch from DB for table %s: %v\n", tableName, err)
			continue
		}

		updated := make(components, len(dbResults))

		for pk, row := range dbResults {
			comp := component(row)  // convert model.Result to component (map[string]any)
			pkStr := fmt.Sprint(pk) // convert primary key to string
			updated[pkStr] = &comp
		}

		// Update jsonStore directly
		jsonStore[tableName] = updated

		// Optionally dump to file
		if err := dumpComponentToJSON(tableName, updated); err != nil {
			fmt.Printf("[Component] Failed to write %s.component.json: %v\n", tableName, err)
		}
	}
}

// dumpComponentToJSON writes the in-memory component data to its JSON file
func dumpComponentToJSON(tableName string, data any) error {
	jsonStoreMu.Lock()
	defer jsonStoreMu.Unlock()
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(componentsDir, tableName+".component.json")
	return os.WriteFile(path, bytes, 0644)
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

func Get(model *model.Struct) *components {
	if model == nil {
		panic("[Component] GetComponentFromModel: modelStruct is nil")
	}

	data, ok := jsonStore[model.TableName]
	if !ok {
		panic("[Component] No component data found in jsonStore for table: " + model.TableName)
	}

	return &data
}

func (c *components) Where(id string) *component {
	return (*c)[id]
}

func (c *components) Is(field string, value any) components {
	result := make(components)
	for id, comp := range *c {
		if v, ok := (*comp)[field]; ok && v == value {
			result[id] = comp
		}
	}
	return result
}
