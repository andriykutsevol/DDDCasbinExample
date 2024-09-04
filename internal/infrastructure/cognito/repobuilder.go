package cognito

import (
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure"
)

func BuildRespositories(infrarepos *infrastructure.InfraRepos) (*infrastructure.InfraRepos, error) {
	infrarepos.AuthRepository = NewRepository()
	return infrarepos, nil
}