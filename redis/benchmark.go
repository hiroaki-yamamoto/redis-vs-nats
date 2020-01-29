package redis

import (
	"context"

	"github.com/go-redis/redis/v7"
)

// Benchmark represents a benchmark class of redis.
type Benchmark struct {
	Ctx      context.Context
	Con      *redis.Client
	NumItr   int
	NumTrial int
	BufSize  int
}
