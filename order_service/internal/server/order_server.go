package server

import (
	"fmt"
	"net/http"
	"order_processing_system/order_service/internal/controllers"
	"order_processing_system/order_service/internal/middleware"

	"github.com/gorilla/mux"
)

func NewServer(c *controllers.Controller) *http.Server {
	r := mux.NewRouter().StrictSlash(true)

	orderRouter := r.PathPrefix("/api/orders").Subrouter()
	orderRouter.Use(middleware.IsAuthenticated)

	orderRouter.HandleFunc("", c.OrderList).Methods("POST")
	orderRouter.HandleFunc("/{id}", c.OrderDetail).Methods("GET")
	orderRouter.HandleFunc("/user/{id}", c.UserOrders).Methods("GET")

	adminRouter := r.PathPrefix("/api/orders").Subrouter()
	adminRouter.Use(middleware.IsAdmin)

	adminRouter.HandleFunc("/{id}/status", c.UpdateOrderStatus).Methods("PUT")

	fmt.Println("http://localhost:8002/api/orders/")

	serv := &http.Server{
		Addr:    ":8002",
		Handler: r,
	}

	return serv
}

func StartServer(serv *http.Server) {
	serv.ListenAndServe()
}

func StopServer(serv *http.Server) {
	serv.Close()
}
