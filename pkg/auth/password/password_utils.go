package password

import "golang.org/x/crypto/bcrypt"

type (

	// PasswordManager Interface that represents a password manager that generates and checks them
	PasswordManager interface {
		GeneratePassword(password string) (string, error)
		CheckPassword(password string, hash string) (bool, error)
	}

	// PasswordManagerBcrypt is a bcrypt implementation of PasswordUtils
	PasswordManagerBcrypt struct{}
)

func NewPasswordManagerBcrypt() *PasswordManagerBcrypt {
	return &PasswordManagerBcrypt{}
}

func (p *PasswordManagerBcrypt) GeneratePassword(password string) (string, error) {
	generatedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(generatedPassword), nil
}

func (p *PasswordManagerBcrypt) CheckPassword(password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
