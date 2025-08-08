package services

import (
	"errors"
	"fmt"
	"log"
	"order_processing_system/db/psql"
	"order_processing_system/db/redis"
	"order_processing_system/order_service/order_utils"
	"order_processing_system/order_service/order_utils/models"
	"strconv"
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

func (s *Service) GetOrderById(id string, is_admin bool, user_id int) (*models.Order, error) {
	o_id, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	order, err := s.PSQLRepo.GetOrder(o_id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !is_admin && order.UserID != user_id {
		return nil, errors.New("forbidden access to another user's order")
	}

	return order, nil
}

func (s *Service) GetOrdersByUserId(id string, is_admin bool, user_id int) ([]models.Order, error) {
	u_id, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	orders, err := s.PSQLRepo.GetUserOrders(u_id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf("no orders found for user id %d", u_id)
	}

	if !is_admin && u_id != user_id {
		return nil, errors.New("forbidden access to another user's orders")
	}

	return orders, nil
}

func (s *Service) UpdateOrderStatus(id string, status string) error {
	o_id, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return err
	}

	return s.PSQLRepo.PutOrderStatus(o_id, status)
}
