package redis

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (r *RedisRepo) Delete(key string) error {
	return r.Client.Del(context.Background(), key).Err()
}

func (r *RedisRepo) SetAccessToken(email string, accessToken string) error {
	return r.Client.Set(context.Background(), accessToken, email, 3*time.Minute).Err()
}

func (r *RedisRepo) SetRefreshToken(email string, refreshToken string) error {
	return r.Client.Set(context.Background(), refreshToken, email, 3*time.Hour).Err()
}

func (r *RedisRepo) GetUserEmail(token string) (string, error) {
	return r.Client.Get(context.Background(), token).Result()
}

type Claims struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
	Root  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func ParseToken(tokenStr string) (*Claims, error) {
	fmt.Println("JWT_SECRET:", os.Getenv("JWT_SECRET"))

	fmt.Println(tokenStr)
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	fmt.Println(token.Claims)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	fmt.Println(claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
