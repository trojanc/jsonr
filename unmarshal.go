package jsonr

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Unmarshal(data string, instanceType string, opts ...UnmarshalOption) (any, error) {
	// Build an unmarshalOptions object from the provided options
	options, err := applyUnmarshalOptions(opts...)
	if err != nil {
		return nil, err
	}

	return deserializeObject(data, instanceType, options)
}

// deserializeObject converts an Object back into the original Go structure
func deserializeObject(data string, instanceType string, opts *unmarshalOptions) (any, error) {
	isPointer := false
	isSlice := false
	isMap := false
	mapKeyType := ""
	var value any
	var ok bool
	if strings.HasPrefix(instanceType, "*") {
		// Remove pointer from type name
		instanceType = instanceType[1:]
		isPointer = true
	}

	if strings.HasPrefix(instanceType, "[]") {
		isSlice = true
		instanceType = instanceType[2:]
	} else if strings.HasPrefix(instanceType, "map[") {
		isMap = true
		e := strings.Index(instanceType, "]")
		mapKeyType = instanceType[4:e]
		instanceType = instanceType[e+1:]
	}

	// creates a new instance with the type from the map
	instance := newInstance(instanceType, opts)
	if instance == nil {
		return nil, fmt.Errorf("unmarshalling unknown type %s", instanceType)
	}

	if isSlice {
		// create a new slice using reflection
		result := reflect.MakeSlice(reflect.TypeOf(instance).Elem(), 0, 0)
		//for i, val := range data {
		//	convertedVal, err := deserializeObject(val, targetType.Elem())
		//	if err != nil {
		//		return nil, err
		//	}
		//	result.Index(i).Set(reflect.ValueOf(convertedVal))
		//}
		return result.Interface(), nil
	} else if isMap {
		mapKeyPointer := strings.HasPrefix(mapKeyType, "*")
		mapValPointer := strings.HasPrefix(instanceType, "*")

		// create a new map using reflection
		mapKeyInstance := newInstance(mapKeyType, opts)
		mapValInstance := newInstance(instanceType, opts)

		if !mapKeyPointer {
			mapKeyInstance, ok = removePointer(mapKeyInstance)
			if !ok {
				return nil, errors.New("could not remove pointer")
			}
		}
		if !mapValPointer {
			mapValInstance, ok = removePointer(mapValInstance)
			if !ok {
				return nil, errors.New("could not remove pointer")
			}
		}

		mapType := reflect.MapOf(reflect.TypeOf(mapKeyInstance), reflect.TypeOf(mapValInstance))
		instance = reflect.New(mapType).Interface()
		// Create the map from the map
		m := reflect.ValueOf(value)
		result := make(map[string]any)
		for _, key := range m.MapKeys() {
			val := m.MapIndex(key)
			cv, err := deserializeObject("val", val.String(), opts) // TODO
			if err != nil {
				return nil, err
			}
			result[key.String()] = cv
		}
	} else {
		err := json.Unmarshal([]byte(data), instance)
		if err != nil {
			return nil, err
		}
	}

	if !isPointer {
		if val, ok := removePointer(instance); ok {
			instance = val
		} else {
			return nil, errors.New("could not remove pointer")
		}
	}
	return instance, nil
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

// newInstance Create a new instance given a type name
func newInstance(typeName string, opts *unmarshalOptions) any {
	if strings.HasPrefix(typeName, "*") {
		// Remove pointer from type name
		typeName = typeName[1:]
	}

	t, exists := opts.typeMapping[typeName]
	if !exists {
		fmt.Printf("could not find %s\n", typeName)
		return nil // Type not found
	}
	return reflect.New(t).Interface() // Create a new instance (as pointer)
}
