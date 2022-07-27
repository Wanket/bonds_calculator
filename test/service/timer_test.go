package service_test

import (
	"bonds_calculator/internal/service"
	"bonds_calculator/test"
	"fmt"
	"github.com/benbjohnson/clock"
	"runtime"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	assert, _ := test.PrepareTest(t)

	timer, mockClock := prepareTimer()
	defer timer.Close()

	doneChan := make(chan struct{})

	timer.SubscribeEvery(time.Minute*5, func() {
		doneChan <- struct{}{}
	})

	runtime.Gosched()

	durations := []time.Duration{
		time.Minute * 5,
		time.Minute*5 + 1,
		time.Minute*5 - 1,
		time.Minute*5 - 1,
	}

	for i := 0; i < 3; i++ {
		mockClock.Add(durations[i])

		select {
		case <-doneChan:
		case <-time.After(time.Second):
			assert.Fail("SubscribeEvery timeout")
		}
	}

	mockClock.Add(durations[3])

	select {
	case <-doneChan:
		assert.Fail("SubscribeEvery callback should not be called")
	case <-time.After(time.Second):
	}
}

func TestTimerStartFrom(t *testing.T) {
	assert, _ := test.PrepareTest(t)

	timer, mockClock := prepareTimer()
	defer timer.Close()

	doneChan := make(chan struct{})

	timer.SubscribeEveryStartFrom(time.Minute*5, mockClock.Now().Add(time.Minute), func() {
		doneChan <- struct{}{}
	})

	runtime.Gosched()

	durations := []time.Duration{
		time.Minute - time.Second,
		time.Second,
		time.Minute * 5,
	}

	mockClock.Add(durations[0])

	select {
	case <-doneChan:
		assert.Fail("SubscribeEveryStartFrom callback should not be called")
	case <-time.After(time.Second):
	}

	for duration := 1; duration < 3; duration++ {
		runtime.Gosched()

		time.Sleep(time.Millisecond) // wait for timer to register event

		mockClock.Add(durations[duration])

		select {
		case <-doneChan:
		case <-time.After(time.Second):
			assert.Fail("SubscribeEveryStartFrom timeout", fmt.Sprintf("duration id: %d", duration))
		}
	}
}

func prepareTimer() (*service.TimerService, *clock.Mock) {
	mockClock := clock.NewMock()

	timer := service.NewTimerService(mockClock)

	return timer, mockClock
}
