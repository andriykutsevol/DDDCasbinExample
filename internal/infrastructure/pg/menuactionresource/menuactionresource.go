package menuactionresource

import (
	"github.com/google/uuid"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuactionresource"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/structure"
)

type Model struct {
	ID             uuid.UUID
	ActionID       uuid.UUID
	Method         string
	Path           string
	IDString       *string
	ActionIDString *string
}

func (Model) TableName() string {
	return "menuactionresource"
}

func (a Model) ToDomain() *menuactionresource.MenuActionResource {
	item := new(menuactionresource.MenuActionResource)
	structure.CopyWithUUID(a, item)
	return item
}

func toDomainList(ms []*Model) []*menuactionresource.MenuActionResource {
	list := make([]*menuactionresource.MenuActionResource, len(ms))
	for i, item := range ms {
		list[i] = item.ToDomain()
	}
	return list
}

func domainToModel(m *menuactionresource.MenuActionResource) *Model {
	item := new(Model)
	structure.CopyWithUUID(m, item)
	return item
}
