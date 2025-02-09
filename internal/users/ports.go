package users

import (
	"errors"
	"github.com/google/uuid"
)

type UsersRepo interface {
	GetUserById(userId uuid.UUID) (User, error)
	GetUserByLogin(login string) (User, error)
	GetUserByLoginAndPassword(login string, password string) (User, error)
	CreateUser(user User) (User, error)
	UpdateUser(user User, userId uuid.UUID) (User, error)
	DeleteUser(userId uuid.UUID) error
}

var ErrNoteNotFound = errors.New("note not found")
