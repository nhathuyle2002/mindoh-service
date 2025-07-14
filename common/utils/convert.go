package utils

import "fmt"

func ParseUint(s string) uint {
	var i uint
	fmt.Sscanf(s, "%d", &i)
	return i
}
