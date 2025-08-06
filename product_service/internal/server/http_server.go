package server

import (
	"fmt"
	"net/http"
	"order_processing_system/product_service/internal/controllers"

	"github.com/gorilla/mux"
)

func NewServer(c *controllers.Controller) *http.Server {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from root! Path: %s\n", r.URL.Path)
	})

	productRouter := r.PathPrefix("/api/products").Subrouter()

	productRouter.HandleFunc("/", c.ProductList).Methods("GET")
	productRouter.HandleFunc("/{id}", c.ProductDetail).Methods("GET")
	productRouter.HandleFunc("/{id}/stock", c.ProductStock).Methods("GET")

	// adminRouter := r.PathPrefix("/api/products").Subrouter()
	// adminRouter.Use(middleware.IsAdmin)

	// adminRouter.HandleFunc("/", c.ProductCreate).Methods("POST")        // admin
	// adminRouter.HandleFunc("/{id}", c.ProductUpdate).Methods("PUT")    // admin
	// adminRouter.HandleFunc("/{id}", c.ProductDelete).Methods("DELETE") // admin

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

func StopServer(serv *http.Server) {
	serv.Close()
}
