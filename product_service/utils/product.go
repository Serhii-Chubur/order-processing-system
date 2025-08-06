package utils

import (
	"fmt"
)

type Product struct {
	ID            int     `db:"id" json:"id"`
	Name          string  `db:"name" json:"name"`
	Description   string  `db:"description" json:"description"`
	Price         float64 `db:"price" json:"price"`
	StockQuantity int     `db:"stock_quantity" json:"stock"`
}

type ProductStock struct {
	ID            int `db:"id" json:"id"`
	StockQuantity int `db:"stock_quantity" json:"stock"`
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if p.StockQuantity <= 0 {
		return fmt.Errorf("stock quantity must be greater than 0")
	}
	return nil
}
