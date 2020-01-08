package main

import (
	"bytes"
	"crypto/rand"
	"errors"
	"time"

	"github.com/go-redis/redis/v7"
)

func measureLatency(cli *redis.Client, len int64) (
	dur time.Duration, err error,
) {
	data := make([]byte, len)
	if _, err = rand.Read(data); err != nil {
		return
	}
	endTimeCh := make(chan time.Time)
	errCh := make(chan error)
	go func() {
		sub := cli.Subscribe("test")
		defer sub.Close()
		defer close(errCh)
		msg := <-sub.Channel()
		end := time.Now()
		if bytes.Compare([]byte(msg.Payload), data) != 0 {
			errCh <- errors.New(
				"The message is differ from the expected message",
			)
		}
		endTimeCh <- end
	}()
	startTime := time.Now()
	cli.Publish("test", string(data))
	select {
	case endTime := <-endTimeCh:
		dur = endTime.Sub(startTime)
		return
	case err = <-errCh:
		return
	}
}
