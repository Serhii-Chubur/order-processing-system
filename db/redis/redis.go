package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type RedisRepo struct {
	Client *redis.Client
}

type Redis interface {
	GetData(id string) (string, error)
	SetCache(cacheKey string, data []byte) error
	DeleteProduct(cacheKey string) error
	SetAccessToken(email string, accessToken string) error
	SetRefreshToken(email string, refreshToken string) error
	GetUserEmail(accessToken string) (string, error)
}

func NewRedisRepo(client *redis.Client) *RedisRepo {
	return &RedisRepo{
		Client: client,
	}
}

func ConnectRedis(config RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Redis")

	return client
}

func (r *RedisRepo) GetData(id string) (string, error) {
	return r.Client.Get(context.Background(), id).Result()
}

func (r *RedisRepo) SetCache(cacheKey string, data []byte) error {
	return r.Client.Set(context.Background(), cacheKey, data, 3*time.Minute).Err()
}

func (r *RedisRepo) DeleteProduct(cacheKey string) error {
	return r.Client.Del(context.Background(), cacheKey).Err()
}

func (r *RedisRepo) SetAccessToken(email string, accessToken string) error {
	return r.Client.Set(context.Background(), accessToken, email, 3*time.Minute).Err()
}

func (r *RedisRepo) SetRefreshToken(email string, refreshToken string) error {
	return r.Client.Set(context.Background(), refreshToken, email, 3*time.Hour).Err()
}

func (r *RedisRepo) GetUserEmail(accessToken string) (string, error) {
	return r.Client.Get(context.Background(), accessToken).Result()
}
