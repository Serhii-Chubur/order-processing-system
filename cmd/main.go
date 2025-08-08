package main

import (
	"fmt"
	"log"

	order_app "order_processing_system/order_service/cmd"
	product_app "order_processing_system/product_service/cmd"
	user_app "order_processing_system/user_service/cmd"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	errChan := make(chan error)

	go func() {
		fmt.Println("Hello product")
		if err := product_app.Run(); err != nil {
			errChan <- fmt.Errorf("product service: %w", err)
		}
	}()

	go func() {
		fmt.Println("Hello order")
		if err := order_app.Run(); err != nil {
			errChan <- fmt.Errorf("order service: %w", err)
		}
	}()

	go func() {
		fmt.Println("Hello user")
		if err := user_app.Run(); err != nil {
			errChan <- fmt.Errorf("user service: %w", err)
		}
	}()

	select {}
}
