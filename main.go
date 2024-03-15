package main

import (
	"github.com/Dubjay18/gobank_auth/app"
	"github.com/Dubjay18/gobank_auth/logger"
)

func main() {
	app.GetEnvVar()
	logger.Info("Starting the application...")
	app.Start()
}
