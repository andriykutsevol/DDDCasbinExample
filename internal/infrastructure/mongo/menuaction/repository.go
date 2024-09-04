package menuaction

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuaction"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"

	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/mongo/storage"
)

func NewRepository(storage *storage.MongoStorage) menuaction.Repository {
	return &repository{
		storage: storage,
	}
}

type repository struct {
	storage *storage.MongoStorage
}

func (r *repository) Query(ctx context.Context, params menuaction.QueryParam) (menuaction.MenuActions, *pagination.Pagination, error) {
	filter := bson.D{{}}
	//pattern := ""

	if v := params.MenuID; v != "" {
		// 	db = db.Where("menu_id=?", v)
		r.storage.AddFilterCondition(&filter, "menu_id", v)
	}

	if v := params.IDs; len(v) > 0 {
		// db = db.Where("id IN (?)", v)
		inFilter := bson.D{{Key: "id", Value: bson.M{"$in": v}}}
		filter = bson.D{{Key: "$and", Value: []bson.D{filter, inFilter}}}

	}

	cursor, err := r.storage.GetCollection().Find(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	var list []*Model
	if err := cursor.All(context.Background(), &list); err != nil {
		return nil, nil, errors.WithStack(errors.New(""))
	}

	return toDomainList(list), nil, nil
}

func (r *repository) Get(ctx context.Context, id string) (*menuaction.MenuAction, error) {
	// db := GetModelDB(ctx, r.db).Where("id=?", id)
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

func (r *repository) GetByIdString(ctx context.Context, id string) (*menuaction.MenuAction, error) {

	return nil, nil
}

func (r *repository) Create(ctx context.Context, item *menuaction.MenuAction) error {
	// result := GetModelDB(ctx, r.db).Create(domainToModel(item))
	// return errors.WithStack(result.Error)

	_, err := r.storage.GetCollection().InsertOne(context.Background(), domainToModel(item))
	return errors.WithStack(err)

}

func (r *repository) Update(ctx context.Context, id string, item *menuaction.MenuAction) error {
	// result := GetModelDB(ctx, r.db).Where("id=?", id).Updates(domainToModel(item))
	// return errors.WithStack(result.Error)

	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	// result := GetModelDB(ctx, r.db).Where("id=?", id).Delete(Model{})
	// return errors.WithStack(result.Error)

	return nil
}

func (r *repository) DeleteByMenuID(ctx context.Context, menuID string) error {
	// result := GetModelDB(ctx, r.db).Where("menu_id=?", menuID).Delete(Model{})
	// return errors.WithStack(result.Error)

	return nil
}

func (r *repository) Purge(ctx context.Context) error {
	_, err := r.storage.GetCollection().DeleteMany(ctx, bson.D{{}})
	return err
}
