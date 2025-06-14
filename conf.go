package goconf

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	data map[string]interface{}
}

var config *Config

func LoadConfig(filename ...string) error {
	f := ".config.json"
	if len(filename) > 0 {
		f = filename[0]
	}

	bytes, err := os.ReadFile(f)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}

	config = &Config{data: raw}

	return nil
}

func (c *Config) getValue(key string) (interface{}, error) {
	parts := strings.Split(key, ".")
	var current interface{} = c.data

	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("path '%s' is not a valid object", part)
		}
		val, exists := m[part]
		if !exists {
			return nil, fmt.Errorf("key '%s' not found", key)
		}
		current = val
	}
	return current, nil
}

func GetField[T any](key string) (T, error) {
	var zero T
	val, err := config.getValue(key)
	if err != nil {
		return zero, err
	}

	converted, ok := convertToType[T](val)
	if !ok {
		return zero, fmt.Errorf("type assertion failed for key '%s'", key)
	}
	return converted, nil
}

func GetOpField[T any](key string, defaultValue T) T {
	val, err := GetField[T](key)
	if err != nil {
		return defaultValue
	}
	return val
}

func convertToType[T any](value interface{}) (T, bool) {
	var zero T
	targetType := reflect.TypeOf(zero)

	// Caso 1: ya es del tipo exacto
	if val, ok := value.(T); ok {
		return val, true
	}

	// Caso 2: el destino es struct, y el valor es un map[string]interface{}
	if targetType.Kind() == reflect.Struct {
		bytes, err := json.Marshal(value)
		if err != nil {
			return zero, false
		}
		var result T
		err = json.Unmarshal(bytes, &result)
		if err != nil {
			return zero, false
		}
		return result, true
	}

	// Caso 3: conversión directa
	val := reflect.ValueOf(value)
	if val.Type().ConvertibleTo(targetType) {
		return val.Convert(targetType).Interface().(T), true
	}

	// Caso 4: conversión de slices
	if targetType.Kind() == reflect.Slice && val.Kind() == reflect.Slice {
		slice := reflect.MakeSlice(targetType, val.Len(), val.Len())
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface()
			elemType := targetType.Elem()
			if reflect.TypeOf(item).ConvertibleTo(elemType) {
				slice.Index(i).Set(reflect.ValueOf(item).Convert(elemType))
			} else {
				return zero, false
			}
		}
		return slice.Interface().(T), true
	}

	// Caso 5: conversión de string
	if targetType.Kind() == reflect.String && reflect.TypeOf(value).Kind() == reflect.String {
		return value.(T), true
	}

	return zero, false
}
