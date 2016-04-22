package golf

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"reflect"
	"strings"
)

// KeyError is thrown when the specified key is not found in the configuration.
type KeyError struct {
	key string
}

// Error method implements Error method of Go standard library "error".
func (err *KeyError) Error() string {
	return fmt.Sprintf("Key not found, key: %s", err.key)
}

// ValueTypeError is thrown when the type of the specified value is not valid.
type ValueTypeError struct {
	key     string
	value   interface{}
	message string
}

// Error method implements Error method of Go standard library "error".
func (err *ValueTypeError) Error() string {
	return fmt.Sprintf("%s, key: %s, value: %v (%s)", err.message, err.key, err.value, reflect.TypeOf(err.value).Name())
}

// Config control for the application.
type Config struct {
	mapping map[string]interface{}
}

// NewConfig creates a new configuration instance.
func NewConfig() *Config {
	mapping := make(map[string]interface{})
	return &Config{mapping}
}

// GetString fetches the string value by indicating the key.
// It returns a ValueTypeError if the value is not a sring.
func (config *Config) GetString(key string, defaultValue string) (string, error) {
	value, err := config.Get(key, defaultValue)
	if err != nil {
		return defaultValue, err
	}
	if result, ok := value.(string); ok {
		return result, nil
	}
	return defaultValue, &ValueTypeError{key: key, value: value, message: "Value is not a string."}
}

// GetInt fetches the int value by indicating the key.
// It returns a ValueTypeError if the value is not a sring.
func (config *Config) GetInt(key string, defaultValue int) (int, error) {
	value, err := config.Get(key, defaultValue)
	if err != nil {
		return defaultValue, err
	}
	if result, ok := value.(int); ok {
		return result, nil
	}
	return defaultValue, &ValueTypeError{key: key, value: value, message: "Value is not an integer."}
}

// GetBool fetches the bool value by indicating the key.
// It returns a ValueTypeError if the value is not a sring.
func (config *Config) GetBool(key string, defaultValue bool) (bool, error) {
	value, err := config.Get(key, defaultValue)
	if err != nil {
		return defaultValue, err
	}
	if result, ok := value.(bool); ok {
		return result, nil
	}
	return defaultValue, &ValueTypeError{key: key, value: value, message: "Value is not an bool."}
}

// GetFloat fetches the float value by indicating the key.
// It returns a ValueTypeError if the value is not a sring.
func (config *Config) GetFloat(key string, defaultValue float64) (float64, error) {
	value, err := config.Get(key, defaultValue)
	if err != nil {
		return defaultValue, err
	}
	if result, ok := value.(float64); ok {
		return result, nil
	}
	return defaultValue, &ValueTypeError{key: key, value: value, message: "Value is not a float."}
}

// Set is used to set the value by indicating the key.
// If you want to set multi-level json, key can be like 'foo/bar'.
// For instance, `Set("foo/bar", 4)` and `Set("foo/bar2", "foo")`.
// If the parent key is not a map, then return a KeyError.
// For instance, can not set ("foo/bar", 4) after setting ("foo", 5).
func (config *Config) Set(key string, value interface{}) error {
	var tmp interface{}
	keys := strings.Split(key, "/")
	tmp = config.mapping
	for i, item := range keys {
		if len(item) == 0 {
			continue
		}
		if mapping, ok := tmp.(map[string]interface{}); ok {
			if i == len(keys)-1 {
				mapping[item] = value
				return nil
			}
			if value, exists := mapping[item]; exists {
				switch t := value.(type) {
				case map[string]interface{}:
					tmp = value
				default:
					_ = t
					mapping[item] = make(map[string]interface{})
					tmp = mapping[item]
				}
			} else {
				mapping[item] = make(map[string]interface{})
				tmp = mapping[item]
			}
		} else {
			return &KeyError{key: path.Join(append(keys[:i], item)...)}
		}
	}
	return nil
}

// Get is used to retrieve the value by indicating a key.
// After calling this method you should indicate the type of the return value.
// If one of the parent node is not a map, then return a ValueTypeError.
// If the key is not found, return a KeyError.
func (config *Config) Get(key string, defaultValue interface{}) (interface{}, error) {
	var (
		tmp interface{}
	)
	keys := strings.Split(key, "/")
	tmp = config.mapping
	for i, item := range keys {
		if len(item) == 0 {
			continue
		}
		if mapping, ok := tmp.(map[string]interface{}); ok {
			if value, exists := mapping[item]; exists {
				tmp = value
			} else if defaultValue != nil {
				return nil, &KeyError{key: path.Join(append(keys[:i], item)...)}
			} else {
				return nil, &KeyError{key: path.Join(append(keys[:i], item)...)}
			}
		} else {
			return nil, &ValueTypeError{key: path.Join(append(keys[:i], item)...), value: tmp, message: "Value is not a map"}
		}
	}
	return tmp, nil
}

// ConfigFromJSON creates a Config instance from a JSON io.reader.
func ConfigFromJSON(reader io.Reader) (*Config, error) {
	jsonBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &obj); err != nil {
		return nil, err
	}
	return &Config{obj}, nil
}
