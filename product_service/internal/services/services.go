package services

import (
	"order_processing_system/db/psql"
)

type Service struct {
	// RedisRepo *redis.RedisRepo
	PSQLRepo *psql.PostgresRepo
}

func NewService(psqlRepo *psql.PostgresRepo) *Service {
	// redisRepo *redis.RedisRepo,
	return &Service{
		// RedisRepo: redisRepo,
		PSQLRepo: psqlRepo,
	}
}

func (s *Service) GetAllProducts() ([]psql.Product, error) {
	return s.PSQLRepo.GetProductsList()
}
