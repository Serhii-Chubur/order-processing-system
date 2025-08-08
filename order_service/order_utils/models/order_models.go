package models

import (
	"order_processing_system/product_service/utils"
	"time"
)

type Order struct {
	ID          int            `db:"id" json:"id"`
	UserID      int            `db:"user_id" json:"user_id"`
	OrderDate   time.Time      `db:"order_date" json:"order_date"`
	Status      string         `db:"status" json:"status"`
	TotalAmount float64        `db:"total_amount" json:"total_amount"`
	Products    []OrderProduct `json:"products"`
}

type OrderDetail struct {
	ID          int             `db:"id" json:"id"`
	UserID      int             `db:"user_id" json:"user_id"`
	OrderDate   time.Time       `db:"order_date" json:"order_date"`
	Status      string          `db:"status" json:"status"`
	TotalAmount float64         `db:"total_amount" json:"total_amount"`
	Products    []utils.Product `json:"products"`
}

type OrderInput struct {
	Products []OrderProduct `json:"products"`
}

type OrderProduct struct {
	ProductID int `db:"product_id" json:"product_id"`
	Quantity  int `db:"quantity" json:"quantity"`
}

type StatusUpdate struct {
	Status string `json:"status"`
}
