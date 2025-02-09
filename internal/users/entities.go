package users

import (
	"github.com/google/uuid"
	"notes_service/internal/notes"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key;" json:"id"`
	Login    string    `gorm:"type:varchar(20);not null;" json:"login"`
	Password string    `gorm:"not null;" json:"password"`

	Notes []notes.Note `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}
