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

func parseStruct(source map[string]any, targetType reflect.Type) {
	for _, fieldInfo := range reflect.VisibleFields(targetType) {
		fieldName := fieldInfo.Name
		val, exists := source[fieldName]
		fmt.Printf("STRUCT: field: %v, val: %v, exists: %v\n", fieldName, val, exists)
		parse(val, fieldInfo.Type)
	}
}

func parseList(source any, targetInnerType reflect.Type) {
	rx := reflect.ValueOf(source)
	for i := range rx.Len() {
		val := rx.Index(i).Interface()
		fmt.Printf("LIST: %v", val)
		parse(val, targetInnerType)
	}
}

func parseNormal(source any, targetType reflect.Type) {
	sourceType := reflect.TypeOf(source)
	if !sourceType.ConvertibleTo(targetType) {
		perr(fmt.Errorf("source: %v target: %v", sourceType, targetType))
	}
	fmt.Println("NORMAL: ", source)
}

func parse(source any, targetType reflect.Type) {
	// sourceType := reflect.TypeOf(source)
	// sourceKind := sourceType.Kind()
	targetKind := targetType.Kind()
	switch targetKind {
	case reflect.Array, reflect.Slice:
		parseList(source, targetType.Elem())
	case reflect.Struct:
		parseStruct(source.(map[string]any), targetType)
	default:
		parseNormal(source, targetType)
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
