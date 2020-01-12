package data

import (
	"fmt"
	"sort"
	"time"
)

//Result is a structure to store min, max, and avg time.
type Result struct {
	Min time.Duration `json:"min"`
	Max time.Duration `json:"max"`
	Sum time.Duration `json:"sum"`
	Avg time.Duration `json:"avg"`
}

// SetData set the delay duration data and calculate stats.
func (me *Result) SetData(data []time.Duration) {
	sort.Slice(data, func(i, j int) bool {
		return int64(data[i]) < int64(data[j])
	})
	me.Min = data[0]
	me.Max = data[len(data)-1]
	me.Sum = func() time.Duration {
		var res int64
		for _, v := range data {
			res += int64(v)
		}
		return time.Duration(time.Duration(res))
	}()
	me.Avg = time.Duration(int64(me.Sum) / int64(len(data)))
}

// String converts the result into string.
func (me Result) String() string {
	res := fmt.Sprintln("Min: ", me.Min.String())
	res += fmt.Sprintln("Max: ", me.Max.String())
	res += fmt.Sprintln("Sum: ", me.Sum.String())
	res += fmt.Sprintln("Avg: ", me.Avg.String())
	return res
}
