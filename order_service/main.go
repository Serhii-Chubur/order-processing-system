package main

import (
	"fmt"
	"log"

	order_app "order_processing_system/order_service/cmd"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	fmt.Println("Hello from order service")
	if err := order_app.Run(); err != nil {
		panic(err)
	}
}
