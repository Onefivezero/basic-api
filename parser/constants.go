package parser

import (
	"reflect"
	"slices"
)

var intKinds = []reflect.Kind{
	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
}

var uintKinds = []reflect.Kind{
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
}

var floatKinds = []reflect.Kind{
	reflect.Float32, reflect.Float64,
}

var complexKinds = []reflect.Kind{
	reflect.Complex64, reflect.Complex128,
}

var boolKinds = []reflect.Kind{reflect.Bool}

var stringKinds = []reflect.Kind{reflect.String}

var structKinds = []reflect.Kind{reflect.Struct}

var sliceKinds = []reflect.Kind{reflect.Slice, reflect.Array}

func isStruct(type_ reflect.Type) bool {
	return slices.Contains(structKinds, type_.Kind())
}

func isSlice(type_ reflect.Type) bool {
	return slices.Contains(sliceKinds, type_.Kind())
}
