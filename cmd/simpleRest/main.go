package main

import (
	"github.com/gofiber/fiber/v2/log"
	"notes_service/config"
	"notes_service/internal/apprunner"
)

func main() {
	mainConfig, err := config.New()
	if err != nil {
		log.Fatal("couldn't load configs", err)
	}

	err = apprunner.RunApp(mainConfig)
	if err != nil {
		log.Fatal(err)
	}
}
