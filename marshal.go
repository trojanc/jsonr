package jsonr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func Marshal(v any, opts ...MarshalOption) ([]byte, error) {

	marshalOpts := &marshalOptions{}
	for _, opt := range opts {
		err := opt(marshalOpts)
		if err != nil {
			return []byte{}, err
		}
	}

	// Nothing to marshal
	if v == nil {
		return []byte{}, nil
	}

	refl := reflect.ValueOf(v)
	kind := refl.Kind()

	// Is this a complex type?
	if slices.Contains(ComplexKinds, kind) {
		j, err := serializeToObject(v)
		if err != nil {
			return []byte{}, fmt.Errorf("could not create Object: %w", err)
		}
		return json.Marshal(j)
	}

	return json.Marshal(v)
}

// serializeToObject serializes any value, using JSON struct tags and omitting default values.
func serializeToObject(input any) (any, error) {
	if input == nil {
		return nil, nil
	}

	typeName := ""

	// Get reflection type and value
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)

	// Dereference pointers
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
		typeName = "*"
	}

	// Handle primitive types (omit if default)
	if isPrimitiveType(t) {
		if isDefaultValue(v) {
			return nil, nil
		}
		return input, nil
	}

	// Get type name for complex structures
	typeName = typeName + getTypeName(t)

	// Handle complex structures
	var jsonData []byte
	var err error

	switch t.Kind() {
	case reflect.Map:
		convertedMap := make(map[string]interface{})
		iter := v.MapRange()
		keyType := t.Key()
		if keyType.Kind() == reflect.Ptr {
			return nil, fmt.Errorf("map keys cannot be pointers")
		} else if keyType.Kind() == reflect.Struct {
			return nil, fmt.Errorf("map keys cannot be structs")
		}
		for iter.Next() {
			key := fmt.Sprintf("%v", iter.Key())
			val, err := serializeToObject(iter.Value().Interface())
			if err != nil {
				return nil, err
			}
			if val != nil {
				convertedMap[key] = val
			}
		}
		if len(convertedMap) == 0 {
			jsonData = []byte("{}")
		} else {
			jsonData, err = json.Marshal(convertedMap)
		}

	case reflect.Slice, reflect.Array:
		var convertedSlice []interface{}
		for i := 0; i < v.Len(); i++ {
			val, err := serializeToObject(v.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			if val != nil {
				convertedSlice = append(convertedSlice, val)
			}
		}
		if len(convertedSlice) == 0 {
			jsonData = []byte("[]")
		} else {
			jsonData, err = json.Marshal(convertedSlice)
		}

	case reflect.Struct:
		convertedStruct := make(map[string]interface{})
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldVal := v.Field(i)

			// Respect JSON tags
			jsonKey, ignore := getJSONTag(field)
			if ignore {
				continue
			}

			// Convert field value
			val, err := serializeToObject(fieldVal.Interface())
			if err != nil {
				return nil, err
			}
			if val != nil {
				convertedStruct[jsonKey] = val
			}
		}
		if len(convertedStruct) == 0 {
			jsonData = []byte("{}")
		} else {
			jsonData, err = json.Marshal(convertedStruct)
		}

	default:
		return nil, fmt.Errorf("unsupported type: %s", t.Kind())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to serialize: %w", err)
	}

	// Wrap complex structures in Object
	return &Object{Type: typeName, Value: string(jsonData)}, nil
}

// isPrimitiveType checks if a type is a primitive (int, string, float, bool, etc.)
func isPrimitiveType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String:
		return true
	}
	return false
}

// isDefaultValue checks if a value is the default for its type
func isDefaultValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	}
	return false
}

// getTypeName returns a structured type name for deeply nested types
func getTypeName(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		return "[]" + getTypeName(t.Elem())
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", getTypeName(t.Key()), getTypeName(t.Elem()))
	case reflect.Ptr:
		return "*" + getTypeName(t.Elem())
	case reflect.Struct:
		return t.PkgPath() + "." + t.Name()
	case reflect.Interface:
		return "any"
	default:
		return t.String()
	}
}

// getJSONTag extracts JSON field name from struct tags while respecting json:"-" ignore rules
func getJSONTag(field reflect.StructField) (string, bool) {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return "", true // Ignore field
	}

	if tag == "" {
		return field.Name, false // Default to field name
	}

	// Extract the actual JSON key (before ",omitempty", ",string", etc.)
	jsonKey := strings.Split(tag, ",")[0]
	return jsonKey, false
}
