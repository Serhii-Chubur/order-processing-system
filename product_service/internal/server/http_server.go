package server

import (
	"fmt"
	"net/http"
	"product_service/internal/controllers"

	"github.com/gorilla/mux"
)

func NewServer(c *controllers.Controller) *http.Server {
	r := mux.NewRouter()

	productRouter := r.PathPrefix("api/products").Subrouter()

	productRouter.HandleFunc("", c.ProductList).Methods("GET")
	productRouter.HandleFunc("/{id}", c.ProductDetail).Methods("GET")
	productRouter.HandleFunc("/{id}/stock", c.ProductStock).Methods("GET")

	adminRouter := r.PathPrefix("api/products").Subrouter()
	adminRouter.Use(c.IsAdmin)

	adminRouter.HandleFunc("", c.ProductCreate).Methods("POST")        // admin
	adminRouter.HandleFunc("/{id}", c.ProductUpdate).Methods("PUT")    // admin
	adminRouter.HandleFunc("/{id}", c.ProductDelete).Methods("DELETE") // admin

	fmt.Println("http://localhost:8001/api/products/")

	serv := &http.Server{
		Addr:    ":8001",
		Handler: r,
	}

	return serv
}

func StartServer(serv *http.Server) {
	serv.ListenAndServe()
}
