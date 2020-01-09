package data

import (
	"fmt"
	"time"
)

//Result is a structure to store min, max, and avg time.
type Result struct {
	Min time.Duration `json:"min"`
	Max time.Duration `json:"max"`
	Avg time.Duration `json:"avg"`
}

// String converts the result into string.
func (me Result) String() string {
	res := fmt.Sprintln("Min: ", me.Min.String())
	res += fmt.Sprintln("Max: ", me.Max.String())
	res += fmt.Sprintln("Avg: ", me.Avg.String())
	return res
}
