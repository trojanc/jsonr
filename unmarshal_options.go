package jsonr

import (
	"errors"
	"fmt"
	"reflect"
)

// variableTypeMapping defines a type that can be used to map type keys to actual relection types
type variableTypeMapping map[string]reflect.Type

// unmarshalOptions Options that will be used while unmarshalling the engine
type unmarshalOptions struct {
	exportTypes bool
	typeMapping variableTypeMapping
}

// UnmarshalOption is a function that modifies the unmarshalOptions
type UnmarshalOption func(*unmarshalOptions) error

// registerTypeOption is a function that registers a type in the type mapping
type registerTypeOption func(map[string]reflect.Type) error

// WithUnmarshalComplexTypes enables type unmarshalling.
// Additional types can be registered using the RegisterType function.
func WithUnmarshalComplexTypes(at ...registerTypeOption) UnmarshalOption {
	// All default complex types that the engine may produce must be registered by default here
	var mappingOptions []registerTypeOption
	mappingOptions = append(mappingOptions, at...)

	return func(opts *unmarshalOptions) error {
		opts.exportTypes = true

		opts.typeMapping["string"] = reflect.TypeOf("")
		opts.typeMapping["int"] = reflect.TypeOf(0)
		opts.typeMapping["float64"] = reflect.TypeOf(0.0)
		opts.typeMapping["bool"] = reflect.TypeOf(true)
		opts.typeMapping["any"] = reflect.TypeOf(new(any)).Elem()

		for _, mo := range mappingOptions {
			err := mo(opts.typeMapping)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// RegisterType registers a type that can be unmarshalled into an instance of the given type.
func RegisterType(instance any) func(map[string]reflect.Type) error {
	return func(m map[string]reflect.Type) error {
		t := reflect.TypeOf(instance)

		// Do not allow pointers or any other basic types to be passed in as an instance type
		// Marshalling and Unmarshalling will take care of pointers
		if t.Kind() != reflect.Struct {
			return errors.New("only instance of structs should be used")
		}

		typeKey := fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
		m[typeKey] = t
		return nil
	}
}

// applyUnmarshalOptions Applies the given options and returns the applied unmarshalOptions
func applyUnmarshalOptions(options ...UnmarshalOption) (*unmarshalOptions, error) {
	opts := &unmarshalOptions{
		typeMapping: make(map[string]reflect.Type),
	}
	for _, o := range options {
		err := o(opts)
		if err != nil {
			return nil, fmt.Errorf("could not apply option: %w", err)
		}
	}

	return opts, nil
}
