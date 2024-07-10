package basic_api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	basic_api "github.com/onefivezero/basic-api"
)

func FailTestIfErrorNotNil(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

type StudentInfo struct {
	Name        string  `json:"name"`
	Age         int     `json:"age"`
	Score       float32 `json:"score"`
	LetterScore rune    `json:"letterScore"`
	Passed      bool    `json:"passed"`
}

type StudentIdentifierInfo struct {
	Id string `json:"id"`
}

type StudentCompleteInfo struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Age         int     `json:"age"`
	Score       float32 `json:"score"`
	LetterScore rune    `json:"letterScore"`
	Passed      bool    `json:"passed"`
}

func CombineStudentInfo(
	queryParameters *StudentIdentifierInfo,
	requestData *StudentInfo,
) *StudentCompleteInfo {
	return &StudentCompleteInfo{
		Id:          queryParameters.Id,
		Name:        requestData.Name,
		Age:         requestData.Age,
		Score:       requestData.Score,
		LetterScore: requestData.LetterScore,
		Passed:      requestData.Passed,
	}
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
	if body != expectedBody {
		t.Fatal("expected body not found.")
	}
}
