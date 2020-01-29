package main

import (
	"flag"
	"os"
	"time"

	"github.com/hiroaki-yamamoto/redis-vs-nats/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var rootFlags *flag.FlagSet

// BenchmarkInterface indicates an interface to measure message query.
type BenchmarkInterface interface {
	Measure() (itrDur time.Duration, err error)
}

func init() {
	rootFlags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	rootFlags.String(
		"target", "nats", "[nats|redis] target message queue benchmark",
	)
	rootFlags.Int("numTrial", 1, "The number of trial per a iteration")
	rootFlags.Int(
		"numIteration", 2000,
		"The number of iteration. (i.e. the sum of trial is numTrial x inum)",
	)
	rootFlags.Int("bufSize", 1048576, "")
}

func main() {
	var cfg config.Flags
	pflag.CommandLine.AddGoFlagSet(rootFlags)
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

}
