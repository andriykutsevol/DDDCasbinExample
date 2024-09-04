package menu

import (
	"context"
	"fmt"
	"strings"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/contextx"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuaction"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuactionresource"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"

	//"github.com/andriykutsevol/DDDCasbinExample/internal/domain/trans"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/uuid"
)

type Service interface {
	Query(ctx context.Context, params QueryParam) (Menus, *pagination.Pagination, error)
	Get(ctx context.Context, id string) (*Menu, error)
	GetByIdString(ctx context.Context, id string) (*Menu, error)
	QueryActions(ctx context.Context, id string) (menuaction.MenuActions, error)
	Create(ctx context.Context, item *Menu) (string, error)
	Update(ctx context.Context, id string, item *Menu) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
	PurgeMmenu(ctx context.Context) error
}

func NewService(
	//transRepo trans.Repository,
	menuRepo Repository,
	menuActionRepo menuaction.Repository,
	menuActionResourceRepo menuactionresource.Repository,
) Service {
	return &service{
		//transRepo:              transRepo,
		menuRepo:               menuRepo,
		menuActionRepo:         menuActionRepo,
		menuActionResourceRepo: menuActionResourceRepo,
	}
}

type service struct {
	//transRepo              trans.Repository
	menuRepo               Repository
	menuActionRepo         menuaction.Repository
	menuActionResourceRepo menuactionresource.Repository
}

func (s *service) Query(ctx context.Context, params QueryParam) (Menus, *pagination.Pagination, error) {
	menuActionResult, _, err := s.menuActionRepo.Query(ctx, menuaction.QueryParam{})
	if err != nil {
		return nil, nil, err
	}

	menuResult, pr, err := s.menuRepo.Query(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	menuResult.FillMenuAction(menuActionResult.ToMenuIDMap())
	return menuResult, pr, nil
}

func (s *service) Get(ctx context.Context, id string) (*Menu, error) {

	item, err := s.menuRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.ErrNotFound
	}

	actions, err := s.QueryActions(ctx, id)
	if err != nil {
		return nil, err
	}
	item.Actions = actions

	return item, nil
}

func (s *service) GetByIdString(ctx context.Context, id string) (*Menu, error) {
	menu, err := s.menuRepo.GetByIdString(ctx, id)
	if err != nil {
		return nil, err
	}
	return menu, nil
}

func (s *service) QueryActions(ctx context.Context, id string) (menuaction.MenuActions, error) {
	result, _, err := s.menuActionRepo.Query(ctx, menuaction.QueryParam{
		MenuID: id,
	})
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	resourceResult, _, err := s.menuActionResourceRepo.Query(ctx, menuactionresource.QueryParam{
		MenuID: id,
	})
	if err != nil {
		return nil, err
	}
	result.FillResources(resourceResult.ToMenuActionIDMap())
	return result, nil
}

func (s *service) Create(ctx context.Context, item *Menu) (string, error) {
	if err := s.checkName(ctx, item); err != nil {
		return "", err
	}

	parentPath, err := s.getParentPath(ctx, item.ParentID)
	if err != nil {
		return "", err
	}

	item.ParentPath = parentPath
	item.ID = uuid.MustString()

	item.IDString = new(string)
	*item.IDString = item.Name + "::" + item.Router
	*item.IDString = strings.ToLower(*item.IDString)
	*item.IDString = strings.ReplaceAll(*item.IDString, " ", "_")

	err = s.menuRepo.Create(ctx, item)
	if err != nil {
		return "", err
	}

	err = s.createActions(ctx, item.ID, item.IDString, item.Actions)
	if err != nil {
		return "", err
	}

	return item.ID, nil
}

func (s *service) checkName(ctx context.Context, item *Menu) error {
	menus, _, err := s.menuRepo.Query(ctx, QueryParam{
		PaginationParam: pagination.Param{
			OnlyCount: true,
		},
		ParentID: item.ParentID,
		Name:     item.Name,
	})

	if err != nil {
		return err
	}

	if len(menus) > 0 {
		return errors.New400Response("The menu name already exists")
	}

	return nil

}

func (s *service) getParentPath(ctx context.Context, parentID string) (string, error) {
	if parentID == "" {
		return "", nil
	}

	pitem, err := s.menuRepo.Get(ctx, parentID)
	if err != nil {
		return "", err
	}

	if pitem == nil {
		return "", errors.ErrInvalidParent
	}

	return s.joinParentPath(pitem.ParentPath, pitem.ID), nil
}

func (s *service) joinParentPath(parent, id string) string {
	if parent != "" {
		return parent + "/" + id
	}
	return id
}

func (s *service) createActions(ctx context.Context, menuID string, menuIDString *string, items menuaction.MenuActions) error {
	for _, item := range items {
		item.ID = uuid.MustString()

		item.IDString = new(string)
		*item.IDString = *menuIDString + "::" + item.Code

		item.MenuID = menuID
		item.MenuIDString = menuIDString
		err := s.menuActionRepo.Create(ctx, item)
		if err != nil {
			return err
		}
		for _, ritem := range item.Resources {
			ritem.ID = uuid.MustString()
			ritem.ActionID = item.ID

			ritem.IDString = new(string)
			*ritem.IDString = *item.IDString + "::" + ritem.Method + "::" + ritem.Path

			ritem.ActionIDString = new(string)
			*ritem.ActionIDString = *item.IDString

			err := s.menuActionResourceRepo.Create(ctx, ritem)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *service) updateActions(ctx context.Context, menuID string, menuIDString *string, oldItems, newItems menuaction.MenuActions) error {
	addActions, delActions, updateActions := s.compareActions(oldItems, newItems)

	err := s.createActions(ctx, menuID, menuIDString, addActions)
	if err != nil {
		return err
	}

	for _, item := range delActions {
		err := s.menuActionRepo.Delete(ctx, item.ID)
		if err != nil {
			return err
		}

		err = s.menuActionResourceRepo.DeleteByActionID(ctx, item.ID)
		if err != nil {
			return err
		}
	}

	mOldItems := oldItems.ToMap()
	for _, item := range updateActions {
		oitem := mOldItems[item.Code]
		// only update action name
		if item.Name != oitem.Name {
			oitem.Name = item.Name
			err := s.menuActionRepo.Update(ctx, item.ID, oitem)
			if err != nil {
				return err
			}
		}

		// update new and delete, not update
		addResources, delResources := s.compareResources(oitem.Resources, item.Resources)
		for _, aritem := range addResources {
			aritem.ID = uuid.MustString()
			aritem.ActionID = oitem.ID
			err := s.menuActionResourceRepo.Create(ctx, aritem)
			if err != nil {
				return err
			}
		}

		for _, ditem := range delResources {
			err := s.menuActionResourceRepo.Delete(ctx, ditem.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *service) Update(ctx context.Context, id string, item *Menu) error {
	// if id == item.ParentID {
	// 	return errors.ErrInvalidParent
	// }

	// oldItem, err := s.Get(ctx, id)
	// if err != nil {
	// 	return err
	// }
	// if oldItem == nil {
	// 	return errors.ErrNotFound
	// }
	// if oldItem.Name != item.Name {
	// 	if err := s.checkName(ctx, item); err != nil {
	// 		return err
	// 	}
	// }

	// item.ID = oldItem.ID
	// item.Creator = oldItem.Creator
	// item.CreatedAt = oldItem.CreatedAt

	// if oldItem.ParentID != item.ParentID {
	// 	parentPath, err := s.getParentPath(ctx, item.ParentID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	item.ParentPath = parentPath
	// } else {
	// 	item.ParentPath = oldItem.ParentPath
	// }

	// return s.transRepo.Exec(ctx, func(ctx context.Context) error {
	// 	err := s.updateActions(ctx, id, oldItem.Actions, item.Actions)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	err = s.updateChildParentPath(ctx, oldItem, item)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return s.menuRepo.Update(ctx, id, item)
	// })

	return nil
}

func (s *service) compareActions(oldActions, newActions menuaction.MenuActions) (addList, delList, updateList menuaction.MenuActions) {
	mOldActions := oldActions.ToMap()
	mNewActions := newActions.ToMap()

	for k, item := range mNewActions {
		if _, ok := mOldActions[k]; ok {
			updateList = append(updateList, item)
			delete(mOldActions, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldActions {
		delList = append(delList, item)
	}
	return
}

func (s *service) compareResources(oldResources, newResources menuactionresource.MenuActionResources) (addList, delList menuactionresource.MenuActionResources) {
	mOldResources := oldResources.ToMap()
	mNewResources := newResources.ToMap()

	for k, item := range mNewResources {
		if _, ok := mOldResources[k]; ok {
			delete(mOldResources, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldResources {
		delList = append(delList, item)
	}
	return
}

func (s *service) updateChildParentPath(ctx context.Context, oldItem, newItem *Menu) error {
	if oldItem.ParentID == newItem.ParentID {
		return nil
	}

	opath := s.joinParentPath(oldItem.ParentPath, oldItem.ID)
	result, _, err := s.menuRepo.Query(contextx.NewNoTrans(ctx), QueryParam{
		PrefixParentPath: opath,
	})
	if err != nil {
		return err
	}

	npath := s.joinParentPath(newItem.ParentPath, newItem.ID)
	for _, menu := range result {
		parentPath := menu.ParentPath
		err = s.menuRepo.UpdateParentPath(ctx, menu.ID, fmt.Sprintf("%s%s", npath, parentPath[len(opath):]))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	// oldItem, err := s.menuRepo.Get(ctx, id)
	// if err != nil {
	// 	return err
	// }
	// if oldItem == nil {
	// 	return errors.ErrNotFound
	// }

	// _, pr, err := s.menuRepo.Query(ctx, QueryParam{
	// 	PaginationParam: pagination.Param{OnlyCount: true},
	// 	ParentID:        &id,
	// })
	// if err != nil {
	// 	return err
	// }
	// if pr.Total > 0 {
	// 	return errors.ErrNotAllowDeleteWithChild
	// }

	// return s.transRepo.Exec(ctx, func(ctx context.Context) error {
	// 	err = s.menuActionResourceRepo.DeleteByMenuID(ctx, id)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	err := s.menuActionRepo.DeleteByMenuID(ctx, id)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return s.menuRepo.Delete(ctx, id)
	// })

	return nil
}

func (s *service) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := s.menuRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if oldItem == nil {
		return errors.ErrNotFound
	}

	return s.menuRepo.UpdateStatus(ctx, id, status)
}

func (s *service) PurgeMmenu(ctx context.Context) error {

	_ = s.menuActionResourceRepo.Purge(ctx)
	_ = s.menuActionRepo.Purge(ctx)
	_ = s.menuRepo.Purge(ctx)

	return nil
}
