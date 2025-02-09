package postgres

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"notes_service/internal/adapter/config"
	"notes_service/internal/notes"
	"notes_service/internal/users"
)

type DBInstance struct {
	Db *gorm.DB
}

func New(ctx context.Context, config *config.DB) (*DBInstance, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s sslmode=disable password=%s",
		config.DbHost,
		config.DbUser,
		config.DbName,
		config.DbPassword,
	)

	db, connErr := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if connErr != nil {
		log.Fatal("Failed to connect to databse. \n", connErr)
	}

	log.Printf("connected to the database")
	db.Logger = logger.Default.LogMode(logger.Info)

	/*
		log.Printf("migrating")
		migrateErr := db.AutoMigrate(&models.Note{})
		if migrateErr != nil {
			log.Fatal("Failed to migrate. \n", migrateErr)
		}
	*/

	DB := DBInstance{db}
	return &DB, nil
}

func (db *DBInstance) AutoMigrate() error {
	return db.Db.AutoMigrate(
		&users.User{},
		&notes.Note{},
	)
}
