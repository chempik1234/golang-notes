package usecases

import (
	"github.com/google/uuid"
	"notes_service/internal/models"
	"notes_service/internal/ports"
)

// region Notes

// NoteUseCase represents the use case for managing notes.
// Since it's simple, it actually implements the port
type NoteUseCase struct {
	notesRepo ports.NotesRepo
}

// NewNoteCRUDUseCase create and return a new NoteUseCase
func NewNoteCRUDUseCase(notesRepo ports.NotesRepo) *NoteUseCase {
	return &NoteUseCase{notesRepo: notesRepo}
}

// GetNotesByUserID returns a list of notes for the specified user by their ID.
func (u *NoteUseCase) GetNotesByUserID(userID uuid.UUID) ([]models.Note, error) {
	return u.notesRepo.GetNotesByUserID(userID)
}

// GetNoteByID returns a note by its unique ID.
func (u *NoteUseCase) GetNoteByID(noteID uint) (models.Note, bool, error) {
	return u.notesRepo.GetNoteByID(noteID)
}

// Create creates a new note and returns it.
func (u *NoteUseCase) Create(note models.Note) (models.Note, error) {
	return u.notesRepo.CreateNote(note)
}

// Update updates an existing note by its ID and returns the updated note.
func (u *NoteUseCase) Update(note models.Note, id uint) (models.Note, error) {
	return u.notesRepo.UpdateNote(note, id)
}

// DeleteNote deletes a note by its ID.
func (u *NoteUseCase) DeleteNote(noteID uint) error {
	return u.notesRepo.DeleteNote(noteID)
}

// CountNotesByUser returns the count of notes for the specified user by their UUID.
func (u *NoteUseCase) CountNotesByUser(userID uuid.UUID) (int64, error) {
	return u.notesRepo.CountNotesByUser(userID)
}

//endregion

// region Users

// UserUseCase represents a use case for managing users.
type UserUseCase struct {
	usersRepo ports.UsersRepo // Repository for managing users
}

// NewUserCRUDUseCase creates and returns a new instance of UserUseCase.
func NewUserCRUDUseCase(usersRepo ports.UsersRepo) *UserUseCase {
	return &UserUseCase{usersRepo: usersRepo}
}

// GetUserByID returns a user by their unique UUID.
func (u *UserUseCase) GetUserByID(userID uuid.UUID) (models.User, bool, error) {
	return u.usersRepo.GetUserByID(userID)
}

// GetUserByLogin returns a user by their login.
func (u *UserUseCase) GetUserByLogin(login string) (models.User, bool, error) {
	return u.usersRepo.GetUserByLogin(login)
}

// GetUserByLoginAndPassword searches for a user by login and password.
// Returns the found user or an error if the user is not found.
func (u *UserUseCase) GetUserByLoginAndPassword(login string, password string) (models.User, bool, error) {
	return u.usersRepo.GetUserByLoginAndPassword(login, password)
}

// CreateUser creates a new user and returns it.
func (u *UserUseCase) CreateUser(user models.User) (models.User, error) {
	return u.usersRepo.CreateUser(user)
}

// UpdateUser updates an existing user by their UUID and returns the updated user.
func (u *UserUseCase) UpdateUser(user models.User, userID uuid.UUID) (models.User, error) {
	return u.usersRepo.UpdateUser(user, userID)
}

// DeleteUser deletes a user by their UUID.
func (u *UserUseCase) DeleteUser(userID uuid.UUID) error {
	return u.usersRepo.DeleteUser(userID)
}

//endregion
