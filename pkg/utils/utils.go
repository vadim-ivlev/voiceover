package utils

import (
	"encoding/json"
	"fmt"
)

// CompactJSON - returns a compact JSON string
func CompactJSON(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	return string(bytes)
}

// PrettyJSON - returns a pretty-printed JSON string
func PrettyJSON(params interface{}) string {
	bytes, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	return string(bytes)
}
