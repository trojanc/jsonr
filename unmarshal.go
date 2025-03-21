package jsonr

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

type Unwrapped struct {
	Type  string          `json:"_t"`
	Value json.RawMessage `json:"v"`
}

func Unmarshal(data []byte, options ...UnmarshalOption) (any, error) {

	// Build an unmarshalOptions object from the provided options
	opts, err := applyUnmarshalOptions(options...)
	if err != nil {
		return nil, err
	}

	// Step 1: Extract the type field
	var wrapper Unwrapped
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}

	return Unwrap(wrapper, opts)
}

func Unwrap(wrapper Unwrapped, opts *unmarshalOptions) (any, error) {

	instanceType := wrapper.Type
	isPointer := false
	mapKeyType := ""
	//var value any
	var result any
	if strings.HasPrefix(instanceType, "*") {
		// Remove pointer from type name
		instanceType = instanceType[1:]
		isPointer = true
	}

	if strings.HasPrefix(instanceType, "[]") {
		instanceType = instanceType[2:]
		sliceValPointer := strings.HasPrefix(instanceType, "*")
		sliceValType := getType(instanceType, opts)

		if sliceValPointer {
			sliceValType = reflect.PointerTo(sliceValType)
		}

		sliceType := reflect.SliceOf(sliceValType)
		slice := reflect.MakeSlice(sliceType, 0, 0)

		slicePtr := reflect.New(sliceType)
		slicePtr.Elem().Set(slice)
		result = slicePtr.Interface()

	} else if strings.HasPrefix(instanceType, "map[") {
		e := strings.Index(instanceType, "]")
		mapKeyType = instanceType[4:e]
		instanceType = instanceType[e+1:]

		mapKeyPointer := strings.HasPrefix(mapKeyType, "*")
		mapValPointer := strings.HasPrefix(instanceType, "*")

		// create a new map using reflection
		kt := getType(mapKeyType, opts)
		vt := getType(instanceType, opts)

		if mapKeyPointer {
			kt = reflect.PointerTo(kt)
		}
		if mapValPointer {
			vt = reflect.PointerTo(vt)
		}

		// Create a map type (map[string]MyStruct)
		mapType := reflect.MapOf(kt, vt)

		// Create a new map instance
		mapPtr := reflect.New(mapType)              // Pointer to map[string]MyStruct
		mapPtr.Elem().Set(reflect.MakeMap(mapType)) // Initialize the map

		result = mapPtr.Interface()
	} else {
		result = newInstance(instanceType, opts)
	}

	err := json.Unmarshal(wrapper.Value, result)
	if err != nil {
		return nil, err
	}

	if !isPointer {
		if val, ok := removePointer(result); ok {
			result = val
		} else {
			return nil, errors.New("could not remove pointer")
		}
	}

	return result, nil
}

// removePointer Function to remove pointer from an `any` type variable
func removePointer(v any) (any, bool) {
	// Use type assertion to check if it's a pointer
	if ptr, ok := v.(interface{ Elem() any }); ok {
		return ptr.Elem(), true
	}

	// Use reflection as a fallback for generic cases
	return removePointerReflect(v)
}

// removePointerReflect Function to remove pointer from an `any` type variable
// using reflection
func removePointerReflect(v any) (any, bool) {
	// Use reflection to handle arbitrary pointer types
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		return rv.Elem().Interface(), true
	}
	return v, false
}

func getType(typeName string, opts *unmarshalOptions) reflect.Type {
	typeName = strings.TrimPrefix(typeName, "*")
	t, exists := opts.typeRegistry[typeName]
	if !exists {
		return nil // Type not found
	}
	return t
}

// newInstance Create a new instance given a type name
func newInstance(typeName string, opts *unmarshalOptions) any {
	t := getType(typeName, opts)
	return reflect.New(t).Interface() // Create a new instance (as pointer)
}
