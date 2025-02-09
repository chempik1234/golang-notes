package notes

import (
	"github.com/google/uuid"
	"notes_service/internal/notes"
)

type NoteCRUDUseCase struct {
	notesRepo notes.NotesRepo
}

func NewNoteCRUDUseCase(notesRepo notes.NotesRepo) *NoteCRUDUseCase {
	return &NoteCRUDUseCase{notesRepo: notesRepo}
}

func (u *NoteCRUDUseCase) GetNotesByUserId(userId uint) ([]notes.Note, error) {
	return u.notesRepo.GetNotesByUserId(userId)
}

func (u *NoteCRUDUseCase) GetNoteById(noteId uint) (notes.Note, error) {
	return u.notesRepo.GetNoteById(noteId)
}

func (u *NoteCRUDUseCase) Create(note notes.Note) (notes.Note, error) {
	return u.notesRepo.CreateNote(note)
}

func (u *NoteCRUDUseCase) Update(note notes.Note, id uint) (notes.Note, error) {
	return u.notesRepo.UpdateNote(note, id)
}

func (u *NoteCRUDUseCase) DeleteNote(noteId uint) error {
	return u.notesRepo.DeleteNote(noteId)
}

func (u *NoteCRUDUseCase) CountNotesByUser(userId uuid.UUID) (int64, error) {
	return u.notesRepo.CountNotesByUser(userId)
}
