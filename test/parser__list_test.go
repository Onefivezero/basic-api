package basic_api_test

import (
	"reflect"
	"testing"

	"github.com/onefivezero/basic-api/parser"
	ts "github.com/onefivezero/basic-api/test/structs"
	"github.com/stretchr/testify/assert"
)

func TestParseList(t *testing.T) {
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
	result, err := parser.Parse(mapData, reflect.SliceOf(reflect.TypeFor[ts.L1]()))

	assert.Empty(t, err)
	assert.Equal(t, result, []ts.L1{
		{
			Name: "name",
			Info: ts.L2{
				Sub_info: []ts.L3{
					{
						Sub_sub_info: 1,
						Flooaat:      1.2,
						Buul:         true,
						Listofnums:   []int64{1, 2, 3},
					},
				},
			},
			Moreinfo: []ts.L2{
				{
					Sub_info: []ts.L3{
						{
							Sub_sub_info: 1,
							Flooaat:      1.2,
							Buul:         false,
							Listofnums:   []int64{1, 2, 3},
						},
					},
				},
			},
			Emptyinfo: []ts.L2{},
		},
		{
			Name: "name",
			Info: ts.L2{
				Sub_info: []ts.L3{
					{
						Sub_sub_info: 1,
						Flooaat:      1.2,
						Buul:         true,
						Listofnums:   []int64{1, 2, 3},
					},
				},
			},
			Moreinfo: []ts.L2{
				{
					Sub_info: []ts.L3{
						{
							Sub_sub_info: 1,
							Flooaat:      1.2,
							Buul:         false,
							Listofnums:   []int64{1, 2, 3},
						},
					},
				},
			},
			Emptyinfo: []ts.L2{},
		},
	})
}

func TestParseListInvalid(t *testing.T) {
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
			"Moreinfo": map[string]any{
				"Sub_info": []map[string]any{
					{
						"Sub_sub_info": 1,
						"Flooaat":      1.2,
						"Buul":         false,
						"Listofnums":   []any{1, 2, 3},
					},
				},
			},
			"Emptyinfo": []int{1, 2, 3},
		},
	}
	result, err := parser.Parse(mapData, reflect.SliceOf(reflect.TypeFor[ts.L1]()))

	assert.ElementsMatch(
		t,
		err,
		[]parser.ParseError{
			{
				Path:   "[1].Moreinfo",
				Reason: "Value: map[string]interface {} map[Sub_info:[map[Buul:false Flooaat:1.2 Listofnums:[1 2 3] Sub_sub_info:1]]] can not be converted to Type: []test_parser_structs.L2",
			},
			{
				Path:   "[1].Emptyinfo.[0]",
				Reason: "Value: int 1 can not be converted to Type: test_parser_structs.L2",
			},
			{
				Path:   "[1].Emptyinfo.[1]",
				Reason: "Value: int 2 can not be converted to Type: test_parser_structs.L2",
			},
			{
				Path:   "[1].Emptyinfo.[2]",
				Reason: "Value: int 3 can not be converted to Type: test_parser_structs.L2",
			},
		},
	)
	assert.Equal(t, reflect.Zero(reflect.TypeFor[[]ts.L1]()), result)
}
