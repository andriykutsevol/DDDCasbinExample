package menuactionresource

import "github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"

// TODO:
// DTOs: Use Data Transfer Objects (DTOs) to handle primary keys, keeping the domain model clean.

type MenuActionResource struct {
	ID             string
	ActionID       string
	Method         string
	Path           string
	ActionIDString *string
	IDString       *string
}
type MenuActionResources []*MenuActionResource

type QueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	MenuID          string
	MenuIDs         []string
}

func (a MenuActionResources) ToMenuActionIDMap() map[string]MenuActionResources {
	m := make(map[string]MenuActionResources)
	for _, item := range a {
		m[item.ActionID] = append(m[item.ActionID], item)
	}
	return m
}

func (a MenuActionResources) ToMap() map[string]*MenuActionResource {
	m := make(map[string]*MenuActionResource)
	for _, item := range a {
		m[item.Method+item.Path] = item
	}
	return m
}

func (a MenuActionResources) ToActionIDMap() map[string]MenuActionResources {
	m := make(map[string]MenuActionResources)
	for _, item := range a {
		m[item.ActionID] = append(m[item.ActionID], item)
	}
	return m
}
