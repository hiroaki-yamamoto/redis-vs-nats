package data

import (
	"fmt"
	"time"
)

//Result is a structure to store min, max, and avg time.
type Result struct {
	Min time.Duration
	Max time.Duration
	Avg time.Duration
}

// String converts the result into string.
func (me Result) String() string {
	res := fmt.Sprintln("Min: ", me.Min.String())
	res += fmt.Sprintln("Max: ", me.Max.String())
	res += fmt.Sprintln("Avg: ", me.Max.String())
	return res
}
