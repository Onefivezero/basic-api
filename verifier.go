package basic_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

func ParseNormal(val reflect.Value, type_ reflect.Type) reflect.Value {
	returnVal := val.Convert(type_)
	return returnVal
}

func ParseStruct(val reflect.Value, type_ reflect.Type) reflect.Value {
	var dataMap map[string]any = val.Interface().(map[string]any)
	fields := reflect.VisibleFields(type_)
	result := reflect.New(type_).Elem()
	for _, fieldInfo := range fields {
		fieldName := fieldInfo.Name
		dataVal, exists := dataMap[fieldName]
		if !exists { // do something?
			continue
		}
		parsedVal := Parse(reflect.ValueOf(dataVal), fieldInfo.Type)
		result.FieldByName(fieldName).Set(parsedVal)
	}
	return result
}

func ParseList(val reflect.Value, type_ reflect.Type) reflect.Value {
	resList := reflect.MakeSlice(type_, val.Len(), val.Len())
	for i := range val.Len() {
		elem := val.Index(i)
		parsedVal := Parse(elem.Elem(), type_.Elem())
		resList.Index(i).Set(parsedVal)
	}
	return resList
}

func Parse(val reflect.Value, type_ reflect.Type) reflect.Value {
	switch type_.Kind() {
	case reflect.Array, reflect.Slice:
		val = ParseList(val, type_)
	case reflect.Struct:
		val = ParseStruct(val, type_)
	default:
		val = ParseNormal(val, type_)
	}
	return val
}

func Verify[T any](mapData map[string]any) T {
	res := Parse(reflect.ValueOf(mapData), reflect.TypeFor[T]())
	return res.Interface().(T)
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
