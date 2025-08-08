package services

import (
	"order_processing_system/db/psql"
	"order_processing_system/db/redis"
	"order_processing_system/order_service/order_utils"
	"order_processing_system/order_service/order_utils/models"
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

func (s *Service) CreateOrder(orderData *models.OrderInput, user_id int) (*models.Order, error) {
	var order models.Order
	order.UserID = user_id
	order.Status = "created"
	order.Products = orderData.Products
	order.OrderDate = time.Now()
	amount, err := order_utils.CalculateTotalAmount(&order, s.PSQLRepo)
	if err != nil {
		return nil, err
	}
	order.TotalAmount = amount

	return s.PSQLRepo.PostOrder(&order)
}
