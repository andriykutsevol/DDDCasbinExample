package rolemenu

import (
	"context"
	"log"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/rolemenu"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/storage"
)

func NewRepository(service storage.DatabaseService) rolemenu.Repository {
	return &repository{
		service: service,
	}
}

type repository struct {
	service storage.DatabaseService
}

func (r *repository) Query(ctx context.Context, params rolemenu.QueryParam) (rolemenu.RoleMenus, *pagination.Pagination, error) {
	rows, err := r.service.Pool.Query(context.Background(), "SELECT * FROM casbin.rolemenu")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var rolemenus []*Model
	for rows.Next() {
		var rolemenu Model
		err := rows.Scan(
			&rolemenu.ID,
			&rolemenu.RoleID,
			&rolemenu.MenuID,
			&rolemenu.ActionID,
			&rolemenu.IDString,
			&rolemenu.RoleIDString,
			&rolemenu.MenuIDString,
			&rolemenu.ActionIDString)
		if err != nil {
			return nil, nil, err
		}
		rolemenus = append(rolemenus, &rolemenu)
	}
	return toDomainList(rolemenus), nil, nil
}

func (r *repository) Get(ctx context.Context, id string) (*rolemenu.RoleMenu, error) {
	log.Println("user/rolemenu/repository.go: Get()")
	return nil, nil
}

func (r *repository) Create(ctx context.Context, item *rolemenu.RoleMenu) error {
	model := domainToModel(item)

	queryString := `
		INSERT INTO casbin.rolemenu (
		id, 
		roleid,
		menuid,
		actionid,
		idstring,
		roleidstring,
		menuidstring,
		actionidstring
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.service.Pool.Exec(context.Background(), queryString,
		model.ID, model.RoleID, model.MenuID, model.ActionID,
		model.IDString, model.RoleIDString, model.MenuIDString, model.ActionIDString)
	if err != nil {
		log.Fatalf("Failed to insert data: %v\n", err)
	}
	return nil

}

func (r *repository) Update(ctx context.Context, id string, item *rolemenu.RoleMenu) error {
	log.Println("user/rolemenu/repository.go: Update()")
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	log.Println("user/rolemenu/repository.go: Delete()")
	return nil
}

func (r *repository) DeleteByRoleID(ctx context.Context, roleID string) error {
	log.Println("user/rolemenu/repository.go: DeleteByRoleID()")
	return nil
}

func (r *repository) Purge(ctx context.Context) error {
	_, err := r.service.Pool.Exec(context.Background(), "DELETE FROM casbin.rolemenu")
	if err != nil {
		log.Printf("Failed to delete from table %s: %v\n", "rolemenu", err)
	}

	return nil
}
