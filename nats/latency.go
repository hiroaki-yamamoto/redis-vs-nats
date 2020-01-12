package main

import (
	"bytes"
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/nats-io/nats.go"
)

func measureLatency(ctx context.Context, con *nats.Conn, sz int) (
	dur time.Duration, err error,
) {
	data := make([]byte, sz)
	if _, err = rand.Read(data); err != nil {
		return
	}
	recvCh := make(chan *nats.Msg)
	errCh := make(chan error)
	defer close(errCh)
	defer close(recvCh)
	sub, err := con.ChanSubscribe("test", recvCh)
	if err != nil {
		return
	}
	defer sub.Unsubscribe()
	startTime := time.Now()
	go func() {
		if err := con.Publish("test", []byte(data)); err != nil {
			errCh <- err
			return
		}
	}()

	select {
	case msg := <-recvCh:
		end := time.Now()
		if bytes.Compare(msg.Data, data) != 0 {
			errCh <- errors.New(
				"The received data was differ from the expected",
			)
			return
		}
		dur = end.Sub(startTime)
		return
	case err = <-errCh:
		return
	case <-ctx.Done():
		end := time.Now()
		dur = end.Sub(startTime)
		return
	}
}
