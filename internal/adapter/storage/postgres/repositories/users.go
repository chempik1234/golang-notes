package repository

import (
	"github.com/google/uuid"
	"notes_service/internal/adapter/storage/postgres"
	"notes_service/internal/users"
)

type UsersRepo struct {
	db *postgres.DBInstance
}

func NewUsersRepo(db *postgres.DBInstance) *UsersRepo {
	return &UsersRepo{db}
}

func (r *UsersRepo) GetUserById(userId uuid.UUID) (users.User, error) {
	var user users.User
	r.db.Db.First(&user, userId)
	return user, nil
}

func (r *UsersRepo) GetUserByLogin(login string) (users.User, error) {
	var user users.User
	r.db.Db.First(&user, login)
	return user, nil
}

func (r *UsersRepo) CreateUser(user users.User) (users.User, error) {
	r.db.Db.Create(&user)
	return user, nil
}

func (r *UsersRepo) UpdateUser(user users.User, userId uint) (users.User, error) {
	r.db.Db.Model(&user).Where("id = ?", userId).Updates(user)
	return user, nil
}

func (r *UsersRepo) DeleteUser(userId uuid.UUID) error {
	r.db.Db.Delete(&users.User{}, userId)
	return nil
}
