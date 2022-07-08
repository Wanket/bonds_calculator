//go:build wireinject

//go:generate wire

package internal

import (
	"bonds_calculator/internal/service"
	"github.com/benbjohnson/clock"
	"github.com/google/wire"
)

func NewTimerService() service.TimerService {
	wire.Build(service.NewTimer, clock.New)

	return service.TimerService{}
}
