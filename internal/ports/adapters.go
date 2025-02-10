package ports

import (
	"errors"
	"github.com/go-redis/redis/v7"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"notes_service/internal/models"
	"notes_service/pkg/auth/password"
	"notes_service/pkg/storage/postgres"
	"time"
)

//region Notes

// NotesRepoDB represents a repository interface for interacting with the Notes database.
type NotesRepoDB struct {
	db *postgres.DBInstance // Database instance for executing queries.
}

var _ NotesRepo = (*NotesRepoDB)(nil)

// NewNotesRepoDB initializes and returns a new NotesRepoDB instance.
func NewNotesRepoDB(db *postgres.DBInstance) *NotesRepoDB {
	return &NotesRepoDB{db}
}

// GetNotesByUserID retrieves all notes associated with a specific user ID.
func (r *NotesRepoDB) GetNotesByUserID(userID uuid.UUID) ([]models.Note, error) {
	var notesList []models.Note
	r.db.Db.Find(&notesList, "user_id = ?", userID)
	return notesList, nil
}

// GetNoteByID retrieves a single note by its unique ID.
func (r *NotesRepoDB) GetNoteByID(noteID uint) (models.Note, bool, error) {
	var note models.Note
	result := r.db.Db.First(&note, noteID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return note, false, nil
		}
		return note, false, result.Error
	}
	return note, true, nil
}

// CreateNote creates a new note and saves it to the database.
func (r *NotesRepoDB) CreateNote(note models.Note) (models.Note, error) {
	r.db.Db.Create(&note)
	return note, nil
}

// UpdateNote updates an existing note in the database based on its ID.
func (r *NotesRepoDB) UpdateNote(note models.Note, noteID uint) (models.Note, error) {
	r.db.Db.Model(&note).Where("id = ?", noteID).Updates(note)
	return note, nil
}

// DeleteNote deletes a note from the database by its ID.
func (r *NotesRepoDB) DeleteNote(noteID uint) error {
	r.db.Db.Delete(&models.Note{}, noteID)
	return nil
}

// CountNotesByUser counts the number of notes associated with a specific user ID.
func (r *NotesRepoDB) CountNotesByUser(noteID uuid.UUID) (int64, error) {
	var count int64
	r.db.Db.Model(&models.Note{}).Where("user_id = ?", noteID).Count(&count)
	return count, nil
}

//endregion

//region Users

// UsersRepoDB represents a repository interface for interacting with the Users database.
type UsersRepoDB struct {
	db              *postgres.DBInstance     // Database instance for executing queries.
	passwordManager password.PasswordManager // Password utils instance for generating and checking passwords
}

var _ UsersRepo = (*UsersRepoDB)(nil)

// NewUsersRepoDB initializes and returns a new UsersRepoDB instance.
func NewUsersRepoDB(db *postgres.DBInstance, passwordManager password.PasswordManager) *UsersRepoDB {
	return &UsersRepoDB{
		db,
		passwordManager,
	}
}

// SetPassword hashes the password for a user using bcrypt and sets it in the user object.
func (r *UsersRepoDB) SetPassword(user *models.User) error {
	hashedPassword, err := r.passwordManager.GeneratePassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

// GetUserByID retrieves a user by their unique UUID.
func (r *UsersRepoDB) GetUserByID(userID uuid.UUID) (models.User, bool, error) {
	var user models.User
	result := r.db.Db.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, false, nil
		}
		return user, false, result.Error
	}
	return user, true, nil
}

// GetUserByLogin retrieves a user by their login name.
func (r *UsersRepoDB) GetUserByLogin(login string) (models.User, bool, error) {
	var user models.User
	result := r.db.Db.First(&user, "login = ?", login)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, false, nil
		}
		return user, false, result.Error
	}
	return user, true, nil
}

// GetUserByLoginAndPassword retrieves a user by their login name and password.
func (r *UsersRepoDB) GetUserByLoginAndPassword(login string, password string) (models.User, bool, error) {
	var user models.User
	result := r.db.Db.First(&user, "login = ?", login)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, false, nil
		}
		return user, false, result.Error
	} else if passwordsEqual, err := r.passwordManager.CheckPassword(password, user.Password); !passwordsEqual || err != nil {
		return user, false, nil
	}
	return user, true, nil
}

// CreateUser creates a new user and saves it to the database, hashing the password before saving.
func (r *UsersRepoDB) CreateUser(user models.User) (models.User, error) {
	err := r.SetPassword(&user)
	if err != nil {
		return models.User{}, err
	}
	r.db.Db.Create(&user)
	return user, nil
}

// UpdateUser updates an existing user in the database based on their UUID, rehashing the password if changed.
func (r *UsersRepoDB) UpdateUser(user models.User, userID uuid.UUID) (models.User, error) {
	err := r.SetPassword(&user)
	if err != nil {
		return models.User{}, err
	}
	r.db.Db.Model(&user).Where("id = ?", userID).Updates(user)
	return user, nil
}

// DeleteUser deletes a user from the database by their UUID.
func (r *UsersRepoDB) DeleteUser(userID uuid.UUID) error {
	r.db.Db.Delete(&models.User{}, userID)
	return nil
}

//endregion

//region Redis

// RedisStorage represents a storage interface for interacting with Redis.
type RedisStorage struct {
	db *redis.Client // Redis client instance for executing commands.
}

var _ fiber.Storage = (*RedisStorage)(nil)

// NewRedisStorage initializes and returns a new RedisStorage instance.
func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{db: client}
}

// Get retrieves a value from Redis by its key.
func (r RedisStorage) Get(key string) ([]byte, error) {
	return r.db.Get(key).Bytes()
}

// Set stores a value in Redis with an expiration time.
func (r RedisStorage) Set(key string, val []byte, exp time.Duration) error {
	return r.db.Set(key, string(val), exp).Err()
}

// Delete removes a key-value pair from Redis.
func (r RedisStorage) Delete(key string) error {
	return r.db.Del(key).Err()
}

// Reset clears all data from the Redis database.
func (r RedisStorage) Reset() error {
	return r.db.FlushDB().Err()
}

// Close closes the connection to the Redis server.
func (r RedisStorage) Close() error {
	return r.db.Close()
}

//endregion
