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

// Marshal encodes a Go value into JSON with type information. It wraps the value in a structure that includes
// the Go type, allowing for proper type reconstruction during unmarshalling.
//
// The function handles:
// - Primitive Go types
// - Structs and pointers to structs
// - Maps with primitive keys and any value type
// - Slices of any type
// - Nested combinations of the above
//
// Example usage:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	// Marshal a struct
//	person := Person{Name: "John", Age: 30}
//	data, _ := jsonr.Marshal(person)
//	// data will be {"_t":"github.com/project/example.Person","v":{"Name":"John","Age":30}}
//
//	// Marshal complex types
//	people := map[string]Person{
//	    "john": {Name: "John", Age: 30},
//	    "jane": {Name: "Jane", Age: 25},
//	}
//	data, _ := jsonr.Marshal(people)
//
// data will be {"_t":"map[string]github.com/project/example.Person","v":{"john":{"Name":"John","Age":30},"jane":{"Name":"Jane","Age":25}}}
func Marshal(input any) ([]byte, error) {
	w, err := Wrap(input)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(w)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %s", err.Error())
	}
	return data, nil
}

// Wrap takes a Go value and wraps it in a structure that includes type information. This allows for proper
// type reconstruction during unmarshalling. The function handles various Go types including primitives,
// structs, maps, slices, and their nested combinations.
//
// For maps and slices containing interface{} values, it recursively wraps each element to preserve type
// information throughout the entire data structure.
//
// Example usage:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	// Wrap a struct
//	person := Person{Name: "John", Age: 30}
//	wrapped, _ := jsonr.Wrap(person)
//	// wrapped will contain type information and value
//
//	// Wrap complex types
//	people := map[string]Person{
//		"john": {Name: "John", Age: 30},
//		"jane": {Name: "Jane", Age: 25},
//	}
//
//	wrapped, _ := jsonr.Wrap(people)
//	wrapped will contain type information and value
func Wrap(input any) (*Wrapped, error) {
	if input == nil {
		return nil, nil
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

		// Check the map key type
		switch t.Key().Kind() {
		case reflect.Ptr, reflect.Struct, reflect.Map, reflect.Slice:
			return nil, fmt.Errorf("unsupported map key")
		default:
		}

		// Check the map value type
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
