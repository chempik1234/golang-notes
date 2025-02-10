package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"notes_service/internal/handler/schemas"
	"notes_service/internal/models"
	"notes_service/internal/usecases"
	"strconv"
)

// NotesHandler is a struct of HTTP handlers that relate to notes. Uses usecases.NoteUseCase.
type NotesHandler struct {
	useCase usecases.NoteUseCase
}

// NewNotesHandler creates and returns a new NotesHandler with given useCase
func NewNotesHandler(useCase usecases.NoteUseCase) *NotesHandler {
	return &NotesHandler{useCase}
}

func (h *NotesHandler) checkNoteBelongsToUser(c *fiber.Ctx, noteID uint, userID uuid.UUID) error {
	existingNote, foundNote, err := h.useCase.GetNoteByID(noteID)
	if err != nil {
		return InternalServerError(c, err)
	}
	if !foundNote {
		return NotFoundError(c, "couldn't find note with given id")
	}
	if existingNote.UserID != userID {
		return ForbiddenError(c)
	}
	return nil
}

// ListNotesHandler HTTP handler to list all notes by someone's UUID
func (h *NotesHandler) ListNotesHandler(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return BadRequest(c, "invalid id")
	}
	result, err := h.useCase.GetNotesByUserID(userID)
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

// GetNoteByIDHandler HTTP handler to retrieve note by id.
//
// id param: id of the returned note
//
// Returns 404 if the note couldn't be found.
func (h *NotesHandler) GetNoteByIDHandler(c *fiber.Ctx) error {
	noteID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequest(c, "invalid note ID")
	}
	note, noteFound, err := h.useCase.GetNoteByID(uint(noteID))
	if !noteFound {
		return NotFoundError(c, "couldn't find note with given id")
	}
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(note)
}

// CreateNoteHandler HTTP handler to create notes. Assigns the new record to current user.
func (h *NotesHandler) CreateNoteHandler(c *fiber.Ctx) error {
	var body schemas.NoteBodySchema
	if err := c.BodyParser(&body); err != nil {
		return BadRequest(c, "invalid request body")
	}
	note := models.Note{
		UserID:  c.Locals("userID").(uuid.UUID),
		Title:   body.Title,
		Content: body.Content,
	}
	createdNote, err := h.useCase.Create(note)
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(createdNote)
}

// UpdateNoteHandler HTTP handler to update a note by id.
//
// id param: id of the note being affected.
//
// checks if the note belongs to the authenticated user
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
		UserID:  c.Locals("userID").(uuid.UUID),
		Title:   body.Title,
		Content: body.Content,
	}

	noteIDUint := uint(noteID)

	err = h.checkNoteBelongsToUser(c, noteIDUint, updatedNote.UserID)
	if err != nil {
		return err
	}

	result, err := h.useCase.Update(updatedNote, noteIDUint)
	if err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

// DeleteNoteHandler HTTP handler to delete a note by id.
//
// id param: id of the note being deleted.
//
// checks if the note belongs to the authenticated user
func (h *NotesHandler) DeleteNoteHandler(c *fiber.Ctx) error {
	noteID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequest(c, "invalid note ID")
	}

	noteIDUint := uint(noteID)

	err = h.checkNoteBelongsToUser(c, noteIDUint, c.Locals("userID").(uuid.UUID))
	if err != nil {
		return err
	}

	if err := h.useCase.DeleteNote(uint(noteID)); err != nil {
		return InternalServerError(c, err)
	}
	return c.Status(fiber.StatusNoContent).SendString("")
}

// CountNotesByUserHandler HTTP handler that returns the amount (count) of someone's notes.
//
// id param: string uuid of the user whose notes the handler is going to count
//
// Returns 0 if user couldn't be found
func (h *NotesHandler) CountNotesByUserHandler(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	if userIDStr == "" {
		return BadRequest(c, "user ID is required")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return BadRequest(c, "invalid user ID")
	}

	count, err := h.useCase.CountNotesByUser(userID)
	if err != nil {
		return InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": count,
	})
}
