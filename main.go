package main

import (
	"github.com/divyag/services/app"
	"github.com/divyag/services/logger"
)

func main() {
	logger.Info("Starting Application....")
	app.Start()
}
