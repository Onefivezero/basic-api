package basic_api

import (
	"reflect"
	"slices"
)

var numKinds = []reflect.Kind{
	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
}

var boolKinds = []reflect.Kind{reflect.Bool}

var stringKinds = []reflect.Kind{reflect.String}

var structKinds = []reflect.Kind{reflect.Struct}

var sliceKinds = []reflect.Kind{reflect.Slice, reflect.Array}

func isNormal(type_ reflect.Type) bool {
	return (slices.Contains(numKinds, type_.Kind()) ||
		slices.Contains(stringKinds, type_.Kind()) ||
		slices.Contains(boolKinds, type_.Kind()))
}

func isStruct(type_ reflect.Type) bool {
	return slices.Contains(structKinds, type_.Kind())
}

func isSlice(type_ reflect.Type) bool {
	return slices.Contains(sliceKinds, type_.Kind())
}
