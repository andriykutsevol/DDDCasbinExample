package menu

import (
	"time"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/structure"
)

// We use *string when string can be nil.
type Model struct {
	ID         string     `json:"id" bson:"id"`
	Name       string     `json:"name" bson:"name"`
	Sequence   int        `json:"sequence" bson:"sequence"`
	Icon       *string    `json:"icon" bson:"icon"`
	Router     *string    `json:"router" bson:"router"`
	ParentID   *string    `json:"parentid" bson:"parentid"`
	ParentPath *string    `json:"parentpath" bson:"parentpath"`
	ShowStatus int        `json:"showstatus" bson:"showstatus"`
	Status     int        `json:"status" bson:"status"`
	Memo       *string    `json:"memo" bson:"memo"`
	Creator    string     `json:"creator" bson:"creator"`
	CreatedAt  time.Time  `json:"creatredat" bson:"createdat"`
	UpdatedAt  time.Time  `json:"updatedat" bson:"updatedat"`
	DeletedAt  *time.Time `json:"deletedat" bson:"deletedat"`
}

func (Model) TableName() string {
	return "menus"
}

func (a Model) ToDomain() *menu.Menu {
	item := new(menu.Menu)
	structure.Copy(a, item)
	return item
}

func toDomainList(menus []*Model) []*menu.Menu {
	list := make([]*menu.Menu, len(menus))
	for i, item := range menus {
		list[i] = item.ToDomain()
	}
	return list
}

func domainToModel(m *menu.Menu) *Model {
	item := new(Model)
	structure.Copy(m, item)
	return item
}
