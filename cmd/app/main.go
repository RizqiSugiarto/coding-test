package main

import (
	"log"

	"github.com/RizqiSugiarto/coding-test/config"
	"github.com/RizqiSugiarto/coding-test/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
