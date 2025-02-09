package http

import (
	"github.com/gofiber/fiber/v2"
	"notes_service/internal/handler/schemas"
	"notes_service/internal/models"
	use_cases "notes_service/internal/usecases"
)

type UsersHandler struct {
	useCase use_cases.UserUseCase
}

func NewUsersHandler(useCase use_cases.UserUseCase) *UsersHandler {
	return &UsersHandler{useCase}
}

func (h *UsersHandler) GetUserByIDHandler(c *fiber.Ctx) error {
	userID, err := ParseUUID(c, "id")
	if err != nil {
		return BadRequest(c, "invalid user UUID")
	}
	user, err := h.useCase.GetUserByID(userID)
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UsersHandler) GetUserByLoginHandler(c *fiber.Ctx) error {
	login := c.Params("login")
	if login == "" {
		return BadRequest(c, "login is required")
	}
	user, err := h.useCase.GetUserByLogin(login)
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UsersHandler) UpdateUserHandler(c *fiber.Ctx) error {
	userID, err := ParseUUID(c, "id")
	if err != nil {
		return BadRequest(c, "invalid user UUID")
	}

	var body schemas.UserBodySchema
	if err := c.BodyParser(&body); err != nil {
		return BadRequest(c, "invalid request body")
	}

	user := models.User{
		Login:    body.Login,
		Password: body.Password,
	}
	updatedUser, err := h.useCase.UpdateUser(user, userID)
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(updatedUser)
}

func (h *UsersHandler) DeleteUserHandler(c *fiber.Ctx) error {
	userID, err := ParseUUID(c, "id")
	if err != nil {
		return BadRequest(c, "invalid user UUID")
	}

	if err := h.useCase.DeleteUser(userID); err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusNoContent).SendString("")
}
