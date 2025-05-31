package basic_api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	basic_api "github.com/onefivezero/basic-api"
	"github.com/stretchr/testify/assert"
)

func FailTestIfErrorNotNil(t *testing.T, err error) {
	if err != nil {
		t.Fatal(fmt.Errorf(err.Error()))
	}
}

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
	}, nil
}

func CombineStudentInfoMult(
	queryParameters *StudentIdentifierInfo,
	requestData *[]StudentInfo,
) (*StudentCompleteInfo, *basic_api.ErrorResponse) {
	studentInfo := (*requestData)[0]
	return &StudentCompleteInfo{
		Id:          queryParameters.Id,
		Name:        studentInfo.Name,
		Age:         studentInfo.Age,
		Score:       studentInfo.Score,
		LetterScore: studentInfo.LetterScore,
		Passed:      studentInfo.Passed,
	}, nil
}

func TestBasicUsage(t *testing.T) {
	mux := http.NewServeMux()
	basic_api.CustomHandler("/combine", "POST", CombineStudentInfo, mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	studentInfo := StudentInfo{
		Name:        "Name",
		Age:         20,
		Score:       90.3,
		LetterScore: 'A',
		Passed:      true,
	}
	studentInfoJson, err := json.Marshal(studentInfo)
	FailTestIfErrorNotNil(t, err)

	res, err := http.Post(server.URL+"/combine?id=studentid", "application/json", bytes.NewBuffer(studentInfoJson))
	FailTestIfErrorNotNil(t, err)

	bodyByte, err := io.ReadAll(res.Body)
	FailTestIfErrorNotNil(t, err)

	var body StudentCompleteInfo = StudentCompleteInfo{}
	err = json.Unmarshal(bodyByte, &body)
	FailTestIfErrorNotNil(t, err)

	expectedBody := StudentCompleteInfo{
		Id:          "studentid",
		Name:        "Name",
		Age:         20,
		Score:       90.3,
		LetterScore: 'A',
		Passed:      true,
	}
	assert.Equal(t, expectedBody, body, string(bodyByte))
}

func TestBadRequest(t *testing.T) {
	mux := http.NewServeMux()
	basic_api.CustomHandler("/combine", "POST", CombineStudentInfo, mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	studentInfo := map[string]any{
		"Name":        "Name",
		"Age":         "Adult",
		"Score":       "Enough",
		"LetterScore": 1,
		"Passed":      "maybe",
	}
	studentInfoJson, err := json.Marshal(studentInfo)
	FailTestIfErrorNotNil(t, err)

	res, err := http.Post(server.URL+"/combine?id=studentid", "application/json", bytes.NewBuffer(studentInfoJson))
	FailTestIfErrorNotNil(t, err)

	bodyByte, err := io.ReadAll(res.Body)
	FailTestIfErrorNotNil(t, err)

	var body basic_api.ErrorResponse = basic_api.ErrorResponse{}
	err = json.Unmarshal(bodyByte, &body)
	FailTestIfErrorNotNil(t, err)

	expectedErrorBody := basic_api.ErrorResponse{
		StatusCode: 400,
		ErrorMessage: []any{
			map[string]any{
				"Path":   "Age",
				"Reason": "Value: string Adult can not be converted to Type: int64",
			},
			map[string]any{
				"Path":   "Score",
				"Reason": "Value: string Enough can not be converted to Type: float64",
			},
			map[string]any{
				"Path":   "Passed",
				"Reason": "Value: string maybe can not be converted to Type: bool",
			},
		},
	}

	assert.Equal(t, expectedErrorBody, body)
}

func TestBadRequestList(t *testing.T) {
	mux := http.NewServeMux()
	basic_api.CustomHandler("/combine", "POST", CombineStudentInfoMult, mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	studentInfo := map[string]any{
		"Name":        "Name",
		"Age":         "Adult",
		"Score":       "Enough",
		"LetterScore": 1,
		"Passed":      "maybe",
	}
	studentInfoJson, err := json.Marshal([]map[string]any{studentInfo})
	FailTestIfErrorNotNil(t, err)

	res, err := http.Post(server.URL+"/combine?id=studentid", "application/json", bytes.NewBuffer(studentInfoJson))
	FailTestIfErrorNotNil(t, err)

	bodyByte, err := io.ReadAll(res.Body)
	FailTestIfErrorNotNil(t, err)

	var body basic_api.ErrorResponse = basic_api.ErrorResponse{}
	err = json.Unmarshal(bodyByte, &body)
	FailTestIfErrorNotNil(t, err)

	expectedErrorBody := basic_api.ErrorResponse{
		StatusCode: 400,
		ErrorMessage: []any{
			map[string]any{
				"Path":   "[0].Age",
				"Reason": "Value: string Adult can not be converted to Type: int64",
			},
			map[string]any{
				"Path":   "[0].Score",
				"Reason": "Value: string Enough can not be converted to Type: float64",
			},
			map[string]any{
				"Path":   "[0].Passed",
				"Reason": "Value: string maybe can not be converted to Type: bool",
			},
		},
	}

	assert.Equal(t, expectedErrorBody, body)
}
