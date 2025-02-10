package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func (h *UsersHandler) returnInfoAboutUserById(c *fiber.Ctx, userID uuid.UUID) error {
	user, userFound, err := h.useCase.GetUserByID(userID)
	if !userFound {
		return NotFoundError(c, "couldn't find user with given user UUID")
	}
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UsersHandler) checkIfUserIdBelongsToCurrentUser(c *fiber.Ctx, userID uuid.UUID) error {
	if c.Locals("userID").(uuid.UUID) != userID {
		return BadRequest(c, "user ID is not owned by you")
	}
	return nil
}

func (h *UsersHandler) GetUserByIDHandler(c *fiber.Ctx) error {
	userID, err := ParseUUID(c, "id")
	if err != nil {
		return BadRequest(c, "invalid user UUID")
	}
	return h.returnInfoAboutUserById(c, userID)
}

func (h *UsersHandler) GetCurrentUserHandler(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	return h.returnInfoAboutUserById(c, userID)
}

func (h *UsersHandler) GetUserByLoginHandler(c *fiber.Ctx) error {
	login := c.Params("login")
	if login == "" {
		return BadRequest(c, "login is required")
	}
	user, userFound, err := h.useCase.GetUserByLogin(login)
	if !userFound {
		return NotFoundError(c, "couldn't find user with given login")
	}
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

	err = h.checkIfUserIdBelongsToCurrentUser(c, userID)
	if err != nil {
		return err
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

	err = h.checkIfUserIdBelongsToCurrentUser(c, userID)
	if err != nil {
		return err
	}

	if err := h.useCase.DeleteUser(userID); err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusNoContent).SendString("")
}
