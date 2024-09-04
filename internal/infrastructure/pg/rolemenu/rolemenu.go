package rolemenu

import (
	"github.com/google/uuid"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/rolemenu"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/structure"
)

type Model struct {
	ID             uuid.UUID
	RoleID         string
	MenuID         string
	ActionID       string
	IDString       *string
	RoleIDString   *string
	MenuIDString   *string
	ActionIDString *string
}

func (Model) TableName() string {
	return "rolemenu"
}

func (a Model) ToDomain() *rolemenu.RoleMenu {
	item := new(rolemenu.RoleMenu)
	structure.CopyWithUUID(a, item)
	return item
}

func toDomainList(ms []*Model) []*rolemenu.RoleMenu {
	list := make([]*rolemenu.RoleMenu, len(ms))
	for i, item := range ms {
		list[i] = item.ToDomain()
	}
	return list
}

func domainToModel(u *rolemenu.RoleMenu) *Model {
	item := new(Model)
	structure.CopyWithUUID(u, item)
	return item
}
