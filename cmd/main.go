package main

import (
	"fmt"
	"log"

	product_app "order_processing_system/product_service/cmd"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	fmt.Println("Hello")
	if err := product_app.Run(); err != nil {
		panic(err)
	}
}
