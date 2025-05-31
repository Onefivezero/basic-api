package basic_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/onefivezero/basic-api/parser"
)

func panicErr(w http.ResponseWriter, err error) {
	response := ErrorResponse{
		StatusCode:   400,
		ErrorMessage: &map[string]string{"Error": err.Error()},
	}
	responseBytes, _ := json.Marshal(response)
	io.WriteString(w, string(responseBytes))
}

func returnValidationErr(w http.ResponseWriter, parseErrors []parser.ParseError) {
	response := ErrorResponse{
		StatusCode:   400,
		ErrorMessage: parseErrors,
	}
	responseBytes, _ := json.Marshal(response)
	io.WriteString(w, string(responseBytes))
}

func fillQueryParams[T any](w http.ResponseWriter, mapData map[string][]string) (T, []parser.ParseError) {
	val := new(T)
	type_ := reflect.TypeFor[T]()
	parseErrs := make([]parser.ParseError, 0)

	if type_.Kind() != reflect.Struct {
		panicErr(w, fmt.Errorf("Query struct not param. Invalid.")) // TODO: Proper solution. Check while wrapping the func, perhaps?
	}
	fields := reflect.VisibleFields(type_)
	for _, field := range fields {
		fieldName := field.Name
		data, exists := mapData[strings.ToLower(fieldName)]
		if !exists {
			continue
		}
		if parser.IsSlice(field.Type) {
			reflect.ValueOf(val).Elem().FieldByName(fieldName).Set(reflect.ValueOf(data))
		} else { // Assume its normal: not struct, not list, etc.
			if len(data) > 1 {
				panicErr(w, fmt.Errorf("Too many vals")) // TODO: Proper solution. Check while wrapping the func, perhaps?
			} else {
				parsedVal, parseErrsss := parser.Parse(data[0], field.Type)
				if parseErrs != nil && len(parseErrs) > 0 {
					parseErrs = append(parseErrs, parseErrsss...)
				}
				reflect.ValueOf(val).Elem().FieldByName(fieldName).Set(reflect.ValueOf(parsedVal))
			}
		}
	}
	return *val, parseErrs
}

func fillData[T any](mapData any) (T, []parser.ParseError) {
	result, parseErrs := parser.Parse(mapData, reflect.TypeFor[T]())
	if parseErrs != nil && len(parseErrs) > 0 {
		return *new(T), parseErrs
	}
	typedResult, ok := result.(T)
	if !ok {
		fmt.Println("ERROR: CASTING ERROR FOR SOME REASON.")
		fmt.Printf("%v\n", result)
	}
	return typedResult, nil
}

func fillBody[T any](data io.ReadCloser) (T, []parser.ParseError) {
	decoder := json.NewDecoder(data)
	var mapData any
	decoder.Decode(&mapData)
	return fillData[T](mapData)
}

func CustomHandler[
	QueryModelType any,
	RequestModelType any,
	ResponseModelType any,
](
	url string,
	method Method,
	inFunc func(*QueryModelType, *RequestModelType) (*ResponseModelType, *ErrorResponse),
	serveMux *http.ServeMux,
) {
	wrapperFunc := func(w http.ResponseWriter, rawRequest *http.Request) {
		// create empty structs
		var requestData RequestModelType
		var queryData QueryModelType
		var err error

		// fill query parameters
		queryData, parseErrors := fillQueryParams[QueryModelType](w, rawRequest.URL.Query())
		if parseErrors != nil && len(parseErrors) > 0 {
			returnValidationErr(w, parseErrors)
			return
		}

		// fill body
		requestData, parseErrors = fillBody[RequestModelType](rawRequest.Body)
		if parseErrors != nil && len(parseErrors) > 0 {
			returnValidationErr(w, parseErrors)
			return
		}
		// run code and return response
		response, http_error := inFunc(&queryData, &requestData)

		// if response, return 200 and response. Else, return custom error.
		var responseString []byte
		if response != nil {
			responseString, err = json.Marshal(response)
			if err != nil && err != io.EOF {
				panicErr(w, err)
				return
			}
		} else if http_error != nil {
			responseString, err = json.Marshal(http_error)
			if err != nil && err != io.EOF {
				panicErr(w, err)
				return
			}
		}

		io.WriteString(w, string(responseString))
	}
	if serveMux == nil {
		http.HandleFunc(string(method)+" "+url, wrapperFunc)
	} else {
		serveMux.HandleFunc(string(method)+" "+url, wrapperFunc)
	}
}
