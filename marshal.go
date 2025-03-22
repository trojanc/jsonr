package jsonr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
)

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
func Wrap(input any) (*Object, error) {

	if input == nil {
		return nil, nil // TODO error?
	}

	// Get reflection type and value
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)

	// Is this a complex type?
	if !slices.Contains(ComplexKinds, t.Kind()) {
		return nil, nil // TODO error?
	}

	typeName := getTypeName(t)
	// Dereference pointers
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	switch t.Kind() {
	case reflect.Map:
		keyType := t.Key()
		if keyType.Kind() == reflect.Ptr {
			return nil, fmt.Errorf("map keys cannot be pointers")
		} else if keyType.Kind() == reflect.Struct {
			return nil, fmt.Errorf("map keys cannot be structs")
		}
		//valueType := t.Elem()
		//if valueType.Kind() == reflect.Interface {
		//	return nil, fmt.Errorf("map value can not by any")
		//}
	case reflect.Slice, reflect.Array:
	case reflect.Struct:

	default:
		return nil, fmt.Errorf("unsupported type: %s", t.Kind())
	}

	return &Object{
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
	case reflect.Interface:
		return "any"
	default:
		return t.String()
	}
}
