package userrole

import (
	"context"
	"log"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/userrole"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/storage"
)

func NewRepository(service storage.DatabaseService) userrole.Repository {
	return &repository{
		service: service,
	}
}

type repository struct {
	service storage.DatabaseService
}

func (r *repository) Query(ctx context.Context, params userrole.QueryParam) (userrole.UserRoles, *pagination.Pagination, error) {

	rows, err := r.service.Pool.Query(context.Background(), "SELECT * FROM casbin.userrole")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var userroles []*Model

	for rows.Next() {
		var user Model
		err := rows.Scan(
			&user.ID,
			&user.UserID,
			&user.RoleID,
			&user.IDString,
			&user.UserIDString,
			&user.RoleIDString,
		)
		if err != nil {
			return nil, nil, err
		}
		userroles = append(userroles, &user)
	}

	return toDomainList(userroles), nil, nil
}

func (r *repository) Get(ctx context.Context, id string) (*userrole.UserRole, error) {
	log.Println("user/userrole/repository.go: Get()")
	return nil, nil
}

func (r *repository) Create(ctx context.Context, item *userrole.UserRole) error {

	model := domainToModel(item)
	queryString := `
		INSERT INTO casbin.userrole (
		id, 
		userid,
		roleid,
		idstring,
		useridstring,
		roleidstring
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.service.Pool.Exec(context.Background(), queryString,
		model.ID, model.UserID, model.RoleID,
		model.IDString, model.UserIDString, model.RoleIDString)
	if err != nil {
		log.Fatalf("Failed to insert data: %v\n", err)
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string, item *userrole.UserRole) error {
	log.Println("user/userrole/repository.go: Update()")
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	log.Println("user/userrole/repository.go: Delete()")
	return nil
}

func (r *repository) DeleteByUserID(ctx context.Context, userID string) error {
	log.Println("user/userrole/repository.go: DeleteByUserID()")
	return nil
}

func (r *repository) Purge(ctx context.Context) error {
	_, err := r.service.Pool.Exec(context.Background(), "DELETE FROM "+"casbin.userrole")
	if err != nil {
		log.Printf("Failed to delete from table %s: %v\n", "userrole", err)
	}

	return nil
}
