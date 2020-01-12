package main

import (
	"context"
	"encoding/json"
	"fmt"
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
		fmt.Printf("Starting Iteration %d", i)
		ctx, stop := context.WithTimeout(context.Background(), 10*time.Second)
		defer stop()
		res, err := measureLatency(ctx, con, 100000)
		if err != nil {
			panic(err)
		}
		latencyResults[i] = res
		fmt.Println("...Done")
	}
	sort.Slice(latencyResults, func(i, j int) bool {
		return int64(latencyResults[i]) < int64(latencyResults[j])
	})
	res := &data.Result{}
	res.SetData(latencyResults)
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
