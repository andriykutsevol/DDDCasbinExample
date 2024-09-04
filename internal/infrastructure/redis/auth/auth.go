package auth

import (
	"context"
	"fmt"

	//"reflect"
	"time"

	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/redis/storage"
	//redis "github.com/redis/go-redis/v9"
)

type AuthStorage struct {
	storage.RedisStorage
}

func (s *AuthStorage) wrapperKey(key string) string {
	return fmt.Sprintf("%s%s", s.GetKeyPrefix(), key)
}

func (s *AuthStorage) Set(ctx context.Context, tokenString string, expiration time.Duration) error {
	fmt.Println("s.wrapperKey(tokenString)", s.wrapperKey(tokenString))
	dbs := s.GetDatabaseService()
	cmd := dbs.GetClient().Set(ctx, s.wrapperKey(tokenString), "1", expiration)
	return cmd.Err()
}

func (s *AuthStorage) Check(ctx context.Context, tokenString string) (bool, error) {
	// cmd := s.dbs.client.Exists(ctx, s.wrapperKey(tokenString))
	// if err := cmd.Err(); err != nil {
	// 	return false, err
	// }
	// return cmd.Val() > 0, nil

	return true, nil
}

func (s *AuthStorage) Delete(ctx context.Context, tokenString string) (bool, error) {
	// cmd := s.dbs.client.Del(ctx, s.wrapperKey(tokenString))
	// if err := cmd.Err(); err != nil {
	// 	return false, err
	// }
	// return cmd.Val() > 0, nil

	return true, nil
}

func (s *AuthStorage) Close() error {
	// return s.dbs.client.Close()

	return nil
}
