package schemas

import (
	"github.com/google/uuid"
)

type NoteBodySchema struct {
	UserID  uuid.UUID `json:"user_id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}

type UserBodySchema struct {
	Login    string `gorm:"type:varchar(20);not null;" json:"login"`
	Password string `gorm:"not null;" json:"password"`
}

type SignUpSchema struct {
	Login    string `gorm:"type:varchar(20);not null;" json:"login"`
	Password string `gorm:"not null;" json:"password"`
}

type RefreshTokenSchema struct {
	RefreshToken string `json:"refresh_token"`
}
