package menu

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"

	mongostorage "github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/mongo/storage"
)

func NewRepository(storage *mongostorage.MongoStorage) menu.Repository {
	return &repository{
		storage: storage,
	}
}

type repository struct {
	storage *mongostorage.MongoStorage
}

func (a *repository) Query(ctx context.Context, params menu.QueryParam) (menu.Menus, *pagination.Pagination, error) {

	filter := bson.D{{}}
	pattern := ""

	if v := params.IDs; len(v) > 0 {
		//db = db.Where("id IN (?)", v)
		a.storage.AddFilterCondition(&filter, "id", v)
	}
	if v := params.Name; v != "" {
		//db = db.Where("name=?", v)
		a.storage.AddFilterCondition(&filter, "name", v)
	}
	if v := params.ParentID; v != "" {
		//db = db.Where("parent_id=?", *v)
		a.storage.AddFilterCondition(&filter, "parent_id", v)
	}
	if v := params.PrefixParentPath; v != "" {
		//db = db.Where("parent_path LIKE ?", v+"%")

		// Define the regex pattern
		pattern = v + ".*"
		a.storage.AddFilterCondition(&filter, "parent_path", primitive.Regex{Pattern: pattern, Options: ""})
	}
	if v := params.ShowStatus; v != 0 {
		//db = db.Where("show_status=?", v)
		a.storage.AddFilterCondition(&filter, "show_status", v)
	}
	if v := params.Status; v != 0 {
		//db = db.Where("status=?", v)
		a.storage.AddFilterCondition(&filter, "status", v)
	}

	if v := params.QueryValue; v != "" {

		regex := ".*" + v + ".*"
		sub_pattern := primitive.Regex{Pattern: regex, Options: ""}
		sub_filter := bson.D{
			{Key: "$or", Value: bson.A{
				bson.M{"name": sub_pattern},
				bson.M{"memo": sub_pattern},
			}},
		}

		filter = append(filter, sub_filter...)

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

func (a *repository) Get(ctx context.Context, id string) (*menu.Menu, error) {
	item := &Model{}

	// ok, err := gormx.FindOne(ctx, GetModelDB(ctx, a.db).Where("id=?", id), &item)
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }
	// if !ok {
	// 	return nil, nil
	// }

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: docID}}

	err = a.storage.GetCollection().FindOne(context.TODO(), filter).Decode(item)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return item.ToDomain(), nil
		}
		return nil, errors.New("Error finding document")
	}

	return item.ToDomain(), nil

}

func (a *repository) GetByIdString(ctx context.Context, id string) (*menu.Menu, error) {

	return nil, nil
}

func (a *repository) Create(ctx context.Context, item *menu.Menu) error {

	// result := GetModelDB(ctx, a.db).Create(domainToModel(item))
	// return errors.WithStack(result.Error)

	_, err := a.storage.GetCollection().InsertOne(context.Background(), domainToModel(item))
	return errors.WithStack(err)
}

func (a *repository) Update(ctx context.Context, id string, item *menu.Menu) error {
	// result := GetModelDB(ctx, a.db).Where("id=?", id).Updates(domainToModel(item))
	// return errors.WithStack(result.Error)

	return nil
}

func (a *repository) UpdateParentPath(ctx context.Context, id, parentPath string) error {
	// result := GetModelDB(ctx, a.db).Where("id=?", id).Update("parent_path", parentPath)
	// return errors.WithStack(result.Error)

	return nil
}

func (a *repository) Delete(ctx context.Context, id string) error {
	// result := GetModelDB(ctx, a.db).Where("id=?", id).Delete(Model{})
	// return errors.WithStack(result.Error)

	return nil
}

func (a *repository) UpdateStatus(ctx context.Context, id string, status int) error {
	// result := GetModelDB(ctx, a.db).Where("id=?", id).Update("status", status)
	// return errors.WithStack(result.Error)

	return nil
}

func (a *repository) Purge(ctx context.Context) error {
	_, err := a.storage.GetCollection().DeleteMany(ctx, bson.D{{}})
	return err
}
