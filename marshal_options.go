package jsonr

import "fmt"

// marshalOptions Options that will be used while marshalling the engine
type marshalOptions struct {
	exportTypes bool
}

// MarshalOption is a function that modifies the marshalOptions
type MarshalOption func(*marshalOptions) error

// WithMarshalComplexTypes if added as an option the marshaller will export variables with their specific types.
// When this is used, calls to Unmarshal will need to use RegisterType to configure the types that can be
// unmarshalled into actual instances
// .
// This is useful when you have complex types in your variables that you want to preserve.
//
//	Example:
//	```go
//	// Marshal with type information
//	data := engine.Marshal(WithMarshalComplexTypes())
//
//	// Unmarshal with type mapping
//	engine, _ = Unmarshal(data, RegisterType(MyStruct{}))
//	```
func WithMarshalComplexTypes() MarshalOption {
	return func(opts *marshalOptions) error {
		opts.exportTypes = true
		return nil
	}
}

// applyMarshalOptions Applies the given options and returns the applied marshalOptions
func applyMarshalOptions(options ...MarshalOption) (*marshalOptions, error) {
	opts := &marshalOptions{}
	for _, o := range options {
		err := o(opts)
		if err != nil {
			return nil, fmt.Errorf("could not apply option: %w", err)
		}
	}

	return opts, nil
}
