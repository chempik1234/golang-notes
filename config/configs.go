package config

import (
	"os"
	"strconv"
)

type (
	DB struct {
		DbHost     string
		DbUser     string
		DbName     string
		DbPassword string
	}
	HTTP struct {
		Port string
	}
	JWT struct {
		SecretKey                   string
		AccessTokenLifetimeMinutes  int
		RefreshTokenLifetimeMinutes int
	}
	Limiter struct {
		MaxConnections    int
		ExpirationSeconds int
	}
	Redis struct {
		URL string
	}
	Configs struct {
		DB      *DB
		JWT     *JWT
		Limiter *Limiter
		HTTP    *HTTP
		Redis   *Redis
	}
)

func New() (*Configs, error) {
	accessTokenLifetimeMinutes, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_LIFETIME_MINUTES"))
	if err != nil {
		return nil, err
	}
	refreshTokenLifetimeMinutes, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_LIFETIME_MINUTES"))
	if err != nil {
		return nil, err
	}
	maxConnections, err := strconv.Atoi(os.Getenv("MAX_CONNECTIONS"))
	if err != nil {
		return nil, err
	}
	rateLimiterExpirationSeconds, err := strconv.Atoi(os.Getenv("RATE_LIMITER_EXPIRATION_SECONDS"))
	if err != nil {
		return nil, err
	}

	return &Configs{
		&DB{
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PASSWORD"),
		},
		&JWT{
			os.Getenv("JWT_SECRET_KEY"),
			accessTokenLifetimeMinutes,
			refreshTokenLifetimeMinutes,
		},
		&Limiter{
			MaxConnections:    maxConnections,
			ExpirationSeconds: rateLimiterExpirationSeconds,
		},
		&HTTP{
			os.Getenv("HTTP_PORT"),
		},
		&Redis{
			os.Getenv("REDIS_URL"),
		},
	}, nil
}
