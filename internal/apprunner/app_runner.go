package apprunner

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"notes_service/config"
	"notes_service/internal/handler/http"
	"notes_service/internal/models"
	"notes_service/internal/ports"
	use_cases "notes_service/internal/usecases"
	"notes_service/pkg/storage/postgres"
	"notes_service/pkg/storage/redis"
	"os/user"
	"time"
)

// AutoMigration performs a migration for all models from models module
func AutoMigration(db *postgres.DBInstance) error {
	return db.Db.AutoMigrate(&models.Note{}, &user.User{})
}

// RunApp sets up all connections and servers and returns an error when stops
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

	redisStorage := ports.NewRedisStorage(redisClient)

	notesRepo := ports.NewNotesRepoDB(db)
	notesUseCase := use_cases.NewNoteCRUDUseCase(notesRepo)
	notesHandler := http.NewNotesHandler(*notesUseCase)

	usersRepo := ports.NewUsersRepoDB(db)
	usersUseCase := use_cases.NewUserCRUDUseCase(usersRepo)

	usersHandler := http.NewUsersHandler(*usersUseCase)

	jwtHandler := http.NewJWTHandler(*usersUseCase, *mainConfig.JWT)

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
	protectedV1.Get("/notes/:id", notesHandler.GetNoteByIDHandler)
	protectedV1.Post("/notes", notesHandler.CreateNoteHandler)
	protectedV1.Put("/notes/:id", notesHandler.UpdateNoteHandler)
	protectedV1.Delete("/notes/:id", notesHandler.DeleteNoteHandler)
	protectedV1.Get("/notes/count-by-user/:id", notesHandler.CountNotesByUserHandler)

	protectedV1.Get("/users/:id", usersHandler.GetUserByIDHandler)
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
