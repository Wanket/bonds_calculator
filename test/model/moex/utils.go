package moex

import (
	"net/http"
	"time"
)

var (
	client = http.Client{}
)

func parseDate(str string) time.Time {
	res, _ := time.Parse("2006-01-02", str)

	return res
}
