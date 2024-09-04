package menu

import (
	//"log"
	"time"

	"github.com/google/uuid"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/structure"
)

// type Model struct {
// 	ID         string		`json:"id" bson:"id"`
// 	Name       string		`json:"name" bson:"name"`
// 	Sequence   int			`json:"sequence" bson:"sequence"`
// 	// you can scan a string from a database query into a *string in Go.
// 	// This is particularly useful when dealing with nullable columns in databases where a field can have a null value.
// 	Icon       *string		`json:"icon" bson:"icon"`
// 	Router     *string		`json:"router" bson:"router"`
// 	ParentID   *string		`json:"parentid" bson:"parentid"`
// 	ParentPath *string		`json:"parentpath" bson:"parentpath"`
// 	ShowStatus int			`json:"showstatus" bson:"showstatus"`
// 	Status     int			`json:"status" bson:"status"`
// 	Memo       *string		`json:"memo" bson:"memo"`
// 	Creator    string		`json:"creator" bson:"creator"`
// 	// postgresql "timestamp" and "data" scan to "time.Time".
// 	CreatedAt  time.Time	`json:"creatredat" bson:"createdat"`
// 	UpdatedAt  *time.Time	`json:"updatedat" bson:"updatedat"`
// 	DeletedAt  *time.Time	`json:"deletedat" bson:"deletedat"`
// 	IDString	*string
// }

type Model struct {
	ID       uuid.UUID
	Name     string
	Sequence int
	// you can scan a string from a database query into a *string in Go.
	// This is particularly useful when dealing with nullable columns in databases where a field can have a null value.
	Icon       *string
	Router     *string
	ParentID   *uuid.UUID
	ParentPath *string
	ShowStatus int
	Status     int
	Memo       *string
	Creator    *uuid.UUID
	// postgresql "timestamp" and "data" scan to "time.Time".
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
	IDString  *string
}

func (Model) TableName() string {
	return "menu"
}

func (m Model) ToDomain() *menu.Menu {
	item := new(menu.Menu)
	structure.CopyWithUUID(m, item)
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
	structure.CopyWithUUID(m, item)
	return item
}
