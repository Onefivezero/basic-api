package parser

import (
	"fmt"
	"reflect"
	"strings"
)

type ParseError struct {
	Path   string
	Reason string
}

type parser struct {
	Errors []ParseError
}

func (p *parser) AddError(path []string, reason string) {
	err := ParseError{
		Path:   strings.Join(path, "."),
		Reason: reason,
	}
	// DEBUG: panic(fmt.Sprintf("%v\n", err))
	p.Errors = append(p.Errors, err)
}

func (p *parser) parseNormal(val reflect.Value, type_ reflect.Type, path *[]string) reflect.Value {
	for !val.CanConvert(type_) && (val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface) {
		val = val.Elem()
	}
	if !val.CanConvert(type_) {
		p.AddError(*path, fmt.Sprintf("Value: %v %v can not be converted to Type: %v", val.Type(), val, type_))
		return reflect.Zero(type_)
	}
	return val.Convert(type_)
}

func (p *parser) parseStruct(val reflect.Value, type_ reflect.Type, path *[]string) reflect.Value {
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
		if parsedVal.IsValid() {
			result.FieldByName(fieldName).Set(parsedVal)
		}
	}
	return result
}

func (p *parser) parseList(val reflect.Value, type_ reflect.Type, path *[]string) reflect.Value {
	if !IsSlice(val.Type()) {
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
		val = p.parseList(val, type_, path)
	case reflect.Struct:
		val = p.parseStruct(val, type_, path)
	default:
		val = p.parseNormal(val, type_, path)
	}
	return val
}

func Parse(data any, type_ reflect.Type) (any, []ParseError) {
	p := parser{Errors: []ParseError{}}
	e := make([]string, 0)
	res := p.Parse(reflect.ValueOf(data), type_, &e)
	if p.Errors == nil || len(p.Errors) > 0 {
		return reflect.Zero(type_), p.Errors
	}
	return res.Interface(), p.Errors
}
