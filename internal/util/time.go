package util

import (
	"github.com/benbjohnson/clock"
	"time"
)

var moexLocation, _ = time.LoadLocation("Europe/Moscow")
var _, moexOffset = time.Now().In(moexLocation).Zone()

func GetMoexMidnight(clock clock.Clock) time.Time {
	return GetMoexNow(clock).Truncate(time.Hour * 24)
}

func GetMoexNow(clock clock.Clock) time.Time {
	return clock.Now().UTC().Add(time.Duration(moexOffset) * time.Second)
}
