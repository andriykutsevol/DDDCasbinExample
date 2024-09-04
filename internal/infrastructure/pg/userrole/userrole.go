package userrole

import (
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/userrole"

	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/structure"
)

// type Model struct {
// 	ID     string `gorm:"column:id;primary_key;size:36;"`
// 	UserID string `gorm:"column:user_id;size:36;index;default:'';not null;"`
// 	RoleID string `gorm:"column:role_id;size:36;index;default:'';not null;"`
// }

type Model struct {
	ID           string
	UserID       string
	RoleID       string
	IDString     *string
	UserIDString *string
	RoleIDString *string
}

func (Model) TableName() string {
	return "user_roles"
}

func (a Model) ToDomain() *userrole.UserRole {
	item := new(userrole.UserRole)
	structure.CopyWithUUID(a, item)
	return item
}

func toDomainList(ms []*Model) []*userrole.UserRole {
	list := make([]*userrole.UserRole, len(ms))
	for i, item := range ms {
		list[i] = item.ToDomain()
	}
	return list
}

func domainToModel(u *userrole.UserRole) *Model {
	item := new(Model)
	structure.CopyWithUUID(u, item)
	return item
}
