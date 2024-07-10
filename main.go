package basic_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type Method string

const (
	GET    Method = "GET"
	HEAD   Method = "HEAD"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
)

func fillStruct(donor map[string][]string, receiver any) error {
	receiverVal := reflect.ValueOf(receiver).Elem()
	if receiverVal.Kind() != reflect.Struct {
		return fmt.Errorf("RECEIVER IS NOT A STRUCT: %v", receiverVal)
	}

	for _, structField := range reflect.VisibleFields(receiverVal.Type()) {
		if !structField.Type.AssignableTo(reflect.TypeOf("")) {
			return fmt.Errorf("not assignable to string: %s", structField.Name)
		}
		fieldName := structField.Name
		jsonFieldName := structField.Tag.Get("json")
		fieldValues := donor[jsonFieldName]
		if len(fieldValues) > 1 {
			return errors.New("list of query parameters not supported at this time")
		}
		var fieldValue any
		if len(fieldValues) == 1 {
			fieldValue = fieldValues[0]
			receiverVal.FieldByName(fieldName).SetString(fieldValue.(string))
		}
	}
	return nil
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func CustomHandler[
	QueryModelType interface{},
	RequestModelType interface{},
	ResponseModelType interface{},
](
	url string,
	method Method,
	inFunc func(*QueryModelType, *RequestModelType) *ResponseModelType,
	serveMux *http.ServeMux,
) {
	wrapperFunc := func(w http.ResponseWriter, rawRequest *http.Request) {
		inFuncRefl := reflect.TypeOf(inFunc)
		queryType := inFuncRefl.In(0).Elem()
		reqType := inFuncRefl.In(1).Elem()

		// create empty structs
		var requestData *RequestModelType = reflect.New(reqType).Interface().(*RequestModelType)
		var queryData *QueryModelType = reflect.New(queryType).Interface().(*QueryModelType)

		// fill query parameters
		err := fillStruct(rawRequest.URL.Query(), queryData)
		panicErr(err)

		// fill request data
		err = json.NewDecoder(rawRequest.Body).Decode(requestData)
		panicErr(err)

		// run code and return response
		response := inFunc(queryData, requestData)
		responseString, err := json.Marshal(response)
		if err != io.EOF {
			panicErr(err)
		}
		io.WriteString(w, string(responseString))
	}
	if serveMux == nil {
		http.HandleFunc(string(method)+" "+url, wrapperFunc)
	} else {
		serveMux.HandleFunc(string(method)+" "+url, wrapperFunc)
	}
}
