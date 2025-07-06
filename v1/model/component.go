package model

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Joson pattern will be
/*
{
 "primary_key": {
	// table components including the primary key
 }
}
*/
// Loads a component JSON file and stores it in the model's components map
func (m *meta) loadComponentFromDisk() {
	fmt.Printf("[component] Loading component for table: %s\n", m.TableName)
	path := filepath.Join(componentsDir, m.TableName+".component.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("[component] No JSON file found for: %s\n", m.TableName)
		m.components = make(components)
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("[component] Error reading %s: %v\n", path, err)
		return
	}

	var raw components
	if err := json.Unmarshal(data, &raw); err != nil {
		log.Printf("[component] Error unmarshaling %s: %v\n", path, err)
		return
	}

	m.components = raw
	fmt.Println("components: ", m.components)
	fmt.Printf("[component] Component for table '%s' loaded with %d items\n", m.TableName, len(raw))
}

// Saves the model's in-memory components to its JSON file
func (m *meta) saveComponentToDisk() error {
	path := filepath.Join(componentsDir, m.TableName+".component.json")
	bytes, err := json.MarshalIndent(m.components, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

// Syncs the in-memory components with the DB table for this model
func (m *meta) syncComponentWithDB() error {
	if len(m.components) == 0 {
		fmt.Printf("[component] No components to sync for: %s\n", m.TableName)
		return nil
	}

	dbResults, err := m.Get().Fetch()
	if err != nil {
		return fmt.Errorf("[component] fetch error for %s: %w", m.TableName, err)
	}

	if len(dbResults) == 0 {
		for _, localItem := range m.components {
			if err := m.InsertRow(localItem); err != nil {
				log.Printf("[component] Insert failed: %v\n", err)
			}
		}
		return nil
	}

	// Add missing
	for k, v := range m.components {
		if _, ok := dbResults[k]; !ok {
			log.Printf("[component] DB missing component '%s' in table %s\n", k, m.TableName)
			_ = m.InsertRow(v)
		}
	}

	// Remove stale
	for k := range dbResults {
		if _, ok := m.components[fmt.Sprint(k)]; !ok {
			_ = m.Delete().Where(m.GetPrimaryKey().Name()).Is(k).Exec()
		}
	}

	// Update component file with DB contents
	updated := make(components)
	for k, v := range dbResults {
		c := component(v)
		updated[fmt.Sprint(k)] = c
	}
	m.components = updated

	return m.saveComponentToDisk()
}

// Refreshes model's in-memory components from DB and rewrites JSON
func (m *meta) refreshComponentFromDB() {
	if !m.HasPrimaryKey() {
		log.Printf("[component] Model %s missing primary key\n", m.TableName)
		return
	}
	results, err := m.Get().Fetch()
	if err != nil {
		log.Printf("[component] Fetch error for %s: %v\n", m.TableName, err)
		return
	}
	updated := make(components)
	for k, v := range results {
		c := component(v)
		updated[fmt.Sprint(k)] = c
	}
	m.components = updated
	_ = m.saveComponentToDisk()
}

func (m *meta) GetComponents() components {
	return m.components
}

func (m *meta) GetComponent(id string) component {
	return m.components[id]
}

func (c component) FieldValue(field string) any {
	return c[field]
}
