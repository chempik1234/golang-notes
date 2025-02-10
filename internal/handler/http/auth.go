package http

import (
	"errors"
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
	GetUserByID(userID uuid.UUID) (models.User, bool, error)
	CreateUser(user models.User) (models.User, error)
	GetUserByLogin(login string) (models.User, bool, error)
	GetUserByLoginAndPassword(login string, password string) (models.User, bool, error)
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
		userID, _, err := h.runChecksForTokenString(tokenString, "access")
		if err != nil {
			return NotAuthenticatedError(c, err)
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}

func (h *JWTHandler) runChecksForTokenString(tokenString string, requiredTokenType string) (uuid.UUID, string, error) {
	token, err := jwtutils.ValidateToken(tokenString, h.secretKey)
	if err != nil || !token.Valid {
		return uuid.UUID{}, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, "", errors.New("invalid token claims")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != requiredTokenType {
		return uuid.UUID{}, "", errors.New("token type must be " + requiredTokenType)
	}

	userIDString, ok := claims["sub"].(string)
	userID, err := uuid.Parse(userIDString)
	if !ok || err != nil {
		return uuid.UUID{}, "", errors.New("invalid user ID in token")
	}

	userLogin := claims["username"].(string)

	// user, userFound, err := h.useCase.GetUserByID(userID)
	// if !userFound || err != nil {
	// 	return models.User{}, errors.New("can't find user with given user ID")
	// }

	return userID, userLogin, nil
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

func (h *JWTHandler) createTokensForUser(userID uuid.UUID, userLogin string) (fiber.Map, error) {
	accessToken, err := jwtutils.GenerateToken(
		userID,
		userLogin,
		"access",
		h.accessTokenLifetime,
		h.secretKey,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwtutils.GenerateToken(
		userID,
		userLogin,
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

// SignUpHandler HTTP handler for creating a new user and retrieving a new token
func (h *JWTHandler) SignUpHandler(c *fiber.Ctx) error {
	user, err := h.parseUser(c)
	if err != nil {
		return BadRequest(c, "invalid user data")
	}

	_, found, err := h.useCase.GetUserByLogin(user.Login)
	if err != nil {
		return InternalServerError(c, err)
	}
	if found {
		return BadRequest(c, "username already exists")
	}

	createdUser, err := h.useCase.CreateUser(user)
	if err != nil {
		return InternalServerError(c, err)
	}

	returnData, err := h.createTokensForUser(createdUser.ID, createdUser.Login)
	if err != nil {
		return InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(returnData)
}

// SignInHandler HTTP Handler for retrieving a new token by login and password
func (h *JWTHandler) SignInHandler(c *fiber.Ctx) error {
	user, err := h.parseUser(c)
	if err != nil {
		return BadRequest(c, "invalid user data")
	}

	foundUser, userFound, err := h.useCase.GetUserByLoginAndPassword(user.Login, user.Password)
	if err != nil {
		return InternalServerError(c, err)
	}
	if !userFound {
		return NotAuthenticatedError(c, errors.New("user not found"))
	}

	returnData, err := h.createTokensForUser(foundUser.ID, foundUser.Login)
	if err != nil {
		return InternalServerError(c, err)
	}

	return c.JSON(returnData)
}

// RefreshHandler HTTP handler for refreshing JWT's
func (h *JWTHandler) RefreshHandler(c *fiber.Ctx) error {
	var body schemas.RefreshTokenSchema
	if err := c.BodyParser(&body); err != nil {
		return BadRequest(c, "invalid refresh token data")
	}

	tokenString := body.RefreshToken

	userID, userLogin, err := h.runChecksForTokenString(tokenString, "refresh")
	if err != nil {
		return NotAuthenticatedError(c, err)
	}

	returnData, err := h.createTokensForUser(userID, userLogin)
	if err != nil {
		return InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(returnData)
}
