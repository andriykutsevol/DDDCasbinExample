package infrastructure

import (
	authrepo "github.com/andriykutsevol/DDDCasbinExample/internal/domain/auth"
	menurepo "github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu"
	menuactionrepo "github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuaction"
	menuactionresourcerepo "github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuactionresource"
	rbacrepo "github.com/andriykutsevol/DDDCasbinExample/internal/domain/rbac"
	userrepo "github.com/andriykutsevol/DDDCasbinExample/internal/domain/user"
	rolerepo "github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/role"
	rolemenurepo "github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/rolemenu"
	userrolerepo "github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/userrole"
	weatherrepo "github.com/andriykutsevol/DDDCasbinExample/internal/domain/weather"
)

type InfraRepos struct {
	RoleRepository               rolerepo.Repository
	MenuRepository               menurepo.Repository
	RoleMeuRepository            rolemenurepo.Repository
	MenuActionRepository         menuactionrepo.Repository
	MenuActionResourceRepository menuactionresourcerepo.Repository
	UserRepository               userrepo.Repository
	UserRoleRepository           userrolerepo.Repository
	AuthRepository               authrepo.Repository
	RbacRepository               rbacrepo.Repository
	WeatherRepository            weatherrepo.Repository
}
