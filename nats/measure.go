package nats

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/schollz/progressbar/v2"
)

// Measure starts the benchmark
func (me *Benchmark) Measure() (dur []time.Duration, err error) {
	ret := make([]time.Duration, me.NumItr)
	var bar *progressbar.ProgressBar
	if me.NumTrial < 2 {
		bar = progressbar.New(len(ret))
	}
	for i := range ret {
		var itrDur time.Duration
		var trBar *progressbar.ProgressBar
		if me.NumTrial > 1 {
			trBar = progressbar.NewOptions(
				me.NumTrial,
				progressbar.OptionSetDescription(fmt.Sprintf("Iteration: %v", i)),
			)
		}
		for j := 0; j < me.NumTrial; j++ {
			err = func() (trErr error) {
				ctx, stop := context.WithTimeout(me.Ctx, 1*time.Second)
				defer stop()
				var trDur time.Duration
				if trDur, trErr = me.measure(ctx); trErr == nil {
					itrDur += trDur
				}
				return
			}()
			if err != nil {
				return
			}
			if trBar != nil {
				trBar.Add(1)
			}
		}
		if trBar != nil {
			fmt.Println()
		}
		if bar != nil {
			bar.Add(1)
		}
		ret[i] = itrDur
	}
	if bar != nil {
		fmt.Println()
	}
	dur = ret
	return
}

func (me *Benchmark) measure(ctx context.Context) (
	dur time.Duration, err error,
) {
	data := make([]byte, me.BufSize)
	if _, err = rand.Read(data); err != nil {
		return
	}
	recvCh := make(chan *nats.Msg)
	errCh := make(chan error)
	defer close(errCh)
	defer close(recvCh)
	sub, err := me.Con.ChanSubscribe("test", recvCh)
	if err != nil {
		return
	}
	defer sub.Unsubscribe()
	startTime := time.Now()
	go func() {
		if err := me.Con.Publish("test", []byte(data)); err != nil {
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
