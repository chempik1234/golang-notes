package ports

import (
	"github.com/google/uuid"
	"notes_service/internal/models"
)

// UsersRepo Port for users
type UsersRepo interface {
	GetUserByID(userID uuid.UUID) (models.User, bool, error)
	GetUserByLogin(login string) (models.User, bool, error)
	GetUserByLoginAndPassword(login string, password string) (models.User, bool, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User, userID uuid.UUID) (models.User, error)
	DeleteUser(userID uuid.UUID) error
}

// NotesRepo Port for notes
type NotesRepo interface {
	GetNotesByUserID(userID uuid.UUID) ([]models.Note, error)
	GetNoteByID(noteID uint) (models.Note, bool, error)
	CreateNote(note models.Note) (models.Note, error)
	UpdateNote(note models.Note, id uint) (models.Note, error)
	DeleteNote(noteID uint) error
	CountNotesByUser(userID uuid.UUID) (int64, error)
}
