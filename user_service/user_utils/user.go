package user_utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/badoux/checkmail"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	err := godotenv.Load("db/configs/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type User struct {
	ID        int       `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password_hash" json:"password"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	IsAdmin   bool      `db:"is_admin" json:"is_admin"`
}

type UserInfo struct {
	ID        int       `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	IsAdmin   bool      `db:"is_admin" json:"is_admin"`
}

type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `db:"is_admin" json:"is_admin"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
	jwt.RegisteredClaims
}

type TokenRequest struct {
	RefreshToken string `json:"refresh_token"`
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

func CheckPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func GenerateAccessToken(user User) (string, error) {
	accessClaims := &Claims{
		Email: user.Email,
		ID:    user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}

	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return access, nil
}

func GenerateRefreshToken(user User) (string, error) {
	refreshClaims := &Claims{
		Email: user.Email,
		ID:    user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return refresh, nil
}

// func ParseToken(tokenStr string) (*Claims, error) {
// 	fmt.Println("JWT_SECRET:", os.Getenv("JWT_SECRET"))

// 	fmt.Println(tokenStr)
// 	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(os.Getenv("JWT_SECRET")), nil
// 	})
// 	fmt.Println(token.Claims)

// 	if err != nil {
// 		return nil, err
// 	}

// 	claims, ok := token.Claims.(*Claims)
// 	fmt.Println(claims)
// 	if !ok || !token.Valid {
// 		return nil, errors.New("invalid token")
// 	}

// 	return claims, nil
// }
