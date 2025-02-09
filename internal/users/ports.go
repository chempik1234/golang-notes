package users

import (
	"errors"
	"github.com/google/uuid"
)

type UsersRepo interface {
	GetUserById(userId uuid.UUID) (User, error)
	GetUserByLogin(login string) (User, error)
	CreateUser(user User) (User, error)
	UpdateUser(user User, userId uint) (User, error)
	DeleteUser(userId uint) error
}

var ErrNoteNotFound = errors.New("note not found")
