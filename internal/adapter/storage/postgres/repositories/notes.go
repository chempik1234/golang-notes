package repository

import (
	"github.com/google/uuid"
	"notes_service/internal/adapter/storage/postgres"
	"notes_service/internal/notes"
)

type NotesRepo struct {
	db *postgres.DBInstance
}

func NewNotesRepo(db *postgres.DBInstance) *NotesRepo {
	return &NotesRepo{db}
}

func (r *NotesRepo) GetNotesByUserId(userId uint) ([]notes.Note, error) {
	var notesList []notes.Note
	r.db.Db.Find(&notesList, "user_id = ?", userId)
	return notesList, nil
}

func (r *NotesRepo) GetNoteById(noteId uint) (notes.Note, error) {
	var note notes.Note
	r.db.Db.First(&note, noteId)
	return note, nil
}

func (r *NotesRepo) CreateNote(note notes.Note) (notes.Note, error) {
	r.db.Db.Create(&note)
	return note, nil
}

func (r *NotesRepo) UpdateNote(note notes.Note, noteId uint) (notes.Note, error) {
	r.db.Db.Model(&note).Where("id = ?", noteId).Updates(note)
	return note, nil
}

func (r *NotesRepo) DeleteNote(noteId uint) error {
	r.db.Db.Delete(&notes.Note{}, noteId)
	return nil
}

func (r *NotesRepo) CountNotesByUser(noteId uuid.UUID) (int64, error) {
	var count int64
	r.db.Db.Model(&notes.Note{}).Where("user_id = ?", noteId).Count(&count)
	return count, nil
}
