package redis

import (
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure"
	authinfra "github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/redis/auth"
)

func BuildRespositories(infrarepos *infrastructure.InfraRepos) (*infrastructure.InfraRepos, error) {
	authRedisStorage := &authinfra.AuthStorage{}
	infrarepos.AuthRepository = authinfra.NewRepository(authRedisStorage)
	return infrarepos, nil
}
