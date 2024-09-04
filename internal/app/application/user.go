package application

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/auth"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/role"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/userrole"

	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/hash"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/uuid"
)

type User interface {
	Query(ctx context.Context, params user.QueryParams) (user.Users, *pagination.Pagination, error)
	QueryShow(ctx context.Context, params user.QueryParams) (user.Users, *pagination.Pagination, error)
	Get(ctx context.Context, id string) (*user.User, error)
	Create(ctx context.Context, item *user.User, roleIDs []string) (string, error)
	Update(ctx context.Context, id string, item *user.User, roleIDs []string) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}

type userApp struct {
	authRepo     auth.Repository
	rbacAdapter  RbacAdapter
	enforcer     *casbin.SyncedEnforcer
	userRepo     user.Repository
	userRoleRepo userrole.Repository
	roleRepo     role.Repository
}

func NewUser(
	authRepo auth.Repository,
	rbacAdapter RbacAdapter,
	enforcer *casbin.SyncedEnforcer,
	userRepo user.Repository,
	userRoleRepo userrole.Repository,
	roleRepo role.Repository,
) User {
	return &userApp{
		authRepo:     authRepo,
		rbacAdapter:  rbacAdapter,
		enforcer:     enforcer,
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
	}
}

func (a *userApp) Query(ctx context.Context, params user.QueryParams) (user.Users, *pagination.Pagination, error) {
	result, pr, err := a.userRepo.Query(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	return result, pr, nil
}

func (a *userApp) QueryShow(ctx context.Context, params user.QueryParams) (user.Users, *pagination.Pagination, error) {

	// result, pr, err := a.userRepo.Query(ctx, params)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// if result == nil {
	// 	return nil, nil, nil
	// }

	// userRoleResult, _, err := a.userRoleRepo.Query(ctx, userrole.QueryParam{
	// 	UserIDs: result.ToIDs(),
	// })
	// if err != nil {
	// 	return nil, nil, err
	// }

	// roleResult, _, err := a.roleRepo.Query(ctx, role.QueryParam{
	// 	IDs: userRoleResult.ToRoleIDs(),
	// })
	// if err != nil {
	// 	return nil, nil, err
	// }

	// return result.FillRoles(userRoleResult.ToUserIDMap(), roleResult.ToMap()), pr, nil

	return nil, nil, nil
}

func (a *userApp) Get(ctx context.Context, id string) (*user.User, error) {

	return nil, nil
}

func (a *userApp) Create(ctx context.Context, item *user.User, roleIDs []string) (string, error) {
	fmt.Println("User Application Create()")
	err := a.checkUserName(ctx, item)
	if err != nil {
		return "", err
	}

	item.Password = hash.SHA1String(item.Password)
	item.ID = uuid.MustString()

	for _, roleID := range roleIDs {
		urItem := new(userrole.UserRole)
		urItem.ID = uuid.MustString()
		urItem.UserID = item.ID
		urItem.RoleID = roleID
		err := a.userRoleRepo.Create(ctx, urItem)
		if err != nil {
			return "", err
		}
	}

	err = a.userRepo.Create(ctx, item)

	// // TODO
	// a.rbacAdapter.AddPolicyItemToChan(ctx, a.enforcer)

	return item.ID, err
}

func (a *userApp) checkUserName(ctx context.Context, item *user.User) error {
	if rootUser := a.authRepo.FindRootUser(ctx, item.UserName); rootUser != nil {
		return errors.New400Response("The user name is invalid")
	}

	return nil
}

func (a *userApp) Update(ctx context.Context, id string, item *user.User, roleIDs []string) error {

	return nil
}

func (a *userApp) compareUserRoles(oldUserRoles map[string]*userrole.UserRole, roleIDs []string) (addList []string, delList []*userrole.UserRole) {
	return nil, nil
}

func (a *userApp) Delete(ctx context.Context, id string) error {

	return nil
}

func (a *userApp) UpdateStatus(ctx context.Context, id string, status int) error {

	return nil
}
