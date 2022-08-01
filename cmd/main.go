package main

import (
	"bonds_calculator/internal"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := internal.CreateApp().Run(); err != nil {
		log.WithError(err).Fatal("Application: got fatal error")
	}
}
