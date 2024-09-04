package userrole

import "github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"

// TODO:
// DTOs: Use Data Transfer Objects (DTOs) to handle primary keys, keeping the domain model clean.

type UserRole struct {
	ID           string
	UserID       string
	RoleID       string
	IDString     *string
	UserIDString *string
	RoleIDString *string
}

type QueryParam struct {
	PaginationParam pagination.Param
	OrderFields     pagination.OrderFields
	UserID          string
	UserIDs         []string
	IDString        *string
	UserIDString    *string
	RoleIDString    *string
}

type UserRoles []*UserRole

func (a UserRoles) ToRoleIDs() []string {
	ids := make([]string, len(a))
	for i, item := range a {
		ids[i] = item.RoleID
	}
	return ids
}

func (a UserRoles) ToUserIDMap() map[string]UserRoles {
	m := make(map[string]UserRoles)
	for _, item := range a {
		m[item.UserID] = append(m[item.UserID], item)
	}
	return m
}

func (a UserRoles) ToMap() map[string]*UserRole {
	m := make(map[string]*UserRole)
	for _, item := range a {
		m[item.RoleID] = item
	}
	return m
}
