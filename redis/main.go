package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/hiroaki-yamamoto/redis-vs-nats/data"
)

func main() {
	rand.Seed(time.Now().UTC().Unix())
	con := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	_, err := con.Ping().Result()
	if err != nil {
		panic(err)
	}
	defer con.Close()
	if err != nil {
		panic(err)
	}
	latencyResults := make([]time.Duration, 2000)
	for i := range latencyResults {
		fmt.Printf("Starting Iteration %d", i)
		var sum time.Duration
		for j := 0; j < 250; j++ {
			var mres time.Duration
			mres, err = measureLatency(con, 100000)
			if err != nil {
				panic(err)
			}
			sum += mres
		}
		latencyResults[i] = sum
		fmt.Println("...Done")
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
	const fname = "/opt/code/redis.json"
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
