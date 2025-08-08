package cmd

import (
	"fmt"
	"log"
	"order_processing_system/db/psql"
	"order_processing_system/db/redis"
	"order_processing_system/order_service/internal/controllers"
	"order_processing_system/order_service/internal/server"
	"order_processing_system/order_service/internal/services"
	"os"
	"os/signal"
	"strconv"

	"syscall"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("./db/configs/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Run() error {
	errChan := make(chan error, 1)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	psqlConfig := psql.PSQLConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
	}

	redis_db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal(err)
	}
	redisConfig := redis.RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redis_db,
	}

	psqlConn := psql.ConnectPSQL(psqlConfig)
	redisConn := redis.ConnectRedis(redisConfig)

	// repositories
	psqlRepo := psql.NewPSQLRepo(psqlConn)
	redisRepo := redis.NewRedisRepo(redisConn)

	// service
	productService := services.NewService(psqlRepo, redisRepo)

	// controller
	orderController := controllers.NewController(errChan, productService)

	// server
	orderSrv := server.NewServer(orderController)

	go func() {
		select {
		case sig := <-stopChan:
			log.Print(sig)
		case err := <-errChan:
			log.Print(err)
		}

		err := psqlConn.Close()
		if err != nil {
			log.Print(err)
		} else {
			fmt.Println("Product PSQL connection closed")
		}

		err = redisConn.Close()
		if err != nil {
			log.Print(err)
		} else {
			fmt.Println("Product Redis connection closed")
		}

		server.StopServer(orderSrv)
		os.Exit(0)
	}()

	server.StartServer(orderSrv)
	return nil
}
