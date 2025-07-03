package component

import (
	"reflect"
	"sync"

	"github.com/vrianta/Server/model"
)

type Component[T any, PK comparable] struct {
	Model         *model.Struct
	Val           map[PK]T
	DefaultValues []T
	primaryKey    string
	mu            sync.RWMutex
}

var (
	componentsMu sync.RWMutex
	components   []any // store all components for initialization
)

// New creates and registers a new component
func New[T any, PK comparable](model *model.Struct, primaryKey string, defaultValues ...T) *Component[T, PK] {
	c := &Component[T, PK]{
		Model:         model,
		Val:           make(map[PK]T),
		DefaultValues: defaultValues,
		primaryKey:    primaryKey,
	}
	componentsMu.Lock()
	components = append(components, c)
	componentsMu.Unlock()
	return c
}

// Initialize all registered components (should be called at startup)
func InitializeAll() {
	componentsMu.RLock()
	defer componentsMu.RUnlock()
	for _, comp := range components {
		switch c := comp.(type) {
		case interface{ Initialize() error }:
			_ = c.Initialize()
		}
	}
}

// Initialize checks DB, inserts defaults if needed, and loads data into Val
func (c *Component[T, PK]) Initialize() error {
	if c.Model == nil || !c.Model.Initialised {
		return nil // or trigger model init if desired
	}
	// Check if table is empty
	rows, err := c.Model.Get().Limit(1).Fetch()
	if err != nil {
		return err
	}
	if len(rows) == 0 && len(c.DefaultValues) > 0 {
		// Insert default values
		for _, def := range c.DefaultValues {
			q := c.Model.Create()
			v := reflect.ValueOf(def)
			for i := 0; i < v.NumField(); i++ {
				field := v.Type().Field(i)
				q.Set(field.Name).To(v.Field(i).Interface())
			}
			_ = q.Exec()
		}
	}
	// Fetch all data and populate Val
	legacyRows, err := c.Model.Get().Fetch()
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Val = make(map[PK]T, len(legacyRows))
	for _, row := range legacyRows {
		var item T
		rowMap := row.ToMap()
		v := reflect.ValueOf(&item).Elem()
		for k, val := range rowMap {
			f := v.FieldByName(k)
			if f.IsValid() && f.CanSet() {
				f.Set(reflect.ValueOf(val))
			}
		}
		pk := getPrimaryKey[PK, T](item, c.primaryKey)
		c.Val[pk] = item
	}
	return nil
}

// getPrimaryKey uses reflection to extract the primary key value from struct
func getPrimaryKey[PK comparable, T any](obj T, field string) PK {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.FieldByName(field).Interface().(PK)
}
