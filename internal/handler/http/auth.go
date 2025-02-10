package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"notes_service/config"
	"notes_service/internal/handler/schemas"
	"notes_service/internal/models"
	"notes_service/pkg/auth/jwtutils"
	"time"
)

// AuthUseCase Use case interface required for JWTHandler
type AuthUseCase interface {
	GetUserByID(userID uuid.UUID) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	GetUserByLoginAndPassword(login string, password string) (models.User, error)
}

// JWTHandler is the auth header for fiber application.
// It has some values for creating JWT's and an auth use case
type JWTHandler struct {
	secretKey            string
	accessTokenLifetime  time.Duration
	refreshTokenLifetime time.Duration
	useCase              AuthUseCase
}

// NewJWTHandler creates and returns a new instance of NewJWTHandler.
// It accepts the use case and a config.JWT to extract values from.
func NewJWTHandler(useCase AuthUseCase, jwtConfig config.JWT) *JWTHandler {

	return &JWTHandler{
		secretKey:            jwtConfig.SecretKey,
		useCase:              useCase,
		accessTokenLifetime:  time.Duration(jwtConfig.AccessTokenLifetimeMinutes) * time.Minute,
		refreshTokenLifetime: time.Duration(jwtConfig.RefreshTokenLifetimeMinutes) * time.Minute,
	}
}

// JWTMiddleware is an auth middleware that checks the JWT with the validator
// and use case that tries to find the user in the actual database.
func (h *JWTHandler) JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := jwtutils.ValidateToken(tokenString, h.secretKey)
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

		userIDString, ok := claims["sub"].(string)
		userID, err := uuid.Parse(userIDString)
		if !ok || err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user ID in token"})
		}

		_, err = h.useCase.GetUserByID(userID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "can't find user with given user ID"})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}

func (h *JWTHandler) parseUser(c *fiber.Ctx) (models.User, error) {
	var body schemas.UserBodySchema
	if err := c.BodyParser(&body); err != nil {
		return models.User{}, err
	}

	user := models.User{
		Login:    body.Login,
		Password: body.Password,
	}
	return user, nil
}

func (h *JWTHandler) createTokensForUser(user models.User) (fiber.Map, error) {
	accessToken, err := jwtutils.GenerateToken(
		user.ID,
		user.Login,
		"access",
		h.accessTokenLifetime,
		h.secretKey,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwtutils.GenerateToken(
		user.ID,
		user.Login,
		"refresh",
		h.refreshTokenLifetime,
		h.secretKey,
	)
	if err != nil {
		return nil, err
	}

	return fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}, nil
}

func (h *JWTHandler) SignUpHandler(c *fiber.Ctx) error {
	user, err := h.parseUser(c)
	if err != nil {
		return BadRequest(c, "invalid user data")
	}

	createdUser, err := h.useCase.CreateUser(user)
	if err != nil {
		return InternalServerError(c, err)
	}

	returnData, err := h.createTokensForUser(createdUser)
	if err != nil {
		return InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(returnData)
}

func (h *JWTHandler) Login(c *fiber.Ctx) error {
	user, err := h.parseUser(c)
	if err != nil {
		return BadRequest(c, "invalid user data")
	}

	foundUser, err := h.useCase.GetUserByLoginAndPassword(user.Login, user.Password)
	if err != nil {
		return NotAuthenticatedError(c)
	}

	returnData, err := h.createTokensForUser(foundUser)
	if err != nil {
		return InternalServerError(c, err)
	}

	return c.JSON(returnData)
}
