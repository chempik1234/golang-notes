package ports

import (
	"github.com/google/uuid"
	"notes_service/internal/models"
)

// UsersRepo Port for users
type UsersRepo interface {
	GetUserByID(userID uuid.UUID) (models.User, error)
	GetUserByLogin(login string) (models.User, error)
	GetUserByLoginAndPassword(login string, password string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User, userID uuid.UUID) (models.User, error)
	DeleteUser(userID uuid.UUID) error
}

// NotesRepo Port for notes
type NotesRepo interface {
	GetNotesByUserID(userID uint) ([]models.Note, error)
	GetNoteByID(noteID uint) (models.Note, error)
	CreateNote(note models.Note) (models.Note, error)
	UpdateNote(note models.Note, id uint) (models.Note, error)
	DeleteNote(noteID uint) error
	CountNotesByUser(userID uuid.UUID) (int64, error)
}
