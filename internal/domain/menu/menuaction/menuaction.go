package menuaction

import (
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuactionresource"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
)

// TODO:
// DTOs: Use Data Transfer Objects (DTOs) to handle primary keys, keeping the domain model clean.

type MenuAction struct {
	ID           string
	MenuID       string
	Code         string
	Name         string
	Resources    menuactionresource.MenuActionResources
	IDString     *string
	MenuIDString *string
}

type MenuActions []*MenuAction

type QueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	MenuID          string
	IDs             []string
}

func (a MenuActions) ToMenuIDMap() map[string]MenuActions {
	m := make(map[string]MenuActions)
	for _, item := range a {
		m[item.MenuID] = append(m[item.MenuID], item)
	}
	return m
}

func (a MenuActions) FillResources(mResources map[string]menuactionresource.MenuActionResources) {
	for i, item := range a {
		a[i].Resources = mResources[item.ID]
	}
}

func (a MenuActions) ToMap() map[string]*MenuAction {
	m := make(map[string]*MenuAction)
	for _, item := range a {
		m[item.Code] = item
	}
	return m
}
