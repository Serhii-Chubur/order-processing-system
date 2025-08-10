package main

import (
	"log"
	db "order_processing_system/db"
)

func main() {
	if err := db.RunMigration(); err != nil {
		log.Fatal(err)
	}
}
