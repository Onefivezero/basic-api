# Basic API

A Golang REST framework made to simplify writing REST APIs.

## Basic Example

```Go
...

// Define request data and query parameter models.
type StudentInfo struct {
	Name        string
	Age         int64
	Score       float64
	LetterScore rune
	Passed      bool
}

type StudentIdentifierInfo struct {
	Id string
}

type StudentCompleteInfo struct {
	Id          string
	Name        string
	Age         int64
	Score       float64
	LetterScore rune
	Passed      bool
}

// Define an endpoint function.
func CombineStudentInfo(
	queryParameters *StudentIdentifierInfo,
	requestData *StudentInfo,
) (*StudentCompleteInfo, *basic_api.ErrorResponse) {
	return &StudentCompleteInfo{
		Id:          queryParameters.Id,
		Name:        requestData.Name,
		Age:         requestData.Age,
		Score:       requestData.Score,
		LetterScore: requestData.LetterScore,
		Passed:      requestData.Passed,
	}
}


// Create a handler object, call the CustomHandler function to add the endpoint to this handler, and finally start listening.
// mux is optional and can be replaced with nil, in which case the default ServeMux in the http package will be used instead.
func main() {
	mux := http.NewServeMux()
	basic_api.CustomHandler("/combine", "POST", CombineStudentInfo, mux)

	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
```
