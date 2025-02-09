package main

import (
	"github.com/gofiber/fiber/v2/log"
	"notes_service/internal/adapter/config"
	"notes_service/internal/app_runner"
)

func main() {
	mainConfig, err := config.New()
	if err != nil {
		log.Fatal("couldn't load configs", err)
	}

	err = app_runner.RunApp(mainConfig)
	if err != nil {
		log.Fatal(err)
	}
}
