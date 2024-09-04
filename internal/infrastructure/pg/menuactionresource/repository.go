package menuactionresource

import (
	"context"
	"log"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuaction"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuactionresource"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/storage"
)

func NewRepository(service storage.DatabaseService, menuactionsRepo menuaction.Repository) menuactionresource.Repository {
	return &repository{
		service: service,
	}
}

type repository struct {
	service storage.DatabaseService
}

func (r *repository) Query(ctx context.Context, params menuactionresource.QueryParam) (menuactionresource.MenuActionResources, *pagination.Pagination, error) {
	rows, err := r.service.Pool.Query(context.Background(), "SELECT * FROM casbin.menuactionresource")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var menus []*Model

	for rows.Next() {
		var menu Model
		err := rows.Scan(
			&menu.ID,
			&menu.ActionID,
			&menu.Method,
			&menu.Path,
			&menu.IDString,
			&menu.ActionIDString)
		if err != nil {
			return nil, nil, err
		}
		menus = append(menus, &menu)
	}
	return toDomainList(menus), nil, nil
}

func (r *repository) Get(ctx context.Context, id string) (*menuactionresource.MenuActionResource, error) {
	log.Println("menu/menuactionresource/repository.go: Get()")
	return nil, nil
}

func (r *repository) Create(ctx context.Context, item *menuactionresource.MenuActionResource) error {
	model := domainToModel(item)
	queryString := `
		INSERT INTO casbin.menuactionresource (
		id, 
		actionid,
		method,
		path,
		idstring,
		actionidstring
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.service.Pool.Exec(context.Background(), queryString,
		model.ID,
		model.ActionID,
		model.Method,
		model.Path,
		model.IDString,
		model.ActionIDString)
	if err != nil {
		log.Fatalf("Failed to insert data: %v\n", err)
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string, item *menuactionresource.MenuActionResource) error {
	log.Println("menu/menuactionresource/repository.go: Update()")
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	log.Println("menu/menuactionresource/repository.go: Delete()")
	return nil
}

func (r *repository) DeleteByActionID(ctx context.Context, actionID string) error {
	log.Println("menu/menuactionresource/repository.go: DeleteByActionID()")
	return nil
}

func (r *repository) DeleteByMenuID(ctx context.Context, menuID string) error {
	log.Println("menu/menuactionresource/repository.go: DeleteByMenuID()")
	return nil
}

func (r *repository) Purge(ctx context.Context) error {

	_, err := r.service.Pool.Exec(context.Background(), "DELETE FROM casbin.menuactionresource")
	if err != nil {
		log.Printf("Failed to delete from table %s: %v\n", "menuactionresource", err)
	}

	return nil
}
