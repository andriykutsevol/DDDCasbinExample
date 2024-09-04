package role

import (
	"context"
	"log"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/role"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/storage"
)

func NewRepository(service storage.DatabaseService) role.Repository {
	return &repository{
		service: service,
	}
}

type repository struct {
	service storage.DatabaseService
}

func (r *repository) Query(ctx context.Context, params role.QueryParam) (role.Roles, *pagination.Pagination, error) {
	// if v := params.IDs; len(v) > 0 {
	// 	// db = db.Where("id IN (?)", v)
	// }
	// if v := params.Name; v != "" {
	// 	//db = db.Where("name=?", v)
	// }
	// if v := params.UserID; v != "" {
	// 	// todo: serviceへ移動
	// 	// subQuery := userrole.GetModelDB(ctx, a.db).
	// 	// 	Where("deleted_at is null").
	// 	// 	Where("user_id=?", v).
	// 	// 	Select("role_id").SubQuery()
	// 	// db = db.Where("id IN ?", subQuery)
	// }
	// if v := params.QueryValue; v != "" {
	// 	// v = "%" + v + "%"
	// 	// db = db.Where("name LIKE ? OR memo LIKE ?", v, v)
	// }

	rows, err := r.service.Pool.Query(context.Background(), "SELECT * FROM casbin.role")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var roles []*Model

	for rows.Next() {
		var role Model
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Sequence,
			&role.Memo,
			&role.Status,
			&role.Creator,
			&role.CreatedAt,
			&role.UpdatedAt,
			&role.DeletedAt,
			&role.IDString)
		if err != nil {
			log.Println("QWER: ", role.ID)
			return nil, nil, err
		}
		roles = append(roles, &role)
	}
	return toDomainList(roles), nil, nil
}

func (r *repository) Get(ctx context.Context, id string) (*role.Role, error) {
	log.Println("user/role/repository.go: Get()")
	return nil, nil
}

func (r *repository) GetByIdString(ctx context.Context, id string) (*role.Role, error) {
	var role Model
	queryString := "SELECT * FROM casbin.role WHERE idstring = $1"
	err := r.service.Pool.QueryRow(context.Background(), queryString, id).Scan(
		&role.ID,
		&role.Name,
		&role.Sequence,
		&role.Memo,
		&role.Status,
		&role.Creator,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
		&role.IDString,
	)
	if err != nil {
		return nil, err
	}
	return role.ToDomain(), nil
}

func (r *repository) Create(ctx context.Context, item *role.Role) error {
	model := domainToModel(item)
	queryString := `
		INSERT INTO casbin.role (
		id, 
		name, 
		sequence, 
		memo,
		status,
		creator,
		createdat,
		updatedat,
		deletedat,
		idstring
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.service.Pool.Exec(context.Background(), queryString,
		model.ID,
		model.Name,
		model.Sequence,
		model.Memo,
		model.Status,
		model.Creator,
		model.CreatedAt,
		model.UpdatedAt,
		model.DeletedAt,
		model.IDString)
	if err != nil {
		log.Fatalf("Failed to insert data: %v\n", err)
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string, item *role.Role) error {
	log.Println("user/role/repository.go: Update()")
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	log.Println("user/role/repository.go: Delete()")
	return nil
}

func (r *repository) UpdateStatus(ctx context.Context, id string, status int) error {
	log.Println("user/role/repository.go: UpdateStatus()")
	return nil
}

func (r *repository) Purge(ctx context.Context) error {

	_, err := r.service.Pool.Exec(context.Background(), "DELETE FROM "+"casbin.role")
	if err != nil {
		log.Printf("Failed to delete from table %s: %v\n", "role", err)
	}

	return nil
}
