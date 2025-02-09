package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	config2 "notes_service/internal/adapter/config"
	"notes_service/internal/adapter/handler/http"
	"notes_service/internal/adapter/storage/postgres"
	repository "notes_service/internal/adapter/storage/postgres/repositories"
	"notes_service/internal/usecases/notes"
)

func main() {
	mainConfig, err := config2.New()
	if err != nil {
		log.Fatal("couldn't load configs", err)
	}

	ctx := context.Background()

	db, err := postgres.New(ctx, mainConfig.DB)
	if err != nil {
		log.Fatal("couldn't connect to database", err)
	}

	err = db.AutoMigrate()
	if err != nil {
		log.Fatal("couldn't apply migrations", err)
	}

	notesRepo := repository.NewNotesRepo(db)
	notesUseCase := notes.NewNoteCRUDUseCase(notesRepo)
	notesHandler := http.NewNotesHandler(*notesUseCase)

	app := fiber.New()

	app.Get("/notes/by-user/:id", notesHandler.ListNotesHandler)
	app.Get("/notes/:id", notesHandler.GetNoteByIdHandler)
	app.Post("/notes", notesHandler.CreateNoteHandler)
	app.Put("/notes/:id", notesHandler.UpdateNoteHandler)
	app.Delete("/notes/:id", notesHandler.DeleteNoteHandler)
	app.Get("/notes/count-by-user/:id", notesHandler.CountNotesByUserHandler)

	err = app.Listen(":3000")
	log.Info("Listen on :3000")
	if err != nil {
		log.Fatal(err)
	}
}
