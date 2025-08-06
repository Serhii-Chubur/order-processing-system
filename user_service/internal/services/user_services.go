package services

import (
	"order_processing_system/db/psql"
	"order_processing_system/db/redis"
	"order_processing_system/user_service/user_utils"
	"time"
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
