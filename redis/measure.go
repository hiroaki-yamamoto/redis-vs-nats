package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/schollz/progressbar/v2"
)

// Measure runs benchmark.
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
				var trDur time.Duration
				if trDur, trErr = me.measure(); trErr == nil {
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

func (me *Benchmark) measure() (dur time.Duration, err error) {
	ctx, stop := context.WithTimeout(me.Ctx, 1*time.Second)
	defer stop()
	data := make([]byte, me.BufSize)
	if _, err = rand.Read(data); err != nil {
		return
	}
	txt := string(data)
	errCh := make(chan error)
	defer close(errCh)
	sub := me.Con.Subscribe("test")
	defer sub.Close()
	startTime := time.Now()
	go func() {
		if _, err = me.Con.Publish("test", txt).Result(); err != nil {
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
	case <-ctx.Done():
		return
	case err = <-errCh:
		log.Println(err)
		return
	}
}
