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

// newJSONRStruct creates a new jsonrStruct from the given value.
// The type is derived from the type of the given value. If the type is a pointer, it will be prefixed with "*".
// The value is the value itself.
func newJSONRStruct(v any) (*jsonrStruct, error) {
	t := reflect.TypeOf(v)
	typeName := t.Name()
	prefix := ""
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		prefix = "*"
		typeName = t.Name()
	}

	if t.Kind() == reflect.Slice {
		prefix = prefix + "[]"
		t = t.Elem()
		typeName = t.Name()
	}

	if t.Kind() == reflect.Map {
		// TODo what if its a slice of maps
		return processMap(t, v)
	}

	valStr, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return &jsonrStruct{
		Type:  fmt.Sprintf("%s%s.%s", prefix, t.PkgPath(), typeName),
		Value: string(valStr),
	}, nil
}

func processMap(mapType reflect.Type, v any) (*jsonrStruct, error) {
	// use reflection to get the type of the map's key
	key := mapType.Key()
	valueType := mapType.Elem()
	valTypeName := ""
	keyTypeName := key.Kind().String()
	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
		valTypeName = "*"
	}
	valTypeName = valTypeName + valueType.PkgPath() + "." + valueType.Name()

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
