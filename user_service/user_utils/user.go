package user_utils

import (
	"fmt"
	"time"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password_hash" json:"password"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	IsAdmin   bool      `db:"is_admin" json:"is_admin"`
}

type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `db:"is_admin" json:"is_admin"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserInput) Validate() error {
	//password
	if len(u.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters, got %d", len(u.Password))
	}

	//username
	if u.Username == "" {
		return fmt.Errorf("username is required")
	}

	//email
	if u.Email == "" {
		return fmt.Errorf("email is required")
	} else if err := checkmail.ValidateFormat(u.Email); err != nil {
		return fmt.Errorf("email is not valid")
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
