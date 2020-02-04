package config

import "github.com/spf13/viper"

type target string

const (
	// Redis represents the target is redis.
	Redis target = "redis"
	// Nats represents the target is redis.
	Nats target = "nats"
)

// Config represents a config.
type Config struct {
	Target       target
	Addr         string
	NumTrial     int
	NumIteration int
	BufSize      int
}

// New makes a new config.
func New(cfgPath string) (cfg *Config, err error) {
	viper.SetConfigFile(cfgPath)
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	cfg = &Config{}
	if err = viper.Unmarshal(cfg); err != nil {
		cfg = nil
		return
	}
	return
}
