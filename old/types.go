package old

import "reflect"

var (
	// ComplexKinds are the kinds that are considered complex types
	ComplexKinds = []reflect.Kind{
		reflect.Struct,
		reflect.Ptr,
		reflect.Slice,
		reflect.Map,
	}
)

// Object this struct is used when marshalling with WithMarshalComplexTypes to export complex types
type Object struct {
	Type  string `json:"_t"`
	Value any    `json:"v"`
}
