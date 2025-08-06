package services

import (
	"fmt"
	"order_processing_system/db/psql"
	"order_processing_system/db/redis"
	"order_processing_system/user_service/user_utils"
	"strconv"
	"time"

	"github.com/badoux/checkmail"
)

type Service struct {
	RedisRepo *redis.RedisRepo
	PSQLRepo  *psql.PostgresRepo
}

func NewService(psqlRepo *psql.PostgresRepo, redisRepo *redis.RedisRepo) *Service {
	return &Service{
		RedisRepo: redisRepo,
		PSQLRepo:  psqlRepo,
	}
}

func (s *Service) NewUser(userData *user_utils.UserInput) error {
	var user user_utils.User
	user.Username = userData.Username
	user.Email = userData.Email
	hashedPassword, err := user_utils.HashPassword(userData.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	user.IsAdmin = userData.IsAdmin
	user.CreatedAt = time.Now()

	return s.PSQLRepo.PostUser(&user)
}

func (s *Service) GetRegisteredUser(email string) (user_utils.User, error) {
	if email == "" {
		return user_utils.User{}, fmt.Errorf("email is required")
	} else if err := checkmail.ValidateFormat(email); err != nil {
		return user_utils.User{}, fmt.Errorf("email is not valid")
	}
	return s.PSQLRepo.GetUserByEmail(email)
}

func (s *Service) GenerateTokens(user user_utils.User) (string, string, error) {
	accessToken, err := user_utils.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := user_utils.GenerateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	err = s.RedisRepo.SetAccessToken(user.Email, accessToken)
	if err != nil {
		return "", "", err
	}

	err = s.RedisRepo.SetRefreshToken(user.Email, refreshToken)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *Service) GetEmail(token string) (string, error) {
	return s.RedisRepo.GetUserEmail(token)
}
func (s *Service) GetUser(id string) (user_utils.User, error) {
	u_id, err := strconv.Atoi(id)
	if err != nil {
		return user_utils.User{}, err
	}
	return s.PSQLRepo.GetUserById(u_id)
}

func (s *Service) UpdateUserInfo(userData user_utils.UserInput, id string) error {
	u_id, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	hashedPassword, err := user_utils.HashPassword(userData.Password)
	if err != nil {
		return err
	}
	userData.Password = hashedPassword
	return s.PSQLRepo.PutUser(&userData, u_id)

}

func (s *Service) DeleteToken(token string) error {
	return s.RedisRepo.Delete(token)
}

// func (s *Service) RevokeToken(email string) error {
// 	err := s.RedisRepo.DeleteAccessToken(email)
// 	if err != nil {
// 		return err
// 	}
// 	err = s.RedisRepo.DeleteRefreshToken(email)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
