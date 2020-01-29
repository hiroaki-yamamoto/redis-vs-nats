package nats

import "github.com/nats-io/nats.go"

import "context"

// Benchmark represents a benchmark class of nats.
type Benchmark struct {
	Ctx      context.Context
	Con      *nats.Conn
	NumItr   int
	NumTrial int
	BufSize  int
}
