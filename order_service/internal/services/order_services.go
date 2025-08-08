package services

import (
	"errors"
	"fmt"
	"log"
	"order_processing_system/db/psql"
	"order_processing_system/db/redis"
	"order_processing_system/order_service/order_utils"
	"order_processing_system/order_service/order_utils/models"
	"order_processing_system/product_service/utils"
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

func (s *Service) GetOrderById(id string, is_admin bool, user_id int) (*models.OrderDetail, error) {
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

	productsIds, err := s.PSQLRepo.GetOrderProducts(o_id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	orderDeatil := &models.OrderDetail{
		ID:          order.ID,
		UserID:      order.UserID,
		OrderDate:   order.OrderDate,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		Products:    []utils.Product{},
	}

	for _, product := range productsIds {
		product, err := s.PSQLRepo.GetProductByID(product.ProductID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		orderDeatil.Products = append(orderDeatil.Products, product)
	}

	return orderDeatil, nil
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

	if status != "created" && status != "processing" && status != "delivered" && status != "cancelled" {
		return errors.New("invalid status")
	}

	if status == "cancelled" {
		order, err := s.PSQLRepo.GetOrder(o_id)

		if order.Status == status {
			return errors.New("order already cancelled")
		}

		if err != nil {
			log.Println(err)
			return err
		}
		for _, product := range order.Products {
			p_id := product.ProductID
			quantity := product.Quantity
			s.PSQLRepo.IncreaseProductStock(p_id, quantity)
		}
	}

	return s.PSQLRepo.PutOrderStatus(o_id, status)
}
