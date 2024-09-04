package application

import (
	"context"
	"fmt"
	"time"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/userrole"
)

type Demos interface {
	HandleGet(ctx context.Context, dto HandleGet_Dto) error
	HandlePut(ctx context.Context, dto HandlePut_Dto) error
}

type demosApp struct {
	userRepo     user.Repository
	userRoleRepo userrole.Repository
}

func NewDemos(
	userRepo user.Repository,
	userRoleRepo userrole.Repository,
) Demos {
	return &demosApp{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
	}
}

func (a *demosApp) HandleGet(ctx context.Context, dto HandleGet_Dto) error {

	fmt.Println("Application layer HandleGet dto: ", dto.Id)

	//a.userRepo.Get(ctx, dto.Id)

	// Loop 10 times with a 0.5-second delay
	for i := 0; i < 20; i++ {
		fmt.Println("Iteration", i+1)
		time.Sleep(500 * time.Millisecond) // Sleep for 0.5 seconds
	}

	return nil
}

func (a *demosApp) HandlePut(ctx context.Context, dto HandlePut_Dto) error {

	fmt.Println("Application layer HandlePut dto: ", dto.Id)

	//a.userRoleRepo.Get(ctx, dto.Id)

	return nil
}
