package main

import (
	"bytes"
	"errors"
	"math/rand"
	"time"

	"github.com/nats-io/nats.go"
)

func measureLatency(con *nats.Conn, len int) (
	dur time.Duration, err error,
) {
	data := make([]byte, len)
	if _, err = rand.Read(data); err != nil {
		return
	}
	endTimeCh := make(chan time.Time)
	errCh := make(chan error)
	defer close(errCh)
	defer close(endTimeCh)
	sub, err := con.Subscribe("test", func(msg *nats.Msg) {
		t := time.Now()
		if bytes.Compare(msg.Data, data) != 0 {
			errCh <- errors.New(
				"The received data was differ from the expected",
			)
			return
		}
		endTimeCh <- t
	})
	if err != nil {
		errCh <- err
	}
	defer sub.Unsubscribe()
	startTime := time.Now()
	if err = con.Publish("test", data); err != nil {
		return
	}

	select {
	case endTime := <-endTimeCh:
		dur = endTime.Sub(startTime)
		return
	case err = <-errCh:
		return
	}
}
