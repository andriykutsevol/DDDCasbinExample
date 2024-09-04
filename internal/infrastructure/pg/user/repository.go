package user

import (
	"context"
	"log"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/storage"
)

func NewRepository(service storage.DatabaseService) user.Repository {
	return &repository{
		service: service,
	}
}

type repository struct {
	service storage.DatabaseService
}

func (r *repository) Query(ctx context.Context, params user.QueryParams) (user.Users, *pagination.Pagination, error) {
	rows, err := r.service.Pool.Query(context.Background(), "SELECT * FROM casbin.user")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var users []*Model
	for rows.Next() {
		var user Model
		err := rows.Scan(
			&user.ID,
			&user.UserName,
			&user.RealName,
			&user.Password,
			&user.Email,
			&user.Phone,
			&user.Status,
			&user.Creator,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.IDString)
		if err != nil {
			return nil, nil, err
		}
		users = append(users, &user)
	}
	return toDomainList(users), nil, nil
}

func (r *repository) Get(ctx context.Context, id string) (*user.User, error) {
	log.Println("user/repository.go: Get()")
	return nil, nil
}

func (r *repository) Create(ctx context.Context, item *user.User) error {
	model := domainToModel(item)

	queryString := `
		INSERT INTO casbin.user (
		id, 
		username,
		realname,
		password,
		email,
		phone,
		status,
		creator,
		createdat,
		updatedat,
		deletedat,
		idstring
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.service.Pool.Exec(context.Background(), queryString,
		model.ID,
		model.UserName,
		model.RealName,
		model.Password, model.Email, model.Phone, model.Status, model.Creator,
		model.CreatedAt, model.UpdatedAt, model.DeletedAt, model.IDString)
	if err != nil {
		log.Fatalf("Failed to insert data: %v\n", err)
	}
	return nil

}

func (r *repository) Update(ctx context.Context, id string, item *user.User) error {
	log.Println("user/repository.go: Update()")
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	log.Println("user/repository.go: Delete()")
	return nil
}

func (r *repository) UpdateStatus(ctx context.Context, id string, status int) error {
	log.Println("user/repository.go: UpdateStatu()")
	return nil
}

func (r *repository) UpdatePassword(ctx context.Context, id, password string) error {
	log.Println("user/repository.go: UpdatePasswor()")
	return nil
}

func (r *repository) Purge(ctx context.Context) error {

	_, err := r.service.Pool.Exec(context.Background(), "DELETE FROM "+"casbin.user")
	if err != nil {
		log.Printf("Failed to delete from table %s: %v\n", "user", err)
	}

	return nil
}
