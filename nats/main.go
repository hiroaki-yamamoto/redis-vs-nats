package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
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
		var sum time.Duration
		fmt.Printf("Starting Iteration %d", i)
		for j := 0; j < 250; j++ {
			res, err := measureLatency(con, 100000)
			if err != nil {
				panic(err)
			}
			sum += res
		}
		fmt.Println("...Done")
		latencyResults[i] = sum
	}
	sort.Slice(latencyResults, func(i, j int) bool {
		return int64(latencyResults[i]) < int64(latencyResults[j])
	})
	res := data.Result{
		Min: latencyResults[0],
		Max: latencyResults[len(latencyResults)-1],
		Avg: func() time.Duration {
			var res = big.NewFloat(0)
			for _, v := range latencyResults {
				res = res.Add(res, big.NewFloat(float64(v)))
			}
			res.Quo(res, big.NewFloat(float64(len(latencyResults))))
			v, _ := res.Int64()
			return time.Duration(v)
		}(),
	}
	fmt.Print(res.String())
	const fname = "/opt/code/nats.json"
	f, err := os.Create(fname)
	if err != nil {
		panic(nil)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err := enc.Encode(res); err != nil {
		panic(err)
	}
	println("The result file is saved at: ", fname)
}
