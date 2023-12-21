package main

import (
	"github.com/Entreeka/receiver/internal/app"
	"github.com/Entreeka/receiver/internal/config"
	"github.com/Entreeka/receiver/pkg/logger"
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
