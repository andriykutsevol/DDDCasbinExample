package pg

import (
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure"
	menuinfra "github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/menu"
	menuactioninfra "github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/menuaction"
	menuactionresourceinfra "github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/menuactionresource"
	roleinfra "github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/role"
	rolemenuinfra "github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/rolemenu"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/storage"
	userinfra "github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/user"
	userroleinfra "github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/userrole"

	rbacinfra "github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/rbac"
	//authinfra "github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/redis/auth"

	weatherinfra "github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/weather"
)

func BuildRespositories(
	infrarepos  *infrastructure.InfraRepos,
	casbinservice storage.DatabaseService,
	weatherservice storage.DatabaseService) (*infrastructure.InfraRepos, error) {

	//infrarepos := &infrastructure.InfraRepos{}

	infrarepos.MenuRepository = menuinfra.NewRepository(casbinservice)
	infrarepos.MenuActionRepository = menuactioninfra.NewRepository(casbinservice)
	infrarepos.MenuActionResourceRepository = menuactionresourceinfra.NewRepository(casbinservice, infrarepos.MenuActionRepository)
	infrarepos.RoleRepository = roleinfra.NewRepository(casbinservice)
	infrarepos.RoleMeuRepository = rolemenuinfra.NewRepository(casbinservice)
	infrarepos.UserRepository = userinfra.NewRepository(casbinservice)
	infrarepos.UserRoleRepository = userroleinfra.NewRepository(casbinservice)

	// It is conditional now.
	// authRedisStorage := &authinfra.AuthStorage{}
	// infrarepos.AuthRepository = authinfra.NewRepository(authRedisStorage)

	infrarepos.RbacRepository = rbacinfra.NewRepository(
		infrarepos.RoleRepository,
		infrarepos.RoleMeuRepository,
		infrarepos.MenuActionResourceRepository,
		infrarepos.UserRepository,
		infrarepos.UserRoleRepository)

	infrarepos.WeatherRepository = weatherinfra.NewRepository(weatherservice)

	return infrarepos, nil
}
