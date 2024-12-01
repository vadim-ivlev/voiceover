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

// CalcStartEndIndex calculates and adjusts the start and end indices based on the total number of lines.
// It ensures that the indices are within valid bounds and that the start index is less than the end index.
//
// Parameters:
//   - totalLines: The total number of lines available.
//   - startIndex: The initial start index.
//   - endIndex: The initial end index.
//
// Returns:
//   - start: The adjusted start index.
//   - end: The adjusted end index.
//   - err: An error if the start index is greater than or equal to the total number of lines, or if the start index is greater than or equal to the end index.
func CalcStartEndIndex(totalLines, startIndex, endIndex int) (start, end int, err error) {

	start = startIndex
	end = endIndex

	// adjust the end index
	if end < 0 {
		end = 0
	}
	if end > totalLines {
		end = totalLines
	}

	// adjust the start index
	if start < 0 {
		start = 0
	}

	if start >= totalLines {
		return 0, 0, fmt.Errorf("start index %d is greater than or equal to the total number of lines %d", start, totalLines)
	}

	// check if the start index is greater than the end index
	if start >= end {
		return 0, 0, fmt.Errorf("start index %d is greater than or equal to the end index %d", start, end)
	}

	return start, end, nil
}
