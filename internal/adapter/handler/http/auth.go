package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"notes_service/internal/adapter/config"
	"notes_service/internal/adapter/handler/schemas"
	users2 "notes_service/internal/usecases/users"
	"notes_service/internal/users"
	jwt2 "notes_service/pkg/auth/jwt"
	"time"
)

type JWTHandler struct {
	secretKey            string
	accessTokenLifetime  time.Duration
	refreshTokenLifetime time.Duration
	useCase              users2.UserAuthUseCase
}

func NewJWTHandler(useCase users2.UserAuthUseCase, jwtConfig config.JWT) *JWTHandler {

	return &JWTHandler{
		secretKey:            jwtConfig.SecretKey,
		useCase:              useCase,
		accessTokenLifetime:  time.Duration(jwtConfig.AccessTokenLifetimeMinutes) * time.Minute,
		refreshTokenLifetime: time.Duration(jwtConfig.RefreshTokenLifetimeMinutes) * time.Minute,
	}
}

func (h *JWTHandler) JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := jwt2.ValidateToken(tokenString, h.secretKey)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}

		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token type must be access"})
		}

		userIDString, ok := claims["user_id"].(string)
		userID, err := uuid.Parse(userIDString)
		if !ok || err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user ID in token"})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}

func (h *JWTHandler) SignUpHandler(c *fiber.Ctx) error {
	var body schemas.UserBodySchema
	if err := c.BodyParser(&body); err != nil {
		return BadRequest(c, "invalid request body")
	}

	user := users.User{
		Login:    body.Login,
		Password: body.Password,
	}
	createdUser, err := h.useCase.CreateUser(user)
	if err != nil {
		return InternalServerError(c, err)
	}

	accessToken, err := jwt2.GenerateToken(
		createdUser.ID,
		createdUser.Login,
		"access",
		h.accessTokenLifetime,
		h.secretKey,
	)
	if err != nil {
		return InternalServerError(c, err)
	}

	refreshToken, err := jwt2.GenerateToken(
		createdUser.ID,
		createdUser.Login,
		"refresh",
		h.refreshTokenLifetime,
		h.secretKey,
	)
	if err != nil {
		return InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
