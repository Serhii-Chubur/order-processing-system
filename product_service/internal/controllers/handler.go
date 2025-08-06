package controllers

import (
	"encoding/json"
	"net/http"
	"order_processing_system/product_service/internal/services"

	"github.com/gorilla/mux"
)

type Controller struct {
	ch chan error
	s  *services.Service
}

type Handler interface {
	ProductList(w http.ResponseWriter, r *http.Request)
	ProductDetail(w http.ResponseWriter, r *http.Request)
	ProductStock(w http.ResponseWriter, r *http.Request)
	ProductCreate(w http.ResponseWriter, r *http.Request)
	ProductUpdate(w http.ResponseWriter, r *http.Request)
	ProductDelete(w http.ResponseWriter, r *http.Request)
}

func NewController(ch chan error, s *services.Service) *Controller {
	return &Controller{
		ch: ch,
		s:  s,
	}
}

func (c *Controller) ProductList(w http.ResponseWriter, r *http.Request) {

	productsData, err := c.s.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(productsData)
	if err != nil {
		c.ch <- err
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func (c *Controller) ProductDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	productData, err := c.s.GetProduct(id)
	if err != nil {
		http.Error(w, "No such product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(productData)
	if err != nil {
		c.ch <- err
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func (c *Controller) ProductStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	productData, err := c.s.GetProductStock(id)
	if err != nil {
		http.Error(w, "No such product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(productData)
	if err != nil {
		c.ch <- err
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
func (c *Controller) ProductCreate(w http.ResponseWriter, r *http.Request) {}
func (c *Controller) ProductUpdate(w http.ResponseWriter, r *http.Request) {}
func (c *Controller) ProductDelete(w http.ResponseWriter, r *http.Request) {}
