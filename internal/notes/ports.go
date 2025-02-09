package notes

import (
	"errors"
	"github.com/google/uuid"
)

type NotesRepo interface {
	GetNotesByUserId(userId uint) ([]Note, error)
	GetNoteById(noteId uint) (Note, error)
	CreateNote(note Note) (Note, error)
	UpdateNote(note Note, id uint) (Note, error)
	DeleteNote(noteId uint) error
	CountNotesByUser(userId uuid.UUID) (int64, error)
}

var ErrNoteNotFound = errors.New("note not found")
