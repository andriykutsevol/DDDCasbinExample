package userrole

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/userrole"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/mongo/storage"
)

func NewRepository(storage *storage.MongoStorage) userrole.Repository {
	return &repository{
		storage: storage,
	}
}

type repository struct {
	storage *storage.MongoStorage
}

func (a *repository) Query(ctx context.Context, params userrole.QueryParam) (userrole.UserRoles, *pagination.Pagination, error) {

	filter := bson.D{{}}
	//pattern := ""

	if v := params.UserID; v != "" {
		//db = db.Where("user_id=?", v)
		a.storage.AddFilterCondition(&filter, "user_id", v)
	}

	if v := params.UserIDs; len(v) > 0 {
		//db = db.Where("user_id IN (?)", v)
		inFilter := bson.D{{Key: "user_id", Value: bson.M{"$in": v}}}
		filter = bson.D{{Key: "$and", Value: []bson.D{filter, inFilter}}}
	}

	cursor, err := a.storage.GetCollection().Find(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	var list []*Model
	if err := cursor.All(context.Background(), &list); err != nil {
		return nil, nil, errors.WithStack(errors.New(""))
	}

	return toDomainList(list), nil, nil

}

var globalvar string

func (a *repository) Get(ctx context.Context, id string) (*userrole.UserRole, error) {
	// db := GetModelDB(ctx, a.db).Where("id=?", id)
	// item := &Model{}
	// ok, err := gormx.FindOne(ctx, db, &item)
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }
	// if !ok {
	// 	return nil, nil
	// }

	// return item.ToDomain(), nil

	globalvar = id

	for i := 0; i < 4; i++ {
		fmt.Println("Iteration userrole", globalvar)
		time.Sleep(1 * time.Second)
	}

	return nil, nil

}

func (a *repository) Create(ctx context.Context, item *userrole.UserRole) error {
	// result := GetModelDB(ctx, a.db).Create(domainToModel(item))
	// return errors.WithStack(result.Error)

	_, err := a.storage.GetCollection().InsertOne(context.Background(), domainToModel(item))

	return errors.WithStack(err)
}

func (a *repository) Update(ctx context.Context, id string, item *userrole.UserRole) error {
	// result := GetModelDB(ctx, a.db).Where("id=?", id).Updates(domainToModel(item))
	// return errors.WithStack(result.Error)

	return nil
}

func (a *repository) Delete(ctx context.Context, id string) error {
	// result := GetModelDB(ctx, a.db).Where("id=?", id).Delete(Model{})
	// return errors.WithStack(result.Error)

	return nil
}

func (a *repository) DeleteByUserID(ctx context.Context, userID string) error {
	// result := GetModelDB(ctx, a.db).Where("user_id=?", userID).Delete(Model{})
	// return errors.WithStack(result.Error)

	return nil
}

func (a *repository) Purge(ctx context.Context) error {
	_, err := a.storage.GetCollection().DeleteMany(ctx, bson.D{{}})
	return err
}
