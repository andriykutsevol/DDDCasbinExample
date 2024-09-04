package user

import (
	"time"

	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/structure"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/role"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/userrole"
)

// TODO:
// DTOs: Use Data Transfer Objects (DTOs) to handle primary keys, keeping the domain model clean.

type User struct {
	ID        string
	UserName  string
	RealName  string
	Password  string
	Email     *string
	Phone     *string
	Status    int
	Creator   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Roles     role.Roles
	IDString  *string
}

func (a User) FillRoles(userRoles map[string]userrole.UserRoles, roles map[string]*role.Role) *User {
	u := new(User)
	structure.Copy(a, u)
	for _, roleID := range userRoles[a.ID].ToRoleIDs() {
		if v, ok := roles[roleID]; ok {
			u.Roles = append(u.Roles, v)
		}
	}
	return u
}

type Users []*User

type QueryParams struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	UserName        string
	QueryValue      string
	Status          int
	RoleIDs         []string
}

func (a Users) ToIDs() []string {
	idList := make([]string, len(a))
	for i, item := range a {
		idList[i] = item.ID
	}
	return idList
}

func (a Users) FillRoles(userRoles map[string]userrole.UserRoles, roles map[string]*role.Role) Users {
	list := make(Users, len(a))
	for i, item := range a {
		list[i] = item.FillRoles(userRoles, roles)
	}

	return list
}
