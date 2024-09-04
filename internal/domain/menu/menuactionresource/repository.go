package menuactionresource

import (
	"context"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
)

type Repository interface {
	Query(ctx context.Context, params QueryParam) (MenuActionResources, *pagination.Pagination, error)
	Get(ctx context.Context, id string) (*MenuActionResource, error)
	Create(ctx context.Context, item *MenuActionResource) error
	Update(ctx context.Context, id string, item *MenuActionResource) error
	Delete(ctx context.Context, id string) error
	DeleteByActionID(ctx context.Context, actionID string) error
	DeleteByMenuID(ctx context.Context, menuID string) error
	Purge(ctx context.Context) error
}
