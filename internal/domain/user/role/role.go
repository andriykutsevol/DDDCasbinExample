package role

import (
	"time"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/rolemenu"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
)

type Role struct {
	ID        string
	Name      string
	Sequence  int
	Memo      *string
	Status    int
	Creator   *string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
	RoleMenus rolemenu.RoleMenus
	IDString  *string
}

type Roles []*Role

type QueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	IDs             []string
	Name            string
	QueryValue      string
	UserID          string
	Status          int
}

func (a Roles) ToMap() map[string]*Role {
	m := make(map[string]*Role)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}
