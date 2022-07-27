package service

import (
	"context"
	"github.com/benbjohnson/clock"
	"time"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mock/timer_gen.go . ITimerService
type ITimerService interface {
	SubscribeEvery(duration time.Duration, callback func())
	SubscribeEveryStartFrom(duration time.Duration, startTime time.Time, callback func())

	Close() error
}

type TimerService struct {
	context context.Context //nolint:containedctx
	cancel  context.CancelFunc

	clock clock.Clock
}

func NewTimerService(clock clock.Clock) *TimerService {
	ctx, cancel := context.WithCancel(context.Background())

	return &TimerService{
		context: ctx,
		cancel:  cancel,
		clock:   clock,
	}
}

func (t *TimerService) SubscribeEvery(duration time.Duration, callback func()) {
	go func() {
		ticker := t.clock.Ticker(duration)

		for {
			select {
			case <-ticker.C:
				callback()
			case <-t.context.Done():
				ticker.Stop()

				return
			}
		}
	}()
}

func (t *TimerService) SubscribeEveryStartFrom(duration time.Duration, startTime time.Time, callback func()) {
	go func() {
		timer := t.clock.Timer(startTime.Sub(t.clock.Now()))

		select {
		case <-timer.C:
			callback()

			t.SubscribeEvery(duration, callback)
		case <-t.context.Done():
			timer.Stop()
		}
	}()
}

func (t *TimerService) Close() error {
	t.cancel()

	return nil
}
