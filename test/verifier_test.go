package basic_api_test

import (
	"encoding/json"
	"fmt"
	"testing"

	basic_api "github.com/onefivezero/basic-api"
)

type L3 struct {
	Sub_sub_info int64
	Flooaat      float64
	Buul         bool
	Listofnums   []int64
}

type L2 struct {
	Sub_info []L3
}

type L1 struct {
	Name      string
	Info      L2
	Moreinfo  []L2
	Emptyinfo []L2 // Keep this empty
}

func TestVerifier(t *testing.T) {
	body, _ := json.Marshal([]map[string]any{
		{
			"Name": "name",
			"Info": map[string]any{
				"Sub_info": []map[string]any{
					{
						"Sub_sub_info": 1,
						"Flooaat":      1.2,
						"Buul":         true,
						"Listofnums":   []int{1, 2, 3},
					},
				},
			},
			"Moreinfo": []map[string]any{
				{
					"Sub_info": []map[string]any{
						{
							"Sub_sub_info": 1,
							"Flooaat":      1.2,
							"Buul":         true,
							"Listofnums":   []int{1, 2, 3},
						},
					},
				},
			},
			"Emptyinfo": []map[string]any{},
		},
	})
	result, errs := basic_api.VerifyBytes[L1](&body)
	if len(errs) != 0 {
		fmt.Println(errs)
		t.FailNow()
	}
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic("ASDASDASD")
	}
	fmt.Println(string(b))
}

func TestVerifierInvalid(t *testing.T) {
	mapData := []map[string]any{
		{
			"Name": "name",
			"Info": map[string]any{
				"Sub_info": []map[string]any{
					{
						"Sub_sub_info": 1,
						"Flooaat":      1.2,
						"Buul":         true,
						"Listofnums":   []int{1, 2, 3},
					},
				},
			},
			"Moreinfo": []map[string]any{
				{
					"Sub_info": []map[string]any{
						{
							"Sub_sub_info": 1,
							"Flooaat":      1.2,
							"Buul":         false,
							"Listofnums":   []any{1, 2, 3},
						},
					},
				},
			},
			"Emptyinfo": []map[string]any{},
		},
	}
	result, err := basic_api.VerifyList[L1](mapData)
	fmt.Println(err)
	fmt.Println(result)
}
