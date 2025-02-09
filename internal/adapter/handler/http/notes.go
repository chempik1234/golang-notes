package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"notes_service/internal/notes"
	notes2 "notes_service/internal/usecases/notes"
	"strconv"
)

type NotesHandler struct {
	useCase notes2.NoteCRUDUseCase
}

func NewNotesHandler(useCase notes2.NoteCRUDUseCase) *NotesHandler {
	return &NotesHandler{useCase}
}

func (h *NotesHandler) ListNotesHandler(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequest(c, "invalid id")
	}
	result, err := h.useCase.GetNotesByUserId(uint(userId))
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *NotesHandler) GetNoteByIdHandler(c *fiber.Ctx) error {
	noteId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequest(c, "invalid note ID")
	}
	note, err := h.useCase.GetNoteById(uint(noteId))
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(note)
}

func (h *NotesHandler) CreateNoteHandler(c *fiber.Ctx) error {
	var note notes.Note
	if err := c.BodyParser(&note); err != nil {
		return BadRequest(c, "invalid request body")
	}
	createdNote, err := h.useCase.Create(note)
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(createdNote)
}

func (h *NotesHandler) UpdateNoteHandler(c *fiber.Ctx) error {
	noteId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequest(c, "invalid note ID")
	}

	var updatedNote notes.Note
	if err := c.BodyParser(&updatedNote); err != nil {
		return BadRequest(c, "invalid request body")
	}

	result, err := h.useCase.Update(updatedNote, uint(noteId))
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *NotesHandler) DeleteNoteHandler(c *fiber.Ctx) error {
	noteId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequest(c, "invalid note ID")
	}
	if err := h.useCase.DeleteNote(uint(noteId)); err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusNoContent).SendString("")
}

func (h *NotesHandler) CountNotesByUserHandler(c *fiber.Ctx) error {
	userIdStr := c.Params("id")
	if userIdStr == "" {
		return BadRequest(c, "user ID is required")
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return BadRequest(c, "invalid user ID")
	}

	count, err := h.useCase.CountNotesByUser(userId) // Преобразование uint в UUID
	if err != nil {
		return InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": count,
	})
}
