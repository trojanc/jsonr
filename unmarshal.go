package jsonr

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// unwrappedType reference to the type of Unwrapped.
var unwrappedType = reflect.TypeOf(Unwrapped{})

// Unwrapped a structure of an unwrapped type partially read from JSON
type Unwrapped struct {
	Type  string          `json:"_t"`
	Value json.RawMessage `json:"v"`
}

// Unmarshal decodes JSON data into a Go value with type information. It expects JSON data that was previously
// encoded using Marshal() which includes type information in a wrapper structure.
//
// The function takes a byte slice containing the JSON data and optional UnmarshalOptions. These options can be used to
// register types that should be available for unmarshalling.
//
// Example usage:
//
//		type Person struct {
//		    Name string
//		    Age  int
//		}
//
//		// Marshal data with type information
//		data, _ := jsonr.Marshal(Person{Name: "John", Age: 30})
//	 	// data will be {"_t":"github.com/project/example.Person","v":{"Name":"John","Age":30}}
//
//		// Unmarshal back into interface{}
//		result, _ := jsonr.Unmarshal(data, jsonr.RegisterType(example.Person{}))
//		// result will be Person{Name: "John", Age: 30}
//
// The function also supports unmarshalling complex types like maps and slices:
//
//	// Map example
//	data, _ := jsonr.Marshal(map[string]Person{
//	    "John": {Name: "John", Age: 30},
//	    "Jane": {Name: "Jane", Age: 25},
//	})
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

// Unwrap decodes a wrapped JSON structure back into its original Go value. It takes a Unwrapped struct containing
// type information and raw JSON data, along with unmarshal options for type registration.
//
// The function handles complex types like maps and slices, including cases where values are interface{} types
// that need recursive unwrapping. For maps and slices containing interface{} values, it creates new target
// collections with properly unwrapped elements.
//
// Example usage:
//
//	wrapper := Unwrapped{
//	    Type: "github.com/project/example.Person",
//	    Value: json.RawMessage(`{"name":"John","age":30}`),
//	}
//	result, _ := Unwrap(wrapper, opts)
//	// result will be Person{Name: "John", Age: 30}
//
// The function supports:
// - Primitive Go types
// - Structs and pointers to structs
// - Maps with primitive keys and any value type
// - Slices of any type
// - Nested combinations of the above
func Unwrap(wrapper Unwrapped, opts *unmarshalOptions) (any, error) {

	instanceType := wrapper.Type
	isPointer := strings.HasPrefix(instanceType, "*")
	var result any
	t := getType(instanceType, opts)

	if t.Kind() == reflect.Slice {
		slicePtr := reflect.New(t)
		result = slicePtr.Interface()
		err := json.Unmarshal(wrapper.Value, result)
		if err != nil {
			return nil, err
		}
		fmt.Printf("slice : %v\n", result)
		if t.Elem() == unwrappedType {
			slice := slicePtr.Elem()
			// Create a new target map with value of any
			targetSliceType := reflect.SliceOf(reflect.TypeOf((*any)(nil)).Elem())
			targetSlice := reflect.MakeSlice(targetSliceType, slice.Len(), slice.Len())
			targetSlicePtr := reflect.New(targetSliceType) // Pointer to map[?]any
			targetSlicePtr.Elem().Set(targetSlice)         // Initialize the map

			// Iterate over the slice
			for i := 0; i < slice.Len(); i++ {
				elem := slice.Index(i)
				uw := elem.Interface().(Unwrapped)
				unwrappedValue, err := Unwrap(uw, opts)
				if err != nil {
					return nil, err
				}
				mapValue := reflect.ValueOf(unwrappedValue)
				targetSlice.Index(i).Set(mapValue)
			}
			result = targetSlicePtr.Interface()
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

		if sliceValType.Kind() == reflect.Interface {
			sliceValType = unwrappedType // If its any, put it back into a wrapper
		} else if sliceValPointer {
			sliceValType = reflect.PointerTo(sliceValType)
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

	return typ
}
