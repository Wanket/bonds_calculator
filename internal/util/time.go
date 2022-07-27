package util

import (
	"fmt"
	"github.com/benbjohnson/clock"
	"time"
)

const (
	DayMultiplier = 24
)

const (
	Day = time.Hour * DayMultiplier
)

type ITimeHelper interface {
	GetMoexNow() time.Time

	GetMoexMidnight() time.Time
}

type TimeHelper struct {
	clock clock.Clock

	moexLocation *time.Location
	moexOffset   int
}

func NewTimeHelper(clock clock.Clock) *TimeHelper {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(fmt.Errorf("failed to load location: %w", err))
	}

	_, moexOffset := time.Now().In(location).Zone()

	return &TimeHelper{
		clock: clock,

		moexLocation: location,
		moexOffset:   moexOffset,
	}
}

func (helper *TimeHelper) GetMoexMidnight() time.Time {
	return helper.GetMoexNow().Truncate(Day)
}

func (helper *TimeHelper) GetMoexNow() time.Time {
	return helper.clock.Now().UTC().Add(time.Duration(helper.moexOffset) * time.Second)
}
