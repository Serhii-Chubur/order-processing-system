package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"order_processing_system/db/psql"
	"order_processing_system/db/redis"
	"order_processing_system/order_service/internal/natsclient"
	"order_processing_system/order_service/order_utils"
	"order_processing_system/order_service/order_utils/models"
	"order_processing_system/product_service/utils"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

type Service struct {
	RedisRepo  *redis.RedisRepo
	PSQLRepo   *psql.PostgresRepo
	NATSClient *natsclient.OrderNATS
}

func NewService(psqlRepo *psql.PostgresRepo, redisRepo *redis.RedisRepo, natsClient *natsclient.OrderNATS) *Service {
	return &Service{
		RedisRepo:  redisRepo,
		PSQLRepo:   psqlRepo,
		NATSClient: natsClient,
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

	cacheKey := "order_" + id
	cached, err := s.RedisRepo.GetData(cacheKey)
	if err == nil {
		var order *models.OrderDetail
		if err := json.Unmarshal([]byte(cached), &order); err == nil {
			return order, nil
		}
	}

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

	jsonData, err := json.Marshal(orderDeatil)
	if err == nil {
		s.RedisRepo.SetCache(cacheKey, jsonData)
	}

	return orderDeatil, nil
}

func (s *Service) GetOrdersByUserId(id string, is_admin bool, user_id int) ([]models.Order, error) {
	cacheKey := "user_" + id + "_orders"
	cached, err := s.RedisRepo.GetData(cacheKey)
	if err == nil {
		var orders []models.Order
		if err := json.Unmarshal([]byte(cached), &orders); err == nil {
			return orders, nil
		}
	}

	u_id, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !is_admin && u_id != user_id {
		return nil, errors.New("forbidden access to another user's orders")
	}

	orders, err := s.PSQLRepo.GetUserOrders(u_id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf("no orders found for user id %d", u_id)
	}

	for i, order := range orders {
		productsIds, err := s.PSQLRepo.GetOrderProducts(order.ID)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		orders[i].Products = append(orders[i].Products, productsIds...)
	}

	jsonData, err := json.Marshal(orders)
	if err == nil {
		s.RedisRepo.SetCache(cacheKey, jsonData)
	}

	return orders, nil
}

func (s *Service) UpdateOrderStatus(id string, status string) error {
	o_id, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return err
	}

	if status != "processing" && status != "delivered" && status != "cancelled" {
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

func (s *Service) ListenProductUpdates() error {
	s.NATSClient.Subscribe("product.created", func(msg *nats.Msg) {
		var product utils.Product
		err := json.Unmarshal(msg.Data, &product)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Product with id %d created, data: %v", product.ID, product)
	})

	s.NATSClient.Subscribe("product.updated", func(msg *nats.Msg) {
		var product utils.Product
		err := json.Unmarshal(msg.Data, &product)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Product with id %d updated, data: %v", product.ID, product)
	})

	s.NATSClient.Subscribe("product.deleted", func(msg *nats.Msg) {
		id := string(msg.Data)
		log.Printf("Product with id %s deleted", id)
	})

	err := s.NATSClient.Conn.Flush()
	if err != nil {
		return err
	}

	fmt.Println("Listening for 'products' messages...")
	select {}
}
