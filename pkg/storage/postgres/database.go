package postgres

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"notes_service/config"
)

// DBInstance is a structure for storing a database definition
type DBInstance struct {
	Db *gorm.DB
}

// NewDBInstance try to connect to the database and create a new DBInstance to use in repos
func NewDBInstance(ctx context.Context, config *config.DB) (*DBInstance, error) {
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

	log.Info("connected to the postgresql database")
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
