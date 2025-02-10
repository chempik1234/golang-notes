package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": message,
	})
}

func InternalServerError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func NotAuthenticatedError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func NotFoundError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": message,
	})
}

func ForbiddenError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{})
}

func ParseUUID(c *fiber.Ctx, fieldName string) (uuid.UUID, error) {
	uuidField, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuidField, nil
}
