package rolemenu

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/rolemenu"

	storage "github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/mongo/storage"
)

type repository struct {
	storage *storage.MongoStorage
}

func NewRepository(storage *storage.MongoStorage) rolemenu.Repository {
	return &repository{
		storage: storage,
	}
}

func (a *repository) Query(ctx context.Context, params rolemenu.QueryParam) (rolemenu.RoleMenus, *pagination.Pagination, error) {
	// db := GetModelDB(ctx, a.db)
	// if v := params.RoleID; v != "" {
	// 	db = db.Where("role_id=?", v)
	// }
	// if v := params.RoleIDs; len(v) > 0 {
	// 	db = db.Where("role_id IN (?)", v)
	// }
	// db = db.Order(gormx.ParseOrder(params.OrderFields.AddIdSortField()))

	filter := bson.D{{}}
	//pattern := ""

	if v := params.RoleID; v != "" {
		a.storage.AddFilterCondition(&filter, "role_id", v)
	}

	if v := params.RoleIDs; len(v) > 0 {

		inFilter := bson.D{{Key: "role_id", Value: bson.M{"$in": v}}}
		filter = bson.D{{Key: "$and", Value: []bson.D{filter, inFilter}}}

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

func (a *repository) Get(ctx context.Context, id string) (*rolemenu.RoleMenu, error) {
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

	return nil, nil
}

func (a *repository) Create(ctx context.Context, item *rolemenu.RoleMenu) error {
	// result := GetModelDB(ctx, a.db).Create(domainToModel(item))
	// return errors.WithStack(result.Error)
	_, err := a.storage.GetCollection().InsertOne(context.Background(), domainToModel(item))
	return errors.WithStack(err)
}

func (a *repository) Update(ctx context.Context, id string, item *rolemenu.RoleMenu) error {
	// result := GetModelDB(ctx, a.db).Where("id=?", id).Updates(domainToModel(item))
	// return errors.WithStack(result.Error)
	return errors.WithStack(errors.New(""))
}

func (a *repository) Delete(ctx context.Context, id string) error {
	// result := GetModelDB(ctx, a.db).Where("id=?", id).Delete(Model{})
	// return errors.WithStack(result.Error)
	return errors.WithStack(errors.New(""))
}

func (a *repository) DeleteByRoleID(ctx context.Context, roleID string) error {
	// result := GetModelDB(ctx, a.db).Where("role_id=?", roleID).Delete(Model{})
	// return errors.WithStack(result.Error)
	return errors.WithStack(errors.New(""))
}

func (a *repository) Purge(ctx context.Context) error {
	_, err := a.storage.GetCollection().DeleteMany(ctx, bson.D{{}})
	return err
}
