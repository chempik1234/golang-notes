package app_runner

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"notes_service/internal/adapter/config"
	"notes_service/internal/adapter/handler/http"
	repository "notes_service/internal/adapter/storage/postgres/repositories"
	redis2 "notes_service/internal/adapter/storage/redis"
	notes2 "notes_service/internal/notes"
	"notes_service/internal/usecases/notes"
	"notes_service/internal/usecases/users"
	"notes_service/pkg/storage/postgres"
	"notes_service/pkg/storage/redis"
	"os/user"
	"time"
)

func AutoMigration(db *postgres.DBInstance) error {
	return db.Db.AutoMigrate(&notes2.Note{}, &user.User{})
}

func RunApp(mainConfig *config.Configs) error {

	ctx := context.Background()

	db, err := postgres.NewDBInstance(ctx, mainConfig.DB)
	if err != nil {
		log.Fatal("couldn't connect to postgresql database", err)
	}

	err = AutoMigration(db)
	if err != nil {
		log.Fatal("couldn't apply migrations", err)
	}

	redisClient, err := redis.NewRedisClient(mainConfig.Redis.URL)
	if err != nil {
		log.Fatal("couldn't connect to redis database", err)
	}

	redisStorage := redis2.NewRedisStorage(redisClient)

	notesRepo := repository.NewNotesRepo(db)
	notesUseCase := notes.NewNoteCRUDUseCase(notesRepo)
	notesHandler := http.NewNotesHandler(*notesUseCase)

	usersRepo := repository.NewUsersRepo(db)
	usersUseCase := users.NewUserCRUDUseCase(usersRepo)
	usersHandler := http.NewUsersHandler(*usersUseCase)

	authUseCase := users.NewUserAuthUseCase(usersRepo)
	jwtHandler := http.NewJWTHandler(*authUseCase, *mainConfig.JWT)

	app := fiber.New()

	rateLimiter := limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        mainConfig.Limiter.MaxConnections,
		Expiration: time.Duration(mainConfig.Limiter.ExpirationSeconds) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		Storage: redisStorage,
	})

	app.Use(rateLimiter)

	apiV1 := app.Group("/api/v1")

	apiV1.Post("/users", jwtHandler.SignUpHandler)

	protectedV1 := apiV1.Use("/protected", jwtHandler.JWTMiddleware())

	protectedV1.Get("/notes/by-user/:id", notesHandler.ListNotesHandler)
	protectedV1.Get("/notes/:id", notesHandler.GetNoteByIdHandler)
	protectedV1.Post("/notes", notesHandler.CreateNoteHandler)
	protectedV1.Put("/notes/:id", notesHandler.UpdateNoteHandler)
	protectedV1.Delete("/notes/:id", notesHandler.DeleteNoteHandler)
	protectedV1.Get("/notes/count-by-user/:id", notesHandler.CountNotesByUserHandler)

	protectedV1.Get("/users/:id", usersHandler.GetUserByIdHandler)
	protectedV1.Get("/users/by-login/:login", usersHandler.GetUserByLoginHandler)
	protectedV1.Put("/users/:id", usersHandler.UpdateUserHandler)
	protectedV1.Delete("/users/:id", usersHandler.DeleteUserHandler)

	log.Info("Listening on port " + mainConfig.HTTP.Port)
	log.Info("Redis on " + mainConfig.Redis.URL)
	err = app.Listen(":" + mainConfig.HTTP.Port)
	log.Info("finished!")
	if err != nil {
		return err
	}
	return nil
}
