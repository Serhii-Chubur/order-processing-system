package controllers

import (
	"encoding/json"
	"net/http"
	"order_processing_system/user_service/internal/services"
	"order_processing_system/user_service/user_utils"
	"strings"

	"github.com/gorilla/mux"
)

type Controller struct {
	ch chan error
	s  *services.Service
}

type Handler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
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

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var loginData user_utils.LoginInput
	err := json.NewDecoder(r.Body).Decode(&loginData)

	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	user, err := c.s.GetRegisteredUser(loginData.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusNotFound)
		return
	}

	err = user_utils.CheckPassword(loginData.Password, user.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := c.s.GenerateTokens(user)
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}
	resp := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (c *Controller) Refresh(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	const prefix = "Bearer "

	token := strings.TrimPrefix(authHeader, prefix)
	token = strings.TrimSpace(token)

	err := c.s.DeleteToken(token)
	if err != nil {
		http.Error(w, "Failed to delete access token from Redis", http.StatusUnauthorized)
		return
	}

	var tokenData user_utils.TokenRequest
	err = json.NewDecoder(r.Body).Decode(&tokenData)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	currToken := tokenData.RefreshToken
	email, err := c.s.GetEmail(currToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusBadRequest)
		return
	}

	user, err := c.s.GetRegisteredUser(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	accessToken, refreshToken, err := c.s.GenerateTokens(user)
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	err = c.s.DeleteToken(currToken)
	if err != nil {
		http.Error(w, "Failed to delete refresh token from Redis", http.StatusUnauthorized)
		return
	}

	resp := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (c *Controller) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	const prefix = "Bearer "

	token := strings.TrimPrefix(authHeader, prefix)
	token = strings.TrimSpace(token)

	email, err := c.s.GetEmail(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	userInfo, err := c.s.PSQLRepo.GetUserInfo(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userInfo)

}

func (c *Controller) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	const prefix = "Bearer "

	token := strings.TrimPrefix(authHeader, prefix)
	token = strings.TrimSpace(token)

	email, err := c.s.GetEmail(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	requestUser, err := c.s.GetRegisteredUser(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	currUser, err := c.s.GetUser(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var userData user_utils.UserInput
	err = json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if currUser.Email != requestUser.Email && !requestUser.IsAdmin {
		http.Error(w, "Forbidden access to another user", http.StatusForbidden)
		return
	}

	err = userData.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = c.s.UpdateUserInfo(userData, id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	const prefix = "Bearer "

	accessToken := strings.TrimPrefix(authHeader, prefix)
	accessToken = strings.TrimSpace(accessToken)

	var req user_utils.TokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.RefreshToken == "" {
		http.Error(w, "Invalid refresh token", http.StatusBadRequest)
		return
	}

	refreshToken := req.RefreshToken

	err = c.s.DeleteToken(accessToken)
	if err != nil {
		http.Error(w, "Failed to delete access token from Redis", http.StatusUnauthorized)
		return
	}
	err = c.s.DeleteToken(refreshToken)
	if err != nil {
		http.Error(w, "Failed to delete refresh token from Redis", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
