package http

import (
	"github.com/gofiber/fiber/v2"
	"notes_service/internal/adapter/handler/schemas"
	users2 "notes_service/internal/usecases/users"
	"notes_service/internal/users"
)

type UsersHandler struct {
	useCase users2.UserCRUDUseCase
}

func NewUsersHandler(useCase users2.UserCRUDUseCase) *UsersHandler {
	return &UsersHandler{useCase}
}

func (h *UsersHandler) GetUserByIdHandler(c *fiber.Ctx) error {
	userId, err := ParseUUID(c, "id")
	if err != nil {
		return BadRequest(c, "invalid user UUID")
	}
	user, err := h.useCase.GetUserById(userId)
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
	userId, err := ParseUUID(c, "id")
	if err != nil {
		return BadRequest(c, "invalid user UUID")
	}

	var body schemas.UserBodySchema
	if err := c.BodyParser(&body); err != nil {
		return BadRequest(c, "invalid request body")
	}

	user := users.User{
		Login:    body.Login,
		Password: body.Password,
	}
	updatedUser, err := h.useCase.UpdateUser(user, userId)
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(updatedUser)
}

func (h *UsersHandler) DeleteUserHandler(c *fiber.Ctx) error {
	userId, err := ParseUUID(c, "id")
	if err != nil {
		return BadRequest(c, "invalid user UUID")
	}

	if err := h.useCase.DeleteUser(userId); err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusNoContent).SendString("")
}
