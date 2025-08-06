package controllers

import (
	"encoding/json"
	"net/http"
	"order_processing_system/user_service/internal/services"
	"order_processing_system/user_service/user_utils"
)

type Controller struct {
	ch chan error
	s  *services.Service
}

type Handler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetUserProfile(w http.ResponseWriter, r *http.Request)
	UpdateUserProfile(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

func NewController(ch chan error, s *services.Service) *Controller {
	return &Controller{
		ch: ch,
		s:  s,
	}
}

func (c *Controller) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userData user_utils.UserInput
	err := json.NewDecoder(r.Body).Decode(&userData)

	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	err = userData.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.s.NewUser(&userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {}

func (c *Controller) GetUserProfile(w http.ResponseWriter, r *http.Request) {}

func (c *Controller) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {}

func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) {}
