package weather

import "context"

type Repository interface {
	Seed(ctx context.Context, params map[string]string) error
	Retrieve(ctx context.Context, city string) (string, error)
}
