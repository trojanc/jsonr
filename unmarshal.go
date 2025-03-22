package jsonr

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var unwrappedType = reflect.TypeOf(Unwrapped{})

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
	isPointer := strings.HasPrefix(instanceType, "*")
	var result any
	t := getType(instanceType, opts)

	if t.Kind() == reflect.Slice {
		slice := reflect.MakeSlice(t, 0, 0)
		slicePtr := reflect.New(t)
		slicePtr.Elem().Set(slice)
		result = slicePtr.Interface()

		err := json.Unmarshal(wrapper.Value, result)
		if err != nil {
			return nil, err
		}
	} else if t.Kind() == reflect.Map {
		mapPtr := reflect.New(t)              // Pointer to map[string]MyStruct
		mapPtr.Elem().Set(reflect.MakeMap(t)) // Initialize the map
		result = mapPtr.Interface()
		err := json.Unmarshal(wrapper.Value, result)
		if err != nil {
			return nil, err
		}

		if t.Elem() == unwrappedType {
			// Create a new target map with value of any
			targetMapType := reflect.MapOf(t.Key(), reflect.TypeOf((*any)(nil)).Elem())
			targetMap := reflect.MakeMap(targetMapType)
			targetMapPtr := reflect.New(targetMapType) // Pointer to map[?]any
			targetMapPtr.Elem().Set(targetMap)         // Initialize the map

			iter := mapPtr.Elem().MapRange()
			for iter.Next() {
				keyValue := iter.Key()
				originalValue := iter.Value()

				uw := originalValue.Interface().(Unwrapped)
				unwrappedValue, err := Unwrap(uw, opts)
				if err != nil {
					return nil, err
				}
				mapValue := reflect.ValueOf(unwrappedValue)

				targetMap.SetMapIndex(keyValue, mapValue)
			}
			result = targetMapPtr.Interface()
		}

	} else {
		result = reflect.New(t).Interface()
		err := json.Unmarshal(wrapper.Value, result)
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("result %v\n", result)

	if !isPointer {
		if val, ok := removePointer(result); ok {
			result = val
		} else {
			return nil, errors.New("could not remove pointer")
		}
	}

	return result, nil
}

func getPointerValue(val reflect.Value) reflect.Value {
	// If already a pointer, return as-is
	if val.Kind() == reflect.Ptr {
		return val
	}

	// If the value is addressable, use Addr()
	if val.CanAddr() {
		return val.Addr()
	}

	// If not addressable, create a new instance and set the value
	ptr := reflect.New(val.Type()) // Create a pointer to the type
	ptr.Elem().Set(val)            // Copy value into the new pointer
	return ptr
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

func getType(instanceType string, opts *unmarshalOptions) reflect.Type {
	//isPointer := strings.HasPrefix(instanceType, "*")
	instanceType = strings.TrimPrefix(instanceType, "*")
	var typ reflect.Type

	if strings.HasPrefix(instanceType, "[]") {
		instanceType = instanceType[2:]
		sliceValPointer := strings.HasPrefix(instanceType, "*")
		sliceValType := getType(instanceType, opts)

		if sliceValPointer {
			sliceValType = reflect.PointerTo(sliceValType) // The slice should be pointers to the type
		}

		typ = reflect.SliceOf(sliceValType)
	} else if strings.HasPrefix(instanceType, "map[") {
		e := strings.Index(instanceType, "]")
		mapKeyType := instanceType[4:e]
		instanceType = instanceType[e+1:]

		mapKeyPointer := strings.HasPrefix(mapKeyType, "*")
		mapValPointer := strings.HasPrefix(instanceType, "*")

		// create a new map using reflection
		kt := getType(mapKeyType, opts)
		vt := getType(instanceType, opts)

		if vt.Kind() == reflect.Interface {
			vt = unwrappedType // If its any, put it back into a wrapper
		} else if mapValPointer {
			vt = reflect.PointerTo(vt)
		}

		if mapKeyPointer {
			kt = reflect.PointerTo(kt)
		}
		typ = reflect.MapOf(kt, vt)
	} else {
		if t, exists := opts.typeRegistry[instanceType]; exists {
			typ = t
		}
	}

	//if isPointer {
	//	typ = reflect.PointerTo(typ)
	//}

	return typ
}
