package main

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v7"
)

func measureLatency(cli *redis.Client, sz int) (
	dur time.Duration, err error,
) {
	data := make([]byte, sz)
	if _, err = rand.Read(data); err != nil {
		return
	}
	txt := string(data)
	errCh := make(chan error)
	defer close(errCh)
	sub := cli.Subscribe("test")
	defer sub.Close()
	startTime := time.Now()
	go func() {
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
		log.Println(err)
		return
	}
}
