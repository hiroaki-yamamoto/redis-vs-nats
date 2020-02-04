package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/hiroaki-yamamoto/redis-vs-nats/config"
	d "github.com/hiroaki-yamamoto/redis-vs-nats/data"
	natsBench "github.com/hiroaki-yamamoto/redis-vs-nats/nats"
	redisBench "github.com/hiroaki-yamamoto/redis-vs-nats/redis"
	"github.com/nats-io/nats.go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var rootFlags *flag.FlagSet

const natsAddr = "nats://nats:4222"
const redisAddr = "redis:6379"

// BenchmarkInterface indicates an interface to measure message query.
type BenchmarkInterface interface {
	Measure() (dur []time.Duration, err error)
}

func init() {
	rootFlags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	rootFlags.String(
		"f", "/etc/redis-vs-nats/nats-single.yml",
		"The path of the configuration file",
	)
}

func main() {
	pflag.CommandLine.AddGoFlagSet(rootFlags)
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}
	cfgPath := viper.GetString("f")
	cfg, err := config.New(cfgPath)
	if err != nil {
		panic(err)
	}
	rootCtx := context.Background()
	var bench BenchmarkInterface
	switch cfg.Target {
	case "nats":
		if conn, err := nats.Connect(natsAddr); err == nil {
			bench = &natsBench.Benchmark{
				Ctx:      rootCtx,
				NumItr:   cfg.NumIteration,
				NumTrial: cfg.NumTrial,
				BufSize:  cfg.BufSize,
				Con:      conn,
			}
		} else {
			panic(err)
		}
		break
	case "redis":
		cli := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: "",
			DB:       0,
		})
		if _, err := cli.Ping().Result(); err != nil {
			panic(err)
		}
		bench = &redisBench.Benchmark{
			Ctx:      rootCtx,
			Con:      cli,
			NumItr:   cfg.NumIteration,
			NumTrial: cfg.NumTrial,
			BufSize:  cfg.BufSize,
		}
		break
	default:
		panic(fmt.Sprintf("Unspecified Target: %s", cfg.Target))
	}
	fmt.Println("=====Config=====")
	if cfgPret, err := json.MarshalIndent(cfg, "", "  "); err == nil {
		fmt.Println(string(cfgPret))
	} else {
		panic(err)
	}
	fmt.Println("================")
	if dur, err := bench.Measure(); err == nil {
		res := &d.Result{}
		res.SetData(dur)
		fmt.Println(res.String())
	} else {
		panic(err)
	}
}
