package schemas

import (
	"github.com/google/uuid"
)

type NoteBodySchema struct {
	UserId  uuid.UUID `json:"user_id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}
