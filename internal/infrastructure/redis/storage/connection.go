package storage

import (
	"context"
	"fmt"

	redis "github.com/redis/go-redis/v9"

	"github.com/andriykutsevol/DDDCasbinExample/configs"
)

type DatabaseService struct {
	client *redis.Client
}

func NewDatabaseService() (*DatabaseService, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     configs.C.Redis.Addr,
		Password: configs.C.Redis.Password, // no password set
		DB:       0,                        // use default DB
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		fmt.Println("Failed to Connected to Redis", err)
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	} else {
		fmt.Println("Connected to Redis")
	}

	return &DatabaseService{client: client}, nil
}

func (s *DatabaseService) GetClient() *redis.Client {
	return s.client
}
