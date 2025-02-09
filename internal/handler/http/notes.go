package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"notes_service/internal/handler/schemas"
	"notes_service/internal/models"
	"notes_service/internal/usecases"
	"strconv"
)

type NotesHandler struct {
	useCase usecases.NoteUseCase
}

func NewNotesHandler(useCase usecases.NoteUseCase) *NotesHandler {
	return &NotesHandler{useCase}
}

func (h *NotesHandler) ListNotesHandler(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequest(c, "invalid id")
	}
	result, err := h.useCase.GetNotesByUserID(uint(userID))
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *NotesHandler) GetNoteByIDHandler(c *fiber.Ctx) error {
	noteID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequest(c, "invalid note ID")
	}
	note, err := h.useCase.GetNoteByID(uint(noteID))
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(note)
}

func (h *NotesHandler) CreateNoteHandler(c *fiber.Ctx) error {
	var body schemas.NoteBodySchema
	if err := c.BodyParser(&body); err != nil {
		return BadRequest(c, "invalid request body")
	}
	note := models.Note{
		UserID:  body.UserID,
		Title:   body.Title,
		Content: body.Content,
	}
	createdNote, err := h.useCase.Create(note)
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(createdNote)
}

func (h *NotesHandler) UpdateNoteHandler(c *fiber.Ctx) error {
	noteID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequest(c, "invalid note ID")
	}

	var body schemas.NoteBodySchema
	if err := c.BodyParser(&body); err != nil {
		return BadRequest(c, "invalid request body")
	}
	updatedNote := models.Note{
		UserID:  body.UserID,
		Title:   body.Title,
		Content: body.Content,
	}

	result, err := h.useCase.Update(updatedNote, uint(noteID))
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *NotesHandler) DeleteNoteHandler(c *fiber.Ctx) error {
	noteID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequest(c, "invalid note ID")
	}
	if err := h.useCase.DeleteNote(uint(noteID)); err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusNoContent).SendString("")
}

func (h *NotesHandler) CountNotesByUserHandler(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	if userIDStr == "" {
		return BadRequest(c, "user ID is required")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return BadRequest(c, "invalid user ID")
	}

	count, err := h.useCase.CountNotesByUser(userID) // Преобразование uint в UUID
	if err != nil {
		return InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": count,
	})
}
