package basic_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

func perr(err error) {
	fmt.Println(err)
	panic("ENDED")
}

func parseStruct(source map[string]any, targetType reflect.Type) any {
	finalVal := reflect.New(targetType)
	for _, fieldInfo := range reflect.VisibleFields(targetType) {
		fieldName := fieldInfo.Name
		val, exists := source[fieldName]
		fmt.Printf("STRUCT: field: %v, val: %v, exists: %v\n", fieldName, val, exists)
		if exists {
			parsedVal := parse(val, fieldInfo.Type)
			fmt.Printf("1: %v", reflect.ValueOf(finalVal).FieldByName(fieldName))
			fmt.Printf("2: %v", reflect.ValueOf(finalVal).FieldByName(fieldName))
			fmt.Printf("3: %v", fieldName)
			reflect.ValueOf(finalVal).FieldByName(fieldName).Set(reflect.ValueOf(parsedVal))
		}
	}
	return finalVal
}

func parseList(source any, targetType reflect.Type) []any {
	targetInnerType := targetType.Elem()
	rx := reflect.ValueOf(source)
	var resultList []any = make([]any, 0)
	for i := range rx.Len() {
		val := rx.Index(i).Interface()
		Lval := parse(val, targetInnerType)
		resultList = append(resultList, Lval)
	}
	return resultList
}

func parseNormal(source any, targetType reflect.Type) any {
	sourceType := reflect.TypeOf(source)
	if !sourceType.ConvertibleTo(targetType) {
		perr(fmt.Errorf("source: %v target: %v", sourceType, targetType))
	}
	fmt.Println("NORMAL: ", source)
	return source
}

func parse(source any, targetType reflect.Type) any {
	// sourceType := reflect.TypeOf(source)
	// sourceKind := sourceType.Kind()
	targetKind := targetType.Kind()
	switch targetKind {
	case reflect.Array, reflect.Slice:
		return parseList(source, targetType)
	case reflect.Struct:
		return parseStruct(source.(map[string]any), targetType)
	default:
		return parseNormal(source, targetType)
	}
}

func Verify(source any, targetType reflect.Type) {
	parse(source, targetType)
}

func VerifyBytes(source []byte, targetType reflect.Type) {
	decoder := json.NewDecoder(bytes.NewReader(source))

	switch source[0] {
	case '[':
		listVal := make([]map[string]any, 0)
		err := decoder.Decode(&listVal)
		if err != nil {
			fmt.Println("LISTVAL NOT VALID")
			return
		}
		Verify(listVal, targetType)
	case '{':
		mapVal := make(map[string]any)
		err := decoder.Decode(&mapVal)
		if err != nil {
			fmt.Println("MAPVAL NOT VALID")
			fmt.Println(err)
			return
		}
		Verify(mapVal, targetType)
	default:
		fmt.Printf("Not valid: %c", source[0])
		return
	}
}
