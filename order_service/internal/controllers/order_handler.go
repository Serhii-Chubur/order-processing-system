package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"order_processing_system/db/redis"
	"order_processing_system/order_service/internal/services"
	"order_processing_system/order_service/order_utils"
	"order_processing_system/order_service/order_utils/models"
	"strings"

	"github.com/gorilla/mux"
)

type Controller struct {
	ch chan error
	s  *services.Service
}

type Handler interface {
	OrderList(w http.ResponseWriter, r *http.Request)
	OrderDetail(w http.ResponseWriter, r *http.Request)
	UserOrders(w http.ResponseWriter, r *http.Request)
	UpdateOrderStatus(w http.ResponseWriter, r *http.Request)
}

func NewController(ch chan error, s *services.Service) *Controller {
	return &Controller{
		ch: ch,
		s:  s,
	}
}

func (c *Controller) OrderList(w http.ResponseWriter, r *http.Request) {
	var orderData models.OrderInput
	err := json.NewDecoder(r.Body).Decode(&orderData)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
	}

	authHeader := r.Header.Get("Authorization")
	const prefix = "Bearer "
	token := strings.TrimPrefix(authHeader, prefix)
	token = strings.TrimSpace(token)

	userData, err := redis.ParseToken(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = order_utils.Validate(&orderData, c.s.PSQLRepo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := c.s.CreateOrder(&orderData, userData.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, product := range order.Products {
		err = c.s.PSQLRepo.DecreaseProductStock(product.ProductID, product.Quantity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respMsg := fmt.Sprintf("Order %d created successfully", order.ID)
	w.Write([]byte(respMsg))

}

func (c *Controller) OrderDetail(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	authHeader := r.Header.Get("Authorization")
	const prefix = "Bearer "

	token := strings.TrimPrefix(authHeader, prefix)
	token = strings.TrimSpace(token)

	info, _ := redis.ParseToken(token)
	is_admin := info.Root
	u_id := info.ID

	orderData, err := c.s.GetOrderById(id, is_admin, u_id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(orderData)
	if err != nil {
		c.ch <- err
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func (c *Controller) UserOrders(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	authHeader := r.Header.Get("Authorization")
	const prefix = "Bearer "

	token := strings.TrimPrefix(authHeader, prefix)
	token = strings.TrimSpace(token)

	info, _ := redis.ParseToken(token)
	is_admin := info.Root
	u_id := info.ID

	orderData, err := c.s.GetOrdersByUserId(id, is_admin, u_id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(orderData)
	if err != nil {
		c.ch <- err
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func (c *Controller) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var status models.StatusUpdate
	err := json.NewDecoder(r.Body).Decode(&status)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.s.UpdateOrderStatus(id, status.Status)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	respMsg := fmt.Sprintf("Order %s updated successfully", id)
	w.Write([]byte(respMsg))
}
