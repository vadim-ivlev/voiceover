package utils

import (
	"encoding/json"
	"fmt"
)

// PrettyJSON - returns a pretty-printed JSON string
func PrettyJSON(params interface{}) string {
	bytes, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	return string(bytes)
}
