package basic_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
)

func ParseNormal(val any, type_ reflect.Type) any {
	for reflect.TypeOf(val).Name() == "Value" {
		val = val.(reflect.Value).Interface()
	}
	kind := type_.Kind()
	if slices.Contains(intKinds, kind) {
		val = reflect.ValueOf(val).Convert(type_).Int()
	} else if slices.Contains(uintKinds, kind) {
		val = reflect.ValueOf(val).Convert(type_).Uint()
	} else if slices.Contains(floatKinds, kind) {
		val = reflect.ValueOf(val).Convert(type_).Float()
	}
	return val
}

func ParseStruct(val any, type_ reflect.Type) any {
	for reflect.TypeOf(val).Name() == "Value" {
		val = val.(reflect.Value).Interface()
	}
	dataMap := val.(map[string]any)
	fields := reflect.VisibleFields(type_)
	result := reflect.New(type_).Elem()
	for _, fieldInfo := range fields {
		fieldName := fieldInfo.Name
		dataVal, exists := dataMap[fieldName]
		if !exists { // do something?
			continue
		}
		parsedVal := Parse(dataVal, fieldInfo.Type)
		result.FieldByName(fieldName).Set(reflect.ValueOf(parsedVal))
	}
	return result.Interface()
}

func ParseList(val any, type_ reflect.Type) any {
	reflectVal := reflect.ValueOf(val)
	resList := reflect.MakeSlice(type_, reflectVal.Len(), reflectVal.Len())
	for i := range reflectVal.Len() {
		elem := reflectVal.Index(i)
		parsedVal := Parse(elem, type_.Elem())
		resList.Index(i).Set(reflect.ValueOf(parsedVal))
	}
	return resList
}

func Parse(val any, type_ reflect.Type) any {
	var result any

	switch type_.Kind() {
	case reflect.Array, reflect.Slice:
		result = ParseList(val, type_)
	case reflect.Struct:
		result = ParseStruct(val, type_)
	default:
		result = ParseNormal(val, type_)
	}
	fmt.Println("PARSED_RESULT: ", result)
	for reflect.TypeOf(result).Name() == "Value" {
		result = result.(reflect.Value).Interface()
	}
	return result
}

func Verify[T any](mapData map[string]any) T {
	res := Parse(mapData, reflect.TypeFor[T]())
	return res.(T)
}

func VerifyList[T any](mapData []map[string]any) []T {
	result := make([]T, 0)
	for _, i := range mapData {
		result = append(result, Verify[T](i))
	}
	return result
}

func VerifyBytes[T any](dataref *[]byte) any {
	data := *dataref
	var result any
	decoder := json.NewDecoder(bytes.NewReader(data))

	if data[0] == '[' {
		jsonArr := []map[string]any{}
		err := decoder.Decode(&jsonArr)
		if err != nil {
			fmt.Printf("ERRROR PARSING: %v\n", err)
		}
		result = VerifyList[T](jsonArr)
	} else if data[0] == '{' {
		jsonData := map[string]any{}
		err := decoder.Decode(&jsonData)
		if err != nil {
			fmt.Printf("ERRROR PARSING: %v\n", err)
		}
		result = Verify[T](jsonData)
	}
	return result.([]T)
}
