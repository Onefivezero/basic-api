package basic_api

import (
	"log"
	"net/http"
)

// Define API models first and annotate with json struct tags.

type ExampleQueryModel struct {
	String string `json:"queryString"`
}

type ExampleRequestBodyModel struct {
	Number int  `json:"number"`
	Bool   bool `json:"bool"`
}

type ExampleResponseModel struct {
	QueryString   string `json:"queryString"`
	RequestNumber int    `json:"requestNumber"`
	RequestBool   bool   `json:"requestBool"`
}

// Define an endpoint function
// The first parameter's type defines the url parameters, the second defines the request body structure, and finally the return type defines the response model.

func ExampleEndpoint(queryData *ExampleQueryModel, requestData *ExampleRequestBodyModel) *ExampleResponseModel {
	return &ExampleResponseModel{
		QueryString:   queryData.String,
		RequestNumber: requestData.Number,
		RequestBool:   requestData.Bool,
	}
}

// Call the CustomHandler function with url path and method, and finally start the server.

func init() {
	CustomHandler("/example", "POST", ExampleEndpoint)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
