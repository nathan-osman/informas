package server

import (
	"strconv"
)

// atoi converts a string to its integer value and fails silently when an error
// occurs.
func atoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		v = 0
	}
	return v
}
