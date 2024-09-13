package application

import (
	"context"
	//"fmt"

	casbin "github.com/casbin/casbin/v2"
	//"golang.org/x/tools/go/analysis/passes/nilness"

	//"github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
	// "github.com/andriykutsevol/DDDCasbinExample/internal/domain/trans"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/role"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/rolemenu"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/uuid"
)

// Hexagonal architecture - Driving Port
type Role interface {
	Query(ctx context.Context, params role.QueryParam) (role.Roles, *pagination.Pagination, error)
	Get(ctx context.Context, id string) (*role.Role, error)
	QueryRoleMenus(ctx context.Context, roleID string) (rolemenu.RoleMenus, error)
	Create(ctx context.Context, item *role.Role) (string, error)
	Update(ctx context.Context, id string, item *role.Role) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}

type roleApp struct {
	rbacAdapter  RbacAdapter
	enforcer     *casbin.SyncedEnforcer
	roleRepo     role.Repository
	roleMenuRepo rolemenu.Repository
	userRepo     user.Repository
}

func NewRole(
	rbacAdapter RbacAdapter,
	enforcer *casbin.SyncedEnforcer,
	roleRepo role.Repository,
	roleMenuRepo rolemenu.Repository,
	userRepo user.Repository,
) Role {
	return &roleApp{
		rbacAdapter:  rbacAdapter,
		enforcer:     enforcer,
		roleRepo:     roleRepo,
		roleMenuRepo: roleMenuRepo,
		userRepo:     userRepo,
	}
}

func (a *roleApp) Create(ctx context.Context, item *role.Role) (string, error) {

	for _, rmItem := range item.RoleMenus {
		rmItem.ID = uuid.MustString()
		rmItem.RoleID = item.ID
		err := a.roleMenuRepo.Create(ctx, rmItem)
		if err != nil {
			return "err", err
		}
	}

	err := a.roleRepo.Create(ctx, item)
	if err != nil {
		return "", err
	}

	// TODO: implement rbacAdapter
	// a.rbacAdapter.AddPolicyItemToChan(ctx, a.enforcer)

	return item.ID, nil
}

func (a *roleApp) Query(ctx context.Context, params role.QueryParam) (role.Roles, *pagination.Pagination, error) {

	result, pr, err := a.roleRepo.Query(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	return result, pr, nil
}

func (a *roleApp) Get(ctx context.Context, id string) (*role.Role, error) {

	// item, err := a.roleRepo.Get(ctx, id)
	// roleMenus, err := a.QueryRoleMenus(ctx, id)

	return nil, nil
}

func (a *roleApp) Delete(ctx context.Context, id string) error {

	// oldItem, err := a.roleRepo.Get(ctx, id)

	// _, pr, err := a.userRepo.Query(ctx, user.QueryParams{
	// 	PaginationParam: pagination.Param{OnlyCount: true},
	// 	RoleIDs:         []string{id},
	// })

	// err = a.transRepo.Exec(ctx, func(ctx context.Context) error {
	// 	err := a.roleMenuRepo.DeleteByRoleID(ctx, id)
	// 	return a.roleRepo.Delete(ctx, id)
	// }

	// a.rbacAdapter.AddPolicyItemToChan(ctx, a.enforcer)

	return nil
}

func (a *roleApp) QueryRoleMenus(ctx context.Context, roleID string) (rolemenu.RoleMenus, error) {

	// result, _, err := a.roleMenuRepo.Query(ctx, rolemenu.QueryParam{
	// 	RoleID: roleID,
	// })

	return nil, nil

}

func (a *roleApp) checkName(ctx context.Context, item *role.Role) error {

	// _, pr, err := a.roleRepo.Query(ctx, role.QueryParam{
	// 	PaginationParam: pagination.Param{OnlyCount: true},
	// 	Name:            item.Name,
	// })

	return nil
}

func (a *roleApp) Update(ctx context.Context, id string, item *role.Role) error {

	// oldItem, err := a.Get(ctx, id)

	// err = a.transRepo.Exec(ctx, func(ctx context.Context) error {
	// 	return a.roleRepo.Update(ctx, id, item)
	// }

	// a.rbacAdapter.AddPolicyItemToChan(ctx, a.enforcer)

	return nil
}

func (a *roleApp) UpdateStatus(ctx context.Context, id string, status int) error {

	// oldItem, err := a.roleRepo.Get(ctx, id)

	// err = a.roleRepo.UpdateStatus(ctx, id, status)

	// a.rbacAdapter.AddPolicyItemToChan(ctx, a.enforcer)

	return nil
}
