package services

import (
	"order_processing_system/db/psql"
	"strconv"
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

func (s *Service) GetProduct(id string) (psql.Product, error) {
	product_id, err := strconv.Atoi(id)
	if err != nil {
		return psql.Product{}, err
	}
	return s.PSQLRepo.GetProductByID(product_id)

}

func (s *Service) GetProductStock(id string) (psql.ProductStock, error) {
	product_id, err := strconv.Atoi(id)
	if err != nil {
		return psql.ProductStock{}, err
	}
	return s.PSQLRepo.GetProductQuantity(product_id)
}

func (s *Service) CreateProduct(product *psql.Product) error {
	err := product.Validate()
	if err != nil {
		return err
	}
	return s.PSQLRepo.PostProduct(product)
}

func (s *Service) UpdateProduct(newProduct psql.Product) (psql.Product, error) {
	err := newProduct.Validate()
	if err != nil {
		return psql.Product{}, err
	}
	return s.PSQLRepo.PutProduct(newProduct)
}

func (s *Service) RemoveProduct(id string) error {
	product_id, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.PSQLRepo.DeleteProduct(product_id)
}
