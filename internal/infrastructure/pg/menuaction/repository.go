package menuaction

import (
	"context"
	"log"

	//"github.com/google/uuid"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuaction"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/storage"
)

func NewRepository(service storage.DatabaseService) menuaction.Repository {
	return &repository{
		service: service,
	}
}

type repository struct {
	service storage.DatabaseService
}

func (r *repository) Query(ctx context.Context, params menuaction.QueryParam) (menuaction.MenuActions, *pagination.Pagination, error) {
	log.Println("menu/menuaction/repository.go: Query()")

	rows, err := r.service.Pool.Query(context.Background(), "SELECT * FROM casbin.menuaction")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var menuactions []*Model
	for rows.Next() {
		var menuaction Model
		err := rows.Scan(
			&menuaction.ID,
			&menuaction.MenuID,
			&menuaction.Code,
			&menuaction.Name,
		)
		if err != nil {
			return nil, nil, err
		}
		menuactions = append(menuactions, &menuaction)
	}
	return toDomainList(menuactions), nil, nil
}

func (r *repository) Get(ctx context.Context, id string) (*menuaction.MenuAction, error) {
	log.Println("menu/menuaction/repository.go: Get()")
	return nil, nil
}

func (r *repository) GetByIdString(ctx context.Context, id string) (*menuaction.MenuAction, error) {

	var menuaction Model
	queryString := "SELECT * FROM casbin.menuaction WHERE idstring = $1"
	err := r.service.Pool.QueryRow(context.Background(), queryString, id).Scan(
		&menuaction.ID,
		&menuaction.MenuID,
		&menuaction.Code,
		&menuaction.Name,
		&menuaction.IDString,
		&menuaction.MenuIDString)
	if err != nil {
		return nil, err
	}
	return menuaction.ToDomain(), nil
}

func (r *repository) Create(ctx context.Context, item *menuaction.MenuAction) error {
	model := domainToModel(item)

	query := `
    INSERT INTO casbin.menuaction (
        id, menuid, code, name, idstring, menuidstring
    ) VALUES (
        $1, $2, $3, $4, $5, $6
    )`

	_, err := r.service.Pool.Exec(context.Background(), query,
		model.ID,
		model.MenuID,
		model.Code,
		model.Name,
		model.IDString,
		model.MenuIDString)
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string, item *menuaction.MenuAction) error {
	log.Println("menu/menuaction/repository.go: Update()")
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	log.Println("menu/menuaction/repository.go: Delete()")
	return nil
}

func (r *repository) DeleteByMenuID(ctx context.Context, menuID string) error {
	log.Println("menu/menuaction/repository.go: DeleteByMenuID()")
	return nil
}

func (r *repository) Purge(ctx context.Context) error {

	_, err := r.service.Pool.Exec(context.Background(), "DELETE FROM casbin.menuaction")
	if err != nil {
		log.Printf("Failed to delete from table %s: %v\n", "menuaction", err)
	}

	return nil
}
