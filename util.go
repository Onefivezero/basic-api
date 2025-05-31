package basic_api

import (
	"bytes"
	"encoding/json"
)

func bytesToListMap[T any](data *[]byte) ([]map[string]any, error) {
	decoder := json.NewDecoder(bytes.NewReader(*data))
	jsonArr := []map[string]any{}
	err := decoder.Decode(&jsonArr)
	return jsonArr, err
}

func bytesToMap[T any](data *[]byte) (map[string]any, error) {
	decoder := json.NewDecoder(bytes.NewReader(*data))
	jsonArr := map[string]any{}
	err := decoder.Decode(&jsonArr)
	return jsonArr, err
}
