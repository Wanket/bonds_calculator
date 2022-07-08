package service

import (
	"bonds_calculator/internal/service"
	"github.com/benbjohnson/clock"
	asserts "github.com/stretchr/testify/assert"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	t.Parallel()

	assert := asserts.New(t)

	mockClock := clock.NewMock()

	timer := service.NewTimer(mockClock)

	generatedKey := rand.Int()

	doneChan := make(chan struct{})
	timer.SubscribeEvery(generatedKey, time.Minute*5, func(key int) {
		assert.Equal(generatedKey, key)

		doneChan <- struct{}{}
	})

	runtime.Gosched() // let the timer service to start wait for the first tick (otherwise mockClock ignore timer`s tick)

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
	t.Parallel()

	assert := asserts.New(t)

	mockClock := clock.NewMock()

	timer := service.NewTimer(mockClock)

	generatedKey := rand.Int()

	doneChan := make(chan struct{})
	timer.SubscribeEveryStartFrom(generatedKey, time.Minute*5, mockClock.Now().Add(time.Minute), func(key int) {
		assert.Equal(generatedKey, key)

		doneChan <- struct{}{}
	})

	durations := []time.Duration{
		time.Minute - 1,
		1,
		time.Minute * 5,
	}

	mockClock.Add(durations[0])

	select {
	case <-doneChan:
		assert.Fail("SubscribeEveryStartFrom callback should not be called")
	case <-time.After(time.Second):
	}

	for i := 1; i < 3; i++ {
		runtime.Gosched() // let the timer service to start wait for the first tick (otherwise mockClock ignore timer`s tick)

		mockClock.Add(durations[i])

		select {
		case <-doneChan:
		case <-time.After(time.Second):
			assert.Fail("SubscribeEveryStartFrom timeout")
		}
	}
}
