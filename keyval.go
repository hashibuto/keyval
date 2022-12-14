package keyval

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type KeyVal struct {
	root map[string]any
}

// New returns an empty KeyVal instance
func New() *KeyVal {
	return &KeyVal{
		root: map[string]any{},
	}
}

// NewFromJson returns a new KeyVal instance from a JSON source
func NewFromJson(data []byte) (*KeyVal, error) {
	if data == nil {
		data = []byte("{}")
	}
	root := map[string]any{}
	err := json.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}

	return &KeyVal{
		root: root,
	}, nil
}

// NewFromJson returns a new KeyVal instance from a YAML source
func NewFromYaml(data []byte) (*KeyVal, error) {
	if data == nil {
		data = []byte("{}")
	}
	root := map[string]any{}
	err := yaml.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}

	return &KeyVal{
		root: root,
	}, nil
}

// NewFromMap returns a new KeyVal instance from a map[string]any
func NewFromMap(data map[string]any) *KeyVal {
	if data == nil {
		data = map[string]any{}
	}
	return &KeyVal{
		root: data,
	}
}

// GetKeyVal returns a new KeyVal object at the nested key position.
func (kv *KeyVal) GetKeyVal(keys ...string) (*KeyVal, error) {
	v, err := kv.Value(keys...)
	if err != nil {
		return nil, err
	}

	switch t := v.(type) {
	case map[string]any:
		return NewFromMap(t), nil
	default:
		return nil, fmt.Errorf("Data at key was not a generic map")
	}
}

// SplitKey splits a multi-part key string into its separate components.  The default delimiter is "."
func SplitKey(key string, delim ...string) []string {
	delimStr := "."
	if len(delim) > 0 {
		delimStr = delim[0]
	}

	return strings.Split(key, delimStr)
}

// SetValue sets a nested value within the object.  If a parent key cannot be located, an error is returned.
func (kv *KeyVal) SetValue(value any, keys ...string) error {
	var v any
	switch t := value.(type) {
	case int:
		v = float64(t)
	case int64:
		v = float64(t)
	case int32:
		v = float64(t)
	default:
		v = value
	}

	switch len(keys) {
	case 0:
	case 1:
		kv.root[keys[0]] = v
	default:
		target, err := walk(kv.root, false, keys[:len(keys)-1]...)
		if err != nil {
			return err
		}
		target[keys[len(keys)-1]] = v
	}
	return nil
}

// CreateValue sets a nested value within the object.  If a parent key cannot be located, it is created.
// Key collisions are ignored.
func (kv *KeyVal) CreateValue(value any, keys ...string) error {
	var v any
	switch t := value.(type) {
	case int:
		v = float64(t)
	case int64:
		v = float64(t)
	case int32:
		v = float64(t)
	default:
		v = value
	}

	switch len(keys) {
	case 0:
	case 1:
		kv.root[keys[0]] = v
	default:
		target, err := walk(kv.root, true, keys[:len(keys)-1]...)
		if err != nil {
			return err
		}
		target[keys[len(keys)-1]] = v
	}
	return nil
}

// Value returns a value or an error if the value cannot be located
func (kv *KeyVal) Value(keys ...string) (any, error) {
	var obj any = kv.root
	var ok bool
	for _, key := range keys {
		switch t := obj.(type) {
		case map[string]any:
			obj, ok = t[key]
			if !ok {
				return nil, fmt.Errorf("Could not resolve value using key")
			}
		default:
			return nil, fmt.Errorf("Encountered a non-mapping data type while traversing data")
		}
	}

	return obj, nil
}

// String returns a string or an error if the data can't be found, or properly cast
func (kv *KeyVal) String(keys ...string) (string, error) {
	v, err := kv.Value(keys...)
	if err != nil {
		return "", err
	}

	switch t := v.(type) {
	case string:
		return t, nil
	default:
		return "", fmt.Errorf("Value was not a string")
	}
}

// Number returns a float or an error if the data can't be found, or properly cast
func (kv *KeyVal) Number(keys ...string) (float64, error) {
	v, err := kv.Value(keys...)
	if err != nil {
		return 0.0, err
	}

	switch t := v.(type) {
	case float64:
		return t, nil
	default:
		return 0.0, fmt.Errorf("Value was not a number")
	}
}

// Boolean returns a boolean or an error if the data can't be found, or properly cast
func (kv *KeyVal) Boolean(keys ...string) (bool, error) {
	v, err := kv.Value(keys...)
	if err != nil {
		return false, err
	}

	switch t := v.(type) {
	case bool:
		return t, nil
	default:
		return false, fmt.Errorf("Value was not a boolean")
	}
}

// Array returns an array or an error if the data can't be found, or properly cast
func (kv *KeyVal) Array(keys ...string) ([]any, error) {
	v, err := kv.Value(keys...)
	if err != nil {
		return nil, err
	}

	switch t := v.(type) {
	case []any:
		return t, nil
	default:
		return nil, fmt.Errorf("Value was not an array")
	}
}

// Mapping returns an array or an error if the data can't be found, or properly cast
func (kv *KeyVal) Mapping(keys ...string) (map[string]any, error) {
	v, err := kv.Value(keys...)
	if err != nil {
		return nil, err
	}

	switch t := v.(type) {
	case map[string]any:
		return t, nil
	default:
		return nil, fmt.Errorf("Value was not a mapping")
	}
}

// Copy returns a deep copy of KeyVal
func (kv *KeyVal) Copy() *KeyVal {
	return &KeyVal{
		root: deepCopy(kv.root).(map[string]any),
	}
}

// Stack creates a new KeyVal object with the current instance being the base, and layer being stacked atop
func (kv *KeyVal) Stack(layer *KeyVal) *KeyVal {
	base := deepCopy(kv.root).(map[string]any)
	topLayer := deepCopy(layer.root).(map[string]any)

	stack(base, topLayer)
	return &KeyVal{
		root: base,
	}
}

// ToJson marshals the entire data structure to a JSON byte array
func (kv *KeyVal) ToJson() ([]byte, error) {
	return json.Marshal(kv.root)
}

// ToYaml marshals the entire data structure to a YAML byte array
func (kv *KeyVal) ToYaml() ([]byte, error) {
	return yaml.Marshal(kv.root)
}

// deepCopy returns a deep copy of obj
func deepCopy(obj any) any {
	switch t := obj.(type) {
	case []any:
		// Make an array copy
		target := make([]any, len(t))
		for idx, val := range t {
			switch val.(type) {
			case []any:
				val = deepCopy(val)
			case map[string]any:
				val = deepCopy(val)
			}
			target[idx] = val
		}
		return target
	case map[string]any:
		// Make a mapping copy
		target := map[string]any{}
		for key, val := range t {
			switch val.(type) {
			case []any:
				val = deepCopy(val)
			case map[string]any:
				val = deepCopy(val)
			}
			target[key] = val
		}
		return target
	default:
		return obj
	}
}

// isMapping returns true if value is a mapping type
func isMapping(value any) bool {
	switch value.(type) {
	case map[string]any:
		return true
	default:
		return false
	}
}

// stack stacks layerB atop layerA in place
func stack(layerA map[string]any, layerB map[string]any) {
	for key, newVal := range layerB {
		origVal, ok := layerA[key]
		if ok && isMapping(newVal) && isMapping(origVal) {
			stack(origVal.(map[string]any), newVal.(map[string]any))
		} else {
			layerA[key] = newVal
		}
	}
}

// walk walks a path through the object, arriving at the final target or returning an error
func walk(obj map[string]any, fill bool, keys ...string) (map[string]any, error) {
	var pos any = obj
	for _, key := range keys {
		switch t := pos.(type) {
		case map[string]any:
			target, ok := t[key]
			if ok {
				pos = target
				switch pos.(type) {
				case map[string]any:
				default:
					if fill {
						t[key] = map[string]any{}
						pos = t[key]
					} else {
						return nil, fmt.Errorf("Key was not reachable")
					}
				}
			} else {
				if fill {
					t[key] = map[string]any{}
					pos = t[key]
				} else {
					return nil, fmt.Errorf("Key missing during object traversal")
				}
			}
		default:
			return nil, fmt.Errorf("Object at key was incorrect type")
		}
	}

	return pos.(map[string]any), nil
}
