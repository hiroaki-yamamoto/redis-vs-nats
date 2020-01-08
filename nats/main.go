package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/hiroaki-yamamoto/redis-vs-nats/data"
	"github.com/nats-io/nats.go"
)

func main() {
	rand.Seed(time.Now().UTC().Unix())
	con, err := nats.Connect("nats://nats:4222")
	defer con.Close()
	if err != nil {
		panic(err)
	}
	latencyResults := make([]time.Duration, 2000)
	for i := range latencyResults {
		fmt.Printf("Starting Iteration %d", i)
		latencyResults[i], err = measureLatency(con, 100000)
		if err != nil {
			panic(err)
		}
		fmt.Printf("...Done\n")
	}
	sort.Slice(latencyResults, func(i, j int) bool {
		return int64(latencyResults[i]) < int64(latencyResults[j])
	})
	res := data.Result{
		Min: latencyResults[0],
		Max: latencyResults[len(latencyResults)-1],
		Avg: func() time.Duration {
			var res time.Duration
			for _, v := range latencyResults {
				res += v
			}
			res /= time.Duration(len(latencyResults))
			return res
		}(),
	}
	fmt.Print(res.String())
}
