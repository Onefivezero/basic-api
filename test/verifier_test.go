package basic_api_test

import (
	"encoding/json"
	"reflect"
	"testing"

	basic_api "github.com/onefivezero/basic-api"
)

type L3 struct {
	sub_sub_info int
	flooat       float64
	buul         bool
	listofnums   []int
}

type L2 struct {
	sub_info []L3
}

type L1 struct {
	name      string
	info      L2
	moreinfo  []L2
	emptyinfo []L2 // Keep this empty
}

func TestVerifier(t *testing.T) {
	body, _ := json.Marshal([]map[string]any{
		{
			"name": "name",
			"info": map[string]any{
				"sub_info": []map[string]any{
					{
						"sub_sub_info": 1,
						"flooat":       1.2,
						"buul":         true,
						"listofnums":   []int{1, 2, 3},
					},
				},
			},
			"moreinfo": []map[string]any{
				{
					"sub_info": []map[string]any{
						{
							"sub_sub_info": 1,
							"flooat":       1.2,
							"buul":         true,
							"listofnums":   []int{1, 2, 3},
						},
					},
				},
			},
			"emptyinfo": []map[string]any{},
		},
	})
	basic_api.VerifyBytes(body, reflect.TypeFor[[]L1]())
}
