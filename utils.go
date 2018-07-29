package main

import (
	"strconv"
)

func isDigit(char string) bool {
	if _, err := strconv.Atoi(char); err == nil {
		return true
	}
	return false
}
