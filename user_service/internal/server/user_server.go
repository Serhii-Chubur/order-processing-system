package server

import (
	"fmt"
	"net/http"
	"order_processing_system/user_service/internal/controllers"

	"github.com/gorilla/mux"
)

func NewServer(c *controllers.Controller) *http.Server {
	r := mux.NewRouter().StrictSlash(true)

	userRouter := r.PathPrefix("/api/users").Subrouter()

	userRouter.HandleFunc("/register", c.RegisterUser).Methods("POST")
	userRouter.HandleFunc("/login", c.Login).Methods("POST")
	userRouter.HandleFunc("/{id}", c.GetUserProfile).Methods("GET")
	userRouter.HandleFunc("/{id}", c.UpdateUserProfile).Methods("PUT")
	userRouter.HandleFunc("/logout", c.Logout).Methods("POST")

	fmt.Println("http://localhost:8003/api/users/")

	serv := &http.Server{
		Addr:    ":8003",
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
