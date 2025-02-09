package config

import "os"

type (
	DB struct {
		DbHost     string
		DbUser     string
		DbName     string
		DbPassword string
	}
	Configs struct {
		DB *DB
	}
)

func New() (*Configs, error) {
	return &Configs{
		&DB{
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PASSWORD"),
		},
	}, nil
}
