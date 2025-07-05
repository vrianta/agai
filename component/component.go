package component

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
		fmt.Print("[INFO] To do migration of Components please use flag --migrate-component/-mc \n")
		RefreshComponentFromDB()
	}
}

// ReloadComponents reloads all JSON files from disk
func ReloadComponents() {
	loadComponents()
}

func loadComponents() {
	fmt.Print("---------------------------------------------------------\n")
	fmt.Println("[Component] - Loading Components from Local JSON")
	fmt.Print("---------------------------------------------------------\n")

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

	// Step 1: Filter valid component files (*.component.json)
	var componentFiles []os.DirEntry
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".component.json") {
			componentFiles = append(componentFiles, file)
		}
	}

	jsonStore = make(storage, len(componentFiles)) // preallocate only for component files
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, len(componentFiles))

	for _, file := range componentFiles {
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					errCh <- fmt.Errorf("panic while loading %s: %v", file.Name(), r)
				}
			}()

			path := filepath.Join(componentsDir, file.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				errCh <- fmt.Errorf("[Component] Error reading %s: %v", file.Name(), err)
				return
			}

			var raw components
			if err := json.Unmarshal(data, &raw); err != nil {
				if strings.Contains(err.Error(), "cannot unmarshal array into Go value of type component.component") {
					panic(fmt.Sprintf(`[Component] Malformed structure in %s. JSON should follow:
{
  "primaryKeyValue": {
    "field1": "value",
    "field2": "value"
  }
}
Ensure the primary key is present inside the nested object.`,
						file.Name()))
				}
				panic("[Component] Error unmarshaling " + file.Name() + ": " + err.Error())
			}

			tableName := strings.TrimSuffix(file.Name(), ".component.json")

			mu.Lock()
			jsonStore[tableName] = raw
			mu.Unlock()

			fmt.Println("Component for Table:", tableName, "is loaded")
		}(file)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		log.Println(err)
	}

	fmt.Print("---------------------------------------------------------\n")
	fmt.Println("[Component] - Done Loading Components from Local JSON")
	fmt.Print("---------------------------------------------------------\n\n")
}

// InitializeComponent syncs all JSON components with their DB tables.
// For each table in jsonStore:
//   - If the DB table is empty, InsertRow all JSON values into the DB.
//   - If the DB table has data, load from DB, update jsonStore, and write to the JSON file.
func syncComponent() error {
	if len(jsonStore) == 0 {
		fmt.Print("---------------------------------------------------------\n")
		fmt.Println("[Component] No components found to initialize.")
		fmt.Print("---------------------------------------------------------\n")
		return nil
	}

	fmt.Print("---------------------------------------------------------\n")
	fmt.Println("[Component] - Syncing Component with Database")
	fmt.Print("---------------------------------------------------------\n")

	var wg sync.WaitGroup
	errCh := make(chan error, len(jsonStore)) // buffered to avoid deadlocks

	for tableName, localList := range jsonStore {
		wg.Add(1)
		go func(tableName string, localList components) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					errCh <- fmt.Errorf("panic in table %s: %v", tableName, r)
				}
			}()

			fmt.Println("[Info] Initializing", tableName, "Component")
			tableModel := getModelAndInsertRowerByTableName(tableName)
			if tableModel == nil {
				errCh <- fmt.Errorf("[ERROR] No model found for table: %s", tableName)
				return
			}

			dbResults, err := tableModel.Get().Fetch()
			if err != nil {
				errCh <- fmt.Errorf("[Error] fetching from DB for %s: %w", tableName, err)
				return
			}

			if len(dbResults) == 0 {
				// DB is empty, InsertRow all local components
				for _, localItem := range localList {
					addRow := tableModel.Create()
					for key, value := range *localItem {
						fmt.Printf("\t[Info] InsertRowing into '%s': %s = %v\n", tableName, key, value)
						addRow.Set(key).To(value)
					}
					if err := addRow.Exec(); err != nil {
						errCh <- fmt.Errorf("[ERROR] failed to InsertRow component into %s: %w", tableName, err)
					}
				}
				return
			}

			fmt.Println("\t[INFO-Component] Syncing Component:", tableName)

			// Add new components from local file
			for localItemKey, localItem := range localList {
				if _, ok := dbResults[localItemKey]; !ok {
					fmt.Println("[INFO-COMPONENT] InsertRowing new Component", localItemKey, "in DB table:", tableName)
					dbResults[localItemKey] = model.Result(*localItem)
					if err := tableModel.InsertRow(*localItem); err != nil {
						errCh <- fmt.Errorf("[ERROR] failed to InsertRow new component %s into %s: %w", localItemKey, tableName, err)
					}
				}
			}

			// Remove DB entries not present in local
			for pk := range dbResults {
				if _, ok := localList[fmt.Sprint(pk)]; !ok {
					fmt.Printf("\t[INFO-COMPONENT] Deleting %s From %s\n", pk, tableName)
					if err := tableModel.Delete().Where(tableModel.GetPrimaryKey().Name()).Is(pk).Exec(); err != nil {
						errCh <- fmt.Errorf("[ERROR] failed to delete %v in table %s: %w", pk, tableName, err)
					}
				}
			}

			// Update local JSON file with DB state
			updatedLocalList := make(components, len(dbResults))
			for pk, dbItem := range dbResults {
				comp := component(dbItem)
				pkStr := fmt.Sprint(pk)
				updatedLocalList[pkStr] = &comp
			}
			if err := dumpComponentToJSON(tableName, updatedLocalList); err != nil {
				errCh <- fmt.Errorf("[Component] error updating JSON file for '%s': %v", tableName, err)
			}
		}(tableName, localList) // capture loop vars properly
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		log.Println(err)
	}

	fmt.Print("---------------------------------------------------------\n")
	fmt.Println("[Component] - Done Syncing Component with Database")
	fmt.Print("---------------------------------------------------------\n\n")
	return nil
}

func RefreshComponentFromDB() {

	jsonStore = make(storage, len(model.ModelsRegistry))
	errCh := make(chan error, len(model.ModelsRegistry))
	for tableName, tableModel := range model.ModelsRegistry {
		wb.Add(1)
		go func(tableName string, tableModel *model.Table) {
			defer wb.Done()

			fmt.Println("[Component] Loading from DB: ", tableName)

			if !tableModel.HasPrimaryKey() {
				panic("[Component] Skipping table %s: no primary key found" + tableName)
			}

			dbResults, err := tableModel.Get().Fetch()
			if err != nil {
				panic("[Component] Failed to fetch from DB for table " + tableName + " : " + err.Error())
			}

			updated := make(components, len(dbResults))

			for pk, row := range dbResults {
				comp := component(row)  // convert model.Result to component (map[string]any)
				pkStr := fmt.Sprint(pk) // convert primary key to string
				updated[pkStr] = &comp
			}

			// Update jsonStore directly
			jsonStoreMu.Lock()
			jsonStore[tableName] = updated
			jsonStoreMu.Unlock()

			// Optionally dump to file
			if err := dumpComponentToJSON(tableName, updated); err != nil {
				errCh <- fmt.Errorf("[Component] Failed to write %s.component.json: %v", tableName, err)
			}
		}(tableName, tableModel)
	}

	wb.Wait()
	close(errCh)
	for err := range errCh {
		log.Println(err)
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

// getModelAndInsertRowerByTableName returns the model and an InsertRow function for a given table name
func getModelAndInsertRowerByTableName(tableName string) *model.Table {
	if m, ok := model.ModelsRegistry[tableName]; ok {
		if !m.HasPrimaryKey() {
			panic("[ERROR-" + tableName + "] Model need to have Primary Key if that is being used for components")
		}
		return m
	}

	fmt.Println(model.ModelsRegistry)
	panic("[ERROR] No Model Found with the Table Name: " + tableName)
}

func Get(model *model.Table) *components {
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
