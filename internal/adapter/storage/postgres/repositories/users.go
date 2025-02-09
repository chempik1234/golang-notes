package repository

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"notes_service/internal/users"
	"notes_service/pkg/storage/postgres"
)

type UsersRepo struct {
	db *postgres.DBInstance
}

func NewUsersRepo(db *postgres.DBInstance) *UsersRepo {
	return &UsersRepo{db}
}

func SetPassword(user *users.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

func (r *UsersRepo) GetUserById(userId uuid.UUID) (users.User, error) {
	var user users.User
	result := r.db.Db.First(&user, userId)
	err := result.Error
	return user, err
}

func (r *UsersRepo) GetUserByLogin(login string) (users.User, error) {
	var user users.User
	result := r.db.Db.First(&user, "login = ?", login)
	err := result.Error
	return user, err
}

func (r *UsersRepo) GetUserByLoginAndPassword(login string, password string) (users.User, error) {
	var user users.User
	result := r.db.Db.First(&user, "login = ? AND password = ?", login, password)
	err := result.Error
	return user, err
}

func (r *UsersRepo) CreateUser(user users.User) (users.User, error) {
	err := SetPassword(&user)
	if err != nil {
		return users.User{}, err
	}
	r.db.Db.Create(&user)
	return user, nil
}

func (r *UsersRepo) UpdateUser(user users.User, userId uuid.UUID) (users.User, error) {
	err := SetPassword(&user)
	if err != nil {
		return users.User{}, err
	}
	r.db.Db.Model(&user).Where("id = ?", userId).Updates(user)
	return user, nil
}

func (r *UsersRepo) DeleteUser(userId uuid.UUID) error {
	r.db.Db.Delete(&users.User{}, userId)
	return nil
}
