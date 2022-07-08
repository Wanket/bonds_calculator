package service

import (
	"context"
	"github.com/benbjohnson/clock"
	"time"
)

type TimerService struct {
	context context.Context
	cancel  context.CancelFunc

	clock clock.Clock
}

func NewTimer(clock clock.Clock) TimerService {
	ctx, cancel := context.WithCancel(context.Background())

	return TimerService{
		context: ctx,
		cancel:  cancel,
		clock:   clock,
	}
}

func (t *TimerService) SubscribeEvery(key int, duration time.Duration, callback func(int)) {
	go func() {
		ticker := t.clock.Ticker(duration)

		for {
			select {
			case <-ticker.C:
				callback(key)
			case <-t.context.Done():
				ticker.Stop()

				break
			}
		}
	}()
}

func (t *TimerService) SubscribeEveryStartFrom(key int, duration time.Duration, startTime time.Time, callback func(int)) {
	go func() {
		timer := t.clock.Timer(startTime.Sub(t.clock.Now()))

		select {
		case <-timer.C:
			callback(key)

			t.SubscribeEvery(key, duration, callback)
		case <-t.context.Done():
			timer.Stop()
		}
	}()
}

func (t *TimerService) Close() error {
	t.cancel()

	return nil
}
