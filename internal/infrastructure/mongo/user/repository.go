package user

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user"

	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/mongo/storage"
)

type repository struct {
	storage *storage.MongoStorage
}

func NewRepository(storage *storage.MongoStorage) user.Repository {
	return &repository{
		storage: storage,
	}
}

func (a *repository) Query(ctx context.Context, params user.QueryParams) (user.Users, *pagination.Pagination, error) {

	filter := bson.D{{}}

	if v := params.UserName; v != "" {
		//db = db.Where("user_name=?", v)
		a.storage.AddFilterCondition(&filter, "username", v)
	}

	if v := params.Status; v > 0 {
		//db = db.Where("status=?", v)
	}
	if v := params.RoleIDs; len(v) > 0 {
		// // todo: serviceへ移動
		// subQuery := userrole.GetModelDB(ctx, a.db).
		// 	Select("user_id").
		// 	Where("role_id IN (?)", v).
		// 	SubQuery()
		// db = db.Where("id IN ?", subQuery)
	}

	if v := params.QueryValue; v != "" {
		// v = "%" + v + "%"
		// db = db.Where("user_name LIKE ? OR real_name LIKE ? OR phone LIKE ? OR email LIKE ?", v, v, v, v)
	}

	cursor, err := a.storage.GetCollection().Find(ctx, filter)
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

//var globalvar string;

func (a *repository) Get(ctx context.Context, id string) (*user.User, error) {

	//globalvar = id

	for i := 0; i < 4; i++ {
		fmt.Println("Iteration user", id)
		time.Sleep(1 * time.Second)
	}

	return nil, nil
}

func (a *repository) Create(ctx context.Context, item *user.User) error {

	_, err := a.storage.GetCollection().InsertOne(ctx, domainToModel(item))

	return errors.WithStack(err)
}

func (a *repository) Update(ctx context.Context, id string, item *user.User) error {

	return errors.WithStack(errors.New(""))
}

func (a *repository) Delete(ctx context.Context, id string) error {

	return errors.WithStack(errors.New(""))
}

func (a *repository) UpdateStatus(ctx context.Context, id string, status int) error {

	return errors.WithStack(errors.New(""))
}

func (a *repository) UpdatePassword(ctx context.Context, id, password string) error {

	return errors.WithStack(errors.New(""))
}

func (a *repository) Purge(ctx context.Context) error {
	_, err := a.storage.GetCollection().DeleteMany(ctx, bson.D{{}})
	return err
}
