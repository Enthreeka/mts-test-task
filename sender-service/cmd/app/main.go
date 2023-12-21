package main

import (
	"github.com/Entreeka/sender/internal/app"
	"github.com/Entreeka/sender/internal/config"
	"github.com/Entreeka/sender/pkg/logger"
)

func main() {
	log := logger.New()

	cfg, err := config.New()
	if err != nil {
		log.Fatal("failed load config: %v", err)
	}

	if err := app.Run(cfg, log); err != nil {
		log.Fatal("failed to run server: %v", err)
	}
}
