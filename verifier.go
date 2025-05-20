package basic_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type parseError struct {
	Path   string
	Reason string
}

type parser struct {
	Errors []parseError
}

func (p *parser) AddError(path []string, reason string) {
	err := parseError{
		Path:   strings.Join(path, "."),
		Reason: reason,
	}
	// DEBUG: panic(fmt.Sprintf("%v\n", err))
	p.Errors = append(p.Errors, err)
}

func (p *parser) ParseNormal(val reflect.Value, type_ reflect.Type, path *[]string) reflect.Value {
	for !val.CanConvert(type_) && (val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface) {
		val = val.Elem()
	}
	if !val.CanConvert(type_) {
		p.AddError(*path, fmt.Sprintf("Value: %v %v can not be converted to Type: %v", val.Type(), val, type_))
		return reflect.Zero(type_)
	}
	return val.Convert(type_)
}

func (p *parser) ParseStruct(val reflect.Value, type_ reflect.Type, path *[]string) reflect.Value {
	dataMap, ok := val.Interface().(map[string]any)
	if !ok {
		p.AddError(*path, fmt.Sprintf("Value: %v %v can not be converted to Type: %v", val.Type(), val, type_))
		return reflect.ValueOf(nil)
	}

	fields := reflect.VisibleFields(type_)
	result := reflect.New(type_).Elem()

	for _, fieldInfo := range fields {
		fieldName := fieldInfo.Name
		dataVal, exists := dataMap[fieldName]
		if !exists { // TODO: Add to errors if nil not allowed.
			continue
		}
		path_f := append(*path, fieldName)
		parsedVal := p.Parse(reflect.ValueOf(dataVal), fieldInfo.Type, &path_f)
		result.FieldByName(fieldName).Set(parsedVal)
	}
	return result
}

func (p *parser) ParseList(val reflect.Value, type_ reflect.Type, path *[]string) reflect.Value {
	if !isSlice(val.Type()) {
		p.AddError(*path, fmt.Sprintf("Value: %v %v can not be converted to Type: %v", val.Type(), val, type_))
		return reflect.ValueOf(nil)
	}
	resList := reflect.MakeSlice(type_, val.Len(), val.Len())
	for i := range val.Len() {
		elem := val.Index(i)
		path_i := append(*path, fmt.Sprintf("[%d]", i))
		parsedVal := p.Parse(elem, type_.Elem(), &path_i)
		if parsedVal.IsValid() {
			resList.Index(i).Set(parsedVal)
		}
	}
	return resList
}

func (p *parser) Parse(val reflect.Value, type_ reflect.Type, path *[]string) reflect.Value {
	switch type_.Kind() {
	case reflect.Array, reflect.Slice:
		val = p.ParseList(val, type_, path)
	case reflect.Struct:
		val = p.ParseStruct(val, type_, path)
	default:
		val = p.ParseNormal(val, type_, path)
	}
	return val
}

func Verify[T any](mapData map[string]any) (T, []parseError) {
	p := parser{Errors: []parseError{}}
	e := make([]string, 0)
	res := p.Parse(reflect.ValueOf(mapData), reflect.TypeFor[T](), &e)

	if len(p.Errors) > 0 {
		return *new(T), p.Errors
	}
	return res.Interface().(T), p.Errors
}

func VerifyList[T any](mapData []map[string]any) ([]T, []parseError) {
	result := make([]T, 0)
	errList := make([]parseError, 0)
	for _, i := range mapData {
		parsedObj, errs := Verify[T](i)
		result = append(result, parsedObj)
		errList = append(errList, errs...)
	}
	// TODO: Fix parseError paths. They should start with [0] instead of nothing.
	return result, errList
}

func VerifyBytes[T any](dataref *[]byte) (any, []parseError) {
	data := *dataref
	var result any
	var errs []parseError
	decoder := json.NewDecoder(bytes.NewReader(data))

	if data[0] == '[' {
		jsonArr := []map[string]any{}
		err := decoder.Decode(&jsonArr)
		if err != nil {
			fmt.Printf("ERRROR PARSING: %v\n", err)
		}
		result, errs = VerifyList[T](jsonArr)
	} else if data[0] == '{' {
		jsonData := map[string]any{}
		err := decoder.Decode(&jsonData)
		if err != nil {
			fmt.Printf("ERRROR PARSING: %v\n", err)
		}
		result, errs = Verify[T](jsonData)
	}
	return result.([]T), errs
}
