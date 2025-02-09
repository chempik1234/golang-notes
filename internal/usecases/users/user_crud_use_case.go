package users

import "notes_service/internal/users"

type UserCRUDUseCase struct {
	usersRepo users.UsersRepo
}

func (u *UserCRUDUseCase) GetUserById(userId uint) (users.User, error) {
	return u.usersRepo.GetUserById(userId)
}

func (u *UserCRUDUseCase) GetUserByLogin(login string) (users.User, error) {
	return u.usersRepo.GetUserByLogin(login)
}

func (u *UserCRUDUseCase) CreateUser(user users.User) (users.User, error) {
	return u.usersRepo.CreateUser(user)
}

func (u *UserCRUDUseCase) UpdateUser(user users.User, userId uint) (users.User, error) {
	return u.usersRepo.UpdateUser(user, userId)
}

func (u *UserCRUDUseCase) DeleteUser(userId uint) error {
	return u.usersRepo.DeleteUser(userId)
}
