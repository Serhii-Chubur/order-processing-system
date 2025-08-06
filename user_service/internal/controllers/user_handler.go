package controllers

import (
	"net/http"
	"order_processing_system/user_service/internal/services"
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

func (c *Controller) RegisterUser(w http.ResponseWriter, r *http.Request) {}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {}

func (c *Controller) GetUserProfile(w http.ResponseWriter, r *http.Request) {}

func (c *Controller) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {}

func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) {}
