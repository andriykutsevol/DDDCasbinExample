package role

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/role"

	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/mongo/storage"
)

func NewRepository(storage *storage.MongoStorage) role.Repository {
	return &repository{
		storage: storage,
	}
}

type repository struct {
	storage *storage.MongoStorage
}

func (a *repository) Create(ctx context.Context, item *role.Role) error {

	_, err := a.storage.GetCollection().InsertOne(context.Background(), domainToModel(item))

	return errors.WithStack(err)
}

func (a *repository) Query(ctx context.Context, params role.QueryParam) (role.Roles, *pagination.Pagination, error) {

	filter := bson.D{{}}

	// TODO. Implement queries

	if v := params.IDs; len(v) > 0 {
		//db = db.Where("id IN (?)", v)
	}

	if v := params.Name; v != "" {
		//db = db.Where("name=?", v)
	}

	if v := params.UserID; v != "" {

	}

	if v := params.QueryValue; v != "" {

	}

	cursor, err := a.storage.GetCollection().Find(ctx, filter)
	_ = cursor
	if err != nil {
		fmt.Println("Error finding documents:", err)
		return nil, nil, err
	}

	var list []*Model
	if err := cursor.All(context.Background(), &list); err != nil {
		fmt.Println("Error decoding documents:", err)
		return nil, nil, errors.WithStack(errors.New(""))
	}

	return toDomainList(list), nil, nil
}

func (a *repository) Get(ctx context.Context, id string) (*role.Role, error) {

	fmt.Println("roleMongoRepository.go: Get()")

	return nil, nil
}

func (r *repository) GetByIdString(ctx context.Context, id string) (*role.Role, error) {
	fmt.Println("roleMongoRepository.go: GetByIdString()")
	return nil, nil
}

func (a *repository) Update(ctx context.Context, id string, item *role.Role) error {

	fmt.Println("roleMongoRepository.go: Update()")

	return nil
}

func (a *repository) Delete(ctx context.Context, id string) error {

	fmt.Println("roleMongoRepository.go: Delete()")

	return nil
}

func (a *repository) UpdateStatus(ctx context.Context, id string, status int) error {

	fmt.Println("roleMongoRepository.go: UpdateStatus()")

	return nil
}

func (a *repository) Purge(ctx context.Context) error {
	_, err := a.storage.GetCollection().DeleteMany(ctx, bson.D{{}})
	return err
}
