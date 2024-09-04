package menuaction

import (
	"github.com/google/uuid"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuaction"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/structure"
)

// type Model struct {
// 	ID     string `gorm:"column:id;primary_key;size:36;"`
// 	MenuID string `gorm:"column:menu_id;size:36;index;default:'';not null;"`
// 	Code   string `gorm:"column:code;size:100;default:'';not null;"`
// 	Name   string `gorm:"column:name;size:100;default:'';not null;"`
// }

type Model struct {
	ID           uuid.UUID
	MenuID       uuid.UUID
	Code         string
	Name         string
	IDString     *string
	MenuIDString *string
}

// Custom type converter from string to uuid.UUID
func stringToUUID(src string) (uuid.UUID, error) {
	return uuid.Parse(src)
}

func (Model) TableName() string {
	return "menuaction"
}

func (a Model) ToDomain() *menuaction.MenuAction {
	item := new(menuaction.MenuAction)
	structure.CopyWithUUID(a, item)
	return item
}

func toDomainList(ms []*Model) []*menuaction.MenuAction {
	list := make([]*menuaction.MenuAction, len(ms))
	for i, item := range ms {
		list[i] = item.ToDomain()
	}
	return list
}

func domainToModel(m *menuaction.MenuAction) *Model {
	item := new(Model)
	structure.CopyWithUUID(m, item)
	return item
}
