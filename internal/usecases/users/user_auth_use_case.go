package users

import (
	"notes_service/internal/users"
)

type UserAuthUseCase struct {
	usersRepo users.UsersRepo
}

func NewUserAuthUseCase(usersRepo users.UsersRepo) *UserAuthUseCase {
	return &UserAuthUseCase{usersRepo: usersRepo}
}

func (u *UserAuthUseCase) FinUserByLoginAndPassword(login string, password string) (users.User, error) {
	return u.usersRepo.GetUserByLoginAndPassword(login, password)
}

func (u *UserAuthUseCase) CreateUser(user users.User) (users.User, error) {
	return u.usersRepo.CreateUser(user)
}
