package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the User model in the database
type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key;" json:"id"`
	Login    string    `gorm:"type:varchar(20);not null;" json:"login"`
	Password string    `gorm:"not null;" json:"password"`

	Notes []Note `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

// Note represents the Note model in the database
type Note struct {
	gorm.Model
	UserID  uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Title   string    `json:"title" gorm:"text;not null;default:'title'"`
	Content string    `json:"content" gorm:"text;not null;default:null"`
}
