package config

type target string

const (
	// Redis represents the target is redis.
	Redis target = "redis"
	// Nats represents the target is redis.
	Nats target = "nats"
)

// Flags represents a command line flags.
type Flags struct {
	Target       target
	NumTrial     int
	NumIteration int
	BufSize      int
}
