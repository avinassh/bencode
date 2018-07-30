package main

import (
	"encoding/json"
	"strconv"
)

func isDigit(char string) bool {
	if _, err := strconv.Atoi(char); err == nil {
		return true
	}
	return false
}

// always send a pointer
func jsonMustMarshal(v interface{}) []byte {
	resp, _ := json.Marshal(v)
	return resp
}
