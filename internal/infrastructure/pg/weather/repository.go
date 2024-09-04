package weather

import (
	"context"
	"log"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/weather"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/storage"
)

func NewRepository(service storage.DatabaseService) weather.Repository {
	return &repository{
		service: service,
	}
}

type repository struct {
	service storage.DatabaseService
}

func (r *repository) Seed(ctx context.Context, params map[string]string) error {
	log.Println("weather/repository.go: Seed()")
	return nil
}

func (r *repository) Retrieve(ctx context.Context, city string) (string, error) {
	log.Println("weather/repository.go: Retrieve()")
	return "{}", nil
}
