package notes

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	UserId  uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Title   string    `json:"title" gorm:"text;not null;default:'title'"`
	Content string    `json:"content" gorm:"text;not null;default:null"`
}
