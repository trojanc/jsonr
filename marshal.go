package jsonr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
)

// jsonrStruct this struct is used when marshalling with WithMarshalComplexTypes to export complex types
type jsonrStruct struct {
	Type  string `json:"_t"`
	Value string `json:"v"`
}

func Marshal(v any, opts ...marshalOptions) ([]byte, error) {

	// Nothing to marshal
	if v == nil {
		return []byte{}, nil
	}

	refl := reflect.ValueOf(v)
	kind := refl.Kind()

	// Is this a complex type?
	if slices.Contains(ComplexKinds, kind) {
		j, err := newJSONRStruct(v)
		if err != nil {
			return []byte{}, fmt.Errorf("could not create jsonrStruct: %w", err)
		}
		return json.Marshal(j)
	}

	return json.Marshal(v)
}

func getFullTypeName(t reflect.Type) (string, reflect.Type) {
	prefix := ""
	typeName := ""
	if t.Kind() == reflect.Interface {
		return "any", t
	}

	if !slices.Contains(ComplexKinds, t.Kind()) {
		return t.Name(), t
	} else {
		typeName = fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		prefix = "*"
		typeName = fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
	}

	if t.Kind() == reflect.Slice {
		prefix = prefix + "[]"
		typeName, t = getFullTypeName(t.Elem())
	}

	if t.Kind() == reflect.Interface {
		return "any", t
	}

	if t.Kind() == reflect.Map {
		keyTypeName, _ := getMapKeyName(t.Key()) // TODO erro
		valTypeName, _ := getFullTypeName(t.Elem())
		return fmt.Sprintf("%smap[%s]%s", prefix, keyTypeName, valTypeName), t
	}

	return fmt.Sprintf("%s%s", prefix, typeName), t
}

// newJSONRStruct creates a new jsonrStruct from the given value.
// The type is derived from the type of the given value. If the type is a pointer, it will be prefixed with "*".
// The value is the value itself.
func newJSONRStruct(v any) (*jsonrStruct, error) {
	t := reflect.TypeOf(v)

	typeName, t := getFullTypeName(t)

	if t.Kind() == reflect.Map {
		// TODo what if its a slice of maps
		return processMap(t, v)
	}

	valStr, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return &jsonrStruct{
		Type:  typeName,
		Value: string(valStr),
	}, nil
}

func getMapKeyName(keyType reflect.Type) (string, error) {
	if keyType.Kind() == reflect.Ptr {
		return "", fmt.Errorf("map keys cannot be pointers")
	} else if keyType.Kind() == reflect.Struct {
		return "", fmt.Errorf("map keys cannot be structs")
	} else {
		return keyType.Name(), nil
	}
}

func processMap(mapType reflect.Type, v any) (*jsonrStruct, error) {
	// use reflection to get the type of the map's key
	keyType := mapType.Key()
	keyTypeName, err := getMapKeyName(keyType)
	if err != nil {
		return nil, err
	}

	valTypeName, _ := getFullTypeName(mapType.Elem())

	// loop over v as a map and create a slice of complex variables
	m := reflect.ValueOf(v)
	result := make(map[string]any)
	for _, key := range m.MapKeys() {
		val := m.MapIndex(key)
		cv, err := newJSONRStruct(val.Interface())
		if err != nil {
			return nil, err
		}
		result[key.String()] = cv
	}
	valStr, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return &jsonrStruct{
		Type:  fmt.Sprintf("map[%s]%s", keyTypeName, valTypeName),
		Value: string(valStr),
	}, nil
}
