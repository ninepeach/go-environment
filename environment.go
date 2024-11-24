package environment

import (
	"fmt"
	"reflect"
	"sync"
)

// Environment is a thread-safe key-value store.
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

func New() *Environment {
	return &Environment{
		m: make(map[string]interface{}),
	}
}

// Get retrieves the value associated with the given key.
// Returns nil if the key does not exist.
func (e *Environment) Get(key string) interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.m[key]
}

// Set stores the value associated with the given key.
func (e *Environment) Set(key string, value interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.m[key] = value
}

// Clone creates a copy of the current Environment.
func (e *Environment) Clone() *Environment {
	e.mu.RLock()
	defer e.mu.RUnlock()

	clone := &Environment{
		m: make(map[string]interface{}, len(e.m)),
	}
	for key, value := range e.m {
		clone.m[key] = value
	}
	return clone
}

func GetFieldValue(config interface{}, fieldName string) (interface{}, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	v := reflect.ValueOf(config)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("config is not a struct")
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field '%s' not found", fieldName)
	}

	return field.Interface(), nil
}
