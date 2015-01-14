package Golf

import (
	"fmt"
	"path"
	"reflect"
	"strings"
)

type KeyError struct {
	key string
}

func (err *KeyError) Error() string {
	return fmt.Sprintf("Key not found, key: %s", err.key)
}

type ValueTypeError struct {
	key     string
	value   interface{}
	message string
}

func (err *ValueTypeError) Error() string {
	return fmt.Sprintf("%s, key: %s, value: %v (%s)", err.message, err.key, err.value, reflect.TypeOf(err.value).Name())
}

type Config struct {
	app     *Application
	mapping map[string]interface{}
}

func NewConfig(app *Application) *Config {
	mapping := make(map[string]interface{})
	return &Config{app, mapping}
}

func (config *Config) GetString(key string, defaultValue interface{}) (string, error) {
	value, err := config.Get(key, defaultValue)
	if err != nil {
		return "", err
	}
	if result, ok := value.(string); ok {
		return result, nil
	} else {
		return "", &ValueTypeError{key: key, value: value, message: "Value is not a string."}
	}
}

func (config *Config) GetInt(key string, defaultValue interface{}) (int, error) {
	value, err := config.Get(key, defaultValue)
	if err != nil {
		return 0, err
	}
	if result, ok := value.(int); ok {
		return result, nil
	} else {
		return 0, &ValueTypeError{key: key, value: value, message: "Value is not an integer."}
	}
}

func (config *Config) GetBool(key string, defaultValue interface{}) (bool, error) {
	value, err := config.Get(key, defaultValue)
	if err != nil {
		return false, err
	}
	if result, ok := value.(bool); ok {
		return result, nil
	} else {
		return false, &ValueTypeError{key: key, value: value, message: "Value is not an bool."}
	}
}

func (config *Config) GetFloat(key string, defaultValue interface{}) (float64, error) {
	value, err := config.Get(key, defaultValue)
	if err != nil {
		return 0, err
	}
	if result, ok := value.(float64); ok {
		return result, nil
	} else {
		return 0, &ValueTypeError{key: key, value: value, message: "Value is not a float."}
	}
}

func (config *Config) Set(key string, value interface{}) error {
	var tmp interface{}
	keys := strings.Split(key, "/")
	tmp = config.mapping
	for i, item := range keys {
		if len(item) == 0 {
			continue
		}
		if mapping, ok := tmp.(map[string]interface{}); ok {
			if i == len(keys) - 1 {
				mapping[item] = value
				return nil
			} else {
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
			}
		} else {
		    return &KeyError{key: path.Join(append(keys[:i], item)...)}
    }
	}
  return nil
}

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
				return defaultValue, nil
			} else {
				return nil, &KeyError{key: path.Join(append(keys[:i], item)...)}
			}
		} else {
			return nil, &ValueTypeError{key: path.Join(append(keys[:i], item)...), value: tmp, message: "Value is not a map"}
		}
	}
	return tmp, nil
}
