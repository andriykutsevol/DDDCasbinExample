package menuactionresource

import (
	"context"
	//"fmt"

	//"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuaction"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuactionresource"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"

	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/mongo/storage"
)

func NewRepository(storage *storage.MongoStorage, menuactionsRepo menuaction.Repository) menuactionresource.Repository {
	return &repository{
		storage:         storage,
		menuactionsRepo: menuactionsRepo,
	}
}

type repository struct {
	storage         *storage.MongoStorage
	menuactionsRepo menuaction.Repository
}

func (a *repository) Query(ctx context.Context, params menuactionresource.QueryParam) (menuactionresource.MenuActionResources, *pagination.Pagination, error) {

	filter := bson.D{{}}
	//pattern := ""

	//fmt.Println("menuactionresource Qery()")

	if v := params.MenuID; v != "" {

		menuactionQP := new(menuaction.QueryParam)
		menuactionQP.MenuID = params.MenuID

		menuactions, _, _ := a.menuactionsRepo.Query(ctx, *menuactionQP)
		_ = menuactions
		//fmt.Println("menuactions: ", menuactions)

		// TODO. Complete Query implementation

		// subQuery := menuaction.GetModelDB(ctx, a.db).Where("menu_id=?", v).Select("id").SubQuery()
		// db = db.Where("action_id IN ?", subQuery)

		// subfilter := bson.D{{}}

		// a.storage.AddFilterCondition(&filter, "menu_id", v)

		// cursor, err := a.menuactionstorage.GetCollection().Find(ctx, subfilter)
		// if err != nil {
		// 	fmt.Println("Error finding documents:", err)
		// 	return nil, nil, err
		// }

		// var menuactionlist []*menuaction.Model
		// if err := cursor.All(context.Background(), &menuactionlist); err != nil {
		// 	fmt.Println("Error decoding documents:", err)
		// 	return nil, nil, errors.WithStack(errors.New(""))
		// }

		//subQuery :=

	}

	// if v := params.MenuIDs; len(v) > 0 {
	// 	subQuery := menuaction.GetModelDB(ctx, a.db).Where("menu_id IN (?)", v).Select("id").SubQuery()
	// 	db = db.Where("action_id IN ?", subQuery)
	// }

	// db = db.Order(gormx.ParseOrder(params.OrderFields.AddIdSortField()))

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

func (a *repository) Get(ctx context.Context, id string) (*menuactionresource.MenuActionResource, error) {
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

func (a *repository) Create(ctx context.Context, item *menuactionresource.MenuActionResource) error {
	// result := GetModelDB(ctx, a.db).Create(domainToModel(item))
	// return errors.WithStack(result.Error)

	_, err := a.storage.GetCollection().InsertOne(context.Background(), domainToModel(item))
	return errors.WithStack(err)
}

func (a *repository) Update(ctx context.Context, id string, item *menuactionresource.MenuActionResource) error {
	// result := GetModelDB(ctx, a.db).Where("id=?", id).Updates(domainToModel(item))
	// return errors.WithStack(result.Error)

	_, err := a.storage.GetCollection().InsertOne(context.Background(), domainToModel(item))
	return errors.WithStack(err)

}

func (a *repository) Delete(ctx context.Context, id string) error {
	// result := GetModelDB(ctx, a.db).Where("id=?", id).Delete(Model{})
	// return errors.WithStack(result.Error)

	return nil
}

func (a *repository) DeleteByActionID(ctx context.Context, actionID string) error {
	// result := GetModelDB(ctx, a.db).Where("action_id =?", actionID).Delete(Model{})
	// return errors.WithStack(result.Error)

	return nil
}

func (a *repository) DeleteByMenuID(ctx context.Context, menuID string) error {
	// subQuery := menuaction.GetModelDB(ctx, a.db).Where("menu_id=?", menuID).Select("id").SubQuery()
	// result := GetModelDB(ctx, a.db).Where("action_id IN ?", subQuery).Delete(Model{})
	// return errors.WithStack(result.Error)

	return nil
}

func (a *repository) Purge(ctx context.Context) error {
	_, err := a.storage.GetCollection().DeleteMany(ctx, bson.D{{}})
	return err
}
