package main

import (
	"fmt"
	"log"

	user_app "order_processing_system/user_service/cmd"
	// order_app "order_processing_system/order_service/cmd"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	fmt.Println("Hello from user service")

	if err := user_app.Run(); err != nil {
		panic(err)
	}
}
