package main

import (
	"encoding/json"
	"fmt"
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
		latencyResults[i], err = measureLatency(con, 100000)
		if err != nil {
			panic(err)
		}
		fmt.Println("...Done")
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