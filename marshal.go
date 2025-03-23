package jsonr

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Wrapped this struct is used when marshalling with WithMarshalComplexTypes to export complex types
type Wrapped struct {
	Type  string `json:"_t"`
	Value any    `json:"v"`
}

func Marshal(input any) ([]byte, error) {
	w, err := Wrap(input)
	if err != nil {
		return nil, err
	}
	return json.Marshal(w)
}

// Wrap
// If maps are used, the keys must be a primitive type
// Values (of maps and slices too) must be json marshallable
// Structs with `any` fields will not work as expected
func Wrap(input any) (*Wrapped, error) {

	if input == nil {
		return nil, nil // TODO error?
	}

	// Get reflection type and value
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)

	typeName := getTypeName(t)
	// Dereference pointers
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() == reflect.Map {
		switch t.Key().Kind() {
		case reflect.Ptr, reflect.Struct, reflect.Map, reflect.Slice:
			return nil, fmt.Errorf("unsupported map key")
		default:
		}

		switch t.Elem().Kind() {
		case reflect.Interface:
			// rebuild the map by wrapping each value
			m := make(map[string]any)
			for _, k := range v.MapKeys() {
				w, err := Wrap(v.MapIndex(k).Interface())
				if err != nil {
					return nil, err
				}
				m[k.String()] = w
			}
			input = m
		default:
		}
	} else if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		if t.Elem().Kind() == reflect.Interface {
			// rebuild the slice
			s := make([]any, 0)
			for i := 0; i < v.Len(); i++ {
				w, err := Wrap(v.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				s = append(s, w)
			}
			input = s
		}
	}

	return &Wrapped{
		Type:  typeName,
		Value: input,
	}, nil
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
	default:
		return t.Kind().String()
	}
}
