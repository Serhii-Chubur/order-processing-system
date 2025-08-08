package order_utils

import (
	"fmt"
	"order_processing_system/db/psql"
	"order_processing_system/order_service/order_utils/models"
	"order_processing_system/product_service/utils"
)

func GetAvailableProductAmount(productID int, p *psql.PostgresRepo) (utils.ProductStock, error) {
	amount, err := p.GetProductQuantity(productID)
	if err != nil {
		return utils.ProductStock{}, err
	}
	return amount, nil
}

func Validate(o *models.OrderInput, p *psql.PostgresRepo) error {
	for _, product := range o.Products {
		available, err := GetAvailableProductAmount(product.ProductID, p)
		if err != nil {
			return err
		}
		if available.StockQuantity <= 0 {
			return fmt.Errorf("product is out of stock")
		}
		if product.Quantity <= 0 {
			return fmt.Errorf("quantity must be greater than 0")
		}
		if product.Quantity > available.StockQuantity {
			return fmt.Errorf("not enough stock for product, available stock: %d", available.StockQuantity)
		}
	}
	return nil
}

func CalculateTotalAmount(o *models.Order, p *psql.PostgresRepo) (float64, error) {
	totalAmount := 0.0
	for i, product := range o.Products {
		product, err := p.GetProductByID(product.ProductID)
		if err != nil {
			return 0, err
		}
		totalAmount += product.Price * float64(o.Products[i].Quantity)
	}
	return totalAmount, nil
}
