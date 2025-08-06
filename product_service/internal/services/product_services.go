package services

import (
	"encoding/json"
	"fmt"
	"order_processing_system/db/psql"
	"order_processing_system/db/redis"
	"order_processing_system/product_service/utils"
	"strconv"
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

func (s *Service) GetAllProducts() ([]utils.Product, error) {
	cacheKey := "products_all"
	cached, err := s.RedisRepo.GetData(cacheKey)
	if err == nil {
		var products []utils.Product
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			return products, nil
		}
	}

	products, err := s.PSQLRepo.GetProductsList()
	if err != nil {
		return []utils.Product{}, err
	}

	jsonData, err := json.Marshal(products)
	if err == nil {
		s.RedisRepo.SetCache(cacheKey, jsonData)
	}
	return products, nil
}

func (s *Service) GetProduct(id string) (utils.Product, error) {
	product_id, err := strconv.Atoi(id)
	if err != nil {
		return utils.Product{}, err
	}

	cacheKey := "product_" + id
	cached, err := s.RedisRepo.GetData(cacheKey)
	if err == nil {
		var product utils.Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			return product, nil
		}
	}
	product, err := s.PSQLRepo.GetProductByID(product_id)
	if err != nil {
		return utils.Product{}, err
	}

	jsonData, err := json.Marshal(product)
	if err == nil {
		s.RedisRepo.SetCache(cacheKey, jsonData)
	}
	return product, nil

}

func (s *Service) GetProductStock(id string) (utils.ProductStock, error) {
	product_id, err := strconv.Atoi(id)
	if err != nil {
		return utils.ProductStock{}, err
	}
	cacheKey := "product_stock_" + id
	cached, err := s.RedisRepo.GetData(cacheKey)
	if err == nil {
		var stock utils.ProductStock
		if err := json.Unmarshal([]byte(cached), &stock); err == nil {
			return stock, nil
		}
	}

	productStock, err := s.PSQLRepo.GetProductQuantity(product_id)
	if err != nil {
		return utils.ProductStock{}, err
	}

	jsonData, err := json.Marshal(productStock)
	if err == nil {
		s.RedisRepo.SetCache(cacheKey, jsonData)
	}
	return s.PSQLRepo.GetProductQuantity(product_id)
}

func (s *Service) CreateProduct(product *utils.Product) error {
	err := product.Validate()
	if err != nil {
		return err
	}
	return s.PSQLRepo.PostProduct(product)
}

func (s *Service) UpdateProduct(newProduct utils.Product) (utils.Product, error) {
	err := newProduct.Validate()
	if err != nil {
		return utils.Product{}, err
	}
	cacheKey := fmt.Sprintf("product_%d", newProduct.ID)
	s.RedisRepo.DeleteProduct(cacheKey)

	return s.PSQLRepo.PutProduct(newProduct)
}

func (s *Service) RemoveProduct(id string) error {
	product_id, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	cacheKey := "product_" + id
	s.RedisRepo.DeleteProduct(cacheKey)

	return s.PSQLRepo.DeleteProduct(product_id)
}
