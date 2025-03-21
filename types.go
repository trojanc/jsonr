package jsonr

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
