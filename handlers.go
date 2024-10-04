package basic_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
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

func panicErr(w http.ResponseWriter, err error) {
	response := ErrorResponse{
		StatusCode:   400,
		ErrorMessage: &map[string]string{"Error": err.Error()},
	}
	responseBytes, _ := json.Marshal(response)
	io.WriteString(w, string(responseBytes))
}

func CustomHandler[
	QueryModelType interface{},
	RequestModelType interface{},
	ResponseModelType interface{},
](
	url string,
	method Method,
	inFunc func(*QueryModelType, *RequestModelType) (*ResponseModelType, *ErrorResponse),
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
		if err != nil {
			panicErr(w, err)
			return
		}
		// fill request data
		err = json.NewDecoder(rawRequest.Body).Decode(requestData)
		if err != nil {
			panicErr(w, err)
			return
		}
		// run code and return response
		response, http_error := inFunc(queryData, requestData)

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
