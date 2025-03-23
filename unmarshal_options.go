package jsonr

import (
	"errors"
	"fmt"
	"reflect"
)

// typeRegistry defines a type that can be used to map type keys to actual relection types
type typeRegistry map[string]reflect.Type

// unmarshalOptions Options that will be used while unmarshalling the engine
type unmarshalOptions struct {
	// typeRegistry registry of types that can be unmarshalled
	typeRegistry typeRegistry
}

// UnmarshalOption is a function that modifies the unmarshalOptions
type UnmarshalOption func(*unmarshalOptions) error

// RegisterType registers a type that can be unmarshalled into an instance of the given type.
func RegisterType(instance any) UnmarshalOption {
	return func(opts *unmarshalOptions) error {
		t := reflect.TypeOf(instance)

		// Do not allow pointers or any other basic types to be passed in as an instance type
		// Marshalling and Unmarshalling will take care of pointers
		if t.Kind() != reflect.Struct {
			return errors.New("only instance of structs should be used")
		}

		typeKey := fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
		opts.typeRegistry[typeKey] = t
		return nil
	}
}

// applyUnmarshalOptions Applies the given options and returns the applied unmarshalOptions
func applyUnmarshalOptions(options ...UnmarshalOption) (*unmarshalOptions, error) {
	opts := &unmarshalOptions{
		typeRegistry: make(map[string]reflect.Type),
	}

	// Add the primitive types to the type registry
	opts.typeRegistry["int"] = reflect.TypeOf(int(0))
	opts.typeRegistry["int8"] = reflect.TypeOf(int8(0))
	opts.typeRegistry["int16"] = reflect.TypeOf(int16(0))
	opts.typeRegistry["int32"] = reflect.TypeOf(int32(0))
	opts.typeRegistry["int64"] = reflect.TypeOf(int64(0))
	opts.typeRegistry["uint"] = reflect.TypeOf(uint(0))
	opts.typeRegistry["uint8"] = reflect.TypeOf(uint8(0))
	opts.typeRegistry["uint16"] = reflect.TypeOf(uint16(0))
	opts.typeRegistry["uint32"] = reflect.TypeOf(uint32(0))
	opts.typeRegistry["uint64"] = reflect.TypeOf(uint64(0))
	opts.typeRegistry["float32"] = reflect.TypeOf(float32(0))
	opts.typeRegistry["float64"] = reflect.TypeOf(float64(0))
	opts.typeRegistry["complex64"] = reflect.TypeOf(complex64(0))
	opts.typeRegistry["complex128"] = reflect.TypeOf(complex128(0))
	opts.typeRegistry["bool"] = reflect.TypeOf(false)
	opts.typeRegistry["string"] = reflect.TypeOf("")
	opts.typeRegistry["byte"] = reflect.TypeOf(byte(0))
	opts.typeRegistry["rune"] = reflect.TypeOf(rune(0))
	opts.typeRegistry["interface"] = reflect.TypeOf(new(any)).Elem()

	for _, o := range options {
		err := o(opts)
		if err != nil {
			return nil, fmt.Errorf("could not apply option: %w", err)
		}
	}

	return opts, nil
}
