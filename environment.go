package environment

import (
	"fmt"
	"reflect"
	"sync"
)

// Environment provides thread-safe access to a key-value registry.
// Mutable values stored in the registry are not protected.
// An Environment must not be copied after first use.
type Environment struct {
	mu sync.RWMutex
	m  map[string]interface{}
}

var (
	instance *Environment
	once     sync.Once
)

// GetInstance returns the singleton instance of Environment.
func GetInstance() *Environment {
	once.Do(func() {
		instance = New()
	})
	return instance
}

// New creates an empty Environment. The zero value is also ready for use.
func New() *Environment {
	return &Environment{m: make(map[string]interface{})}
}

// Get retrieves a value or nil if the key does not exist.
// Use Lookup to distinguish an absent key from a stored nil.
func (e *Environment) Get(key string) interface{} {
	value, _ := e.Lookup(key)
	return value
}

// Lookup retrieves a value and reports whether the key exists.
func (e *Environment) Lookup(key string) (interface{}, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	value, ok := e.m[key]
	return value, ok
}

// Set stores a value, initializing the map if necessary.
func (e *Environment) Set(key string, value interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.m == nil {
		e.m = make(map[string]interface{})
	}
	e.m[key] = value
}

// Delete removes a key. Removing a missing key is safe.
func (e *Environment) Delete(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.m, key)
}

// Clone creates a shallow copy of the registry.
// Keys are independent, but reference-type values remain shared.
func (e *Environment) Clone() *Environment {
	e.mu.RLock()
	defer e.mu.RUnlock()
	clone := &Environment{m: make(map[string]interface{}, len(e.m))}
	for key, value := range e.m {
		clone.m[key] = value
	}
	return clone
}

// GetFieldValue returns an exported field from a struct or pointer to struct.
func GetFieldValue(config interface{}, fieldName string) (interface{}, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}
	v := reflect.ValueOf(config)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, fmt.Errorf("config is nil")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("config is not a struct")
	}
	// Resolve field metadata before accessing values: Value.FieldByName can
	// panic when a promoted field is reached through a nil embedded pointer.
	fieldInfo, ok := v.Type().FieldByName(fieldName)
	if !ok {
		return nil, fmt.Errorf("field %q not found", fieldName)
	}
	for _, index := range fieldInfo.Index {
		for v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return nil, fmt.Errorf("field %q is behind a nil embedded pointer", fieldName)
			}
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			return nil, fmt.Errorf("field %q cannot be accessed", fieldName)
		}
		v = v.Field(index)
	}
	if !v.CanInterface() {
		return nil, fmt.Errorf("field %q is not exported", fieldName)
	}
	return v.Interface(), nil
}
