package main

import (
	"errors"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v7"
)

func measureLatency(cli *redis.Client, len int) (
	dur time.Duration, err error,
) {
	var data []byte
	if _, err = rand.Read(data); err != nil {
		return
	}
	txt := string(data)
	errCh := make(chan error)
	sub := cli.Subscribe("test")
	defer sub.Close()
	startTime := time.Now()
	go func() {
		defer close(errCh)
		if _, err = cli.Publish("test", txt).Result(); err != nil {
			errCh <- err
			return
		}
	}()
	select {
	case msg := <-sub.Channel():
		end := time.Now()
		if msg.Payload != txt {
			err = errors.New("The message is differ from the expected message")
			return
		}
		dur = end.Sub(startTime)
		return
	case err = <-errCh:
		return
	}
}
