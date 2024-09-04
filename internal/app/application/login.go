package application

import (
	"context"
	"log"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/auth"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/role"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/rolemenu"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuaction"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/userrole"

	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/hash"
)

type Login interface {
	Verify(ctx context.Context, userName, password string) (*user.User, error)
	GenerateToken(ctx context.Context, userID string) (*auth.Auth, error)
	DestroyToken(ctx context.Context, tokenString string) error
	//GetLoginInfo(ctx context.Context, userID string) (*user.User, error)
	//UpdatePassword(ctx context.Context, userID string, oldPassword, newPassword string) error
	//QueryUserMenuTree(ctx context.Context, userID string) (menu.Menus, error)
}

type loginApp struct {
	authRepo       auth.Repository
	userRepo       user.Repository
	roleRepo       role.Repository
	userRoleRepo   userrole.Repository
	userSvc        user.Service
	menuRepo       menu.Repository
	menuActionRepo menuaction.Repository
	roleMenuRepo   rolemenu.Repository
}

func NewLogin(
	authRepo auth.Repository,
	userRepo user.Repository,
	roleRepo role.Repository,
	userRoleRepo userrole.Repository,
	userSvc user.Service,
	menuRepo menu.Repository,
	menuActionRepo menuaction.Repository,
	roleMenuRepo rolemenu.Repository,
) Login {
	return &loginApp{
		authRepo:       authRepo,
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		userRoleRepo:   userRoleRepo,
		userSvc:        userSvc,
		menuRepo:       menuRepo,
		menuActionRepo: menuActionRepo,
		roleMenuRepo:   roleMenuRepo,
	}
}

func (l loginApp) Verify(ctx context.Context, userName, password string) (*user.User, error) {

	log.Println("app: Verify()")

	if rootUser := l.authRepo.FindRootUser(ctx, userName); rootUser != nil {
		if password == rootUser.Password {
			return &user.User{
				UserName: rootUser.UserName,
				Password: rootUser.Password,
			}, nil
		}
	}

	result, _, err := l.userRepo.Query(ctx, user.QueryParams{
		UserName: userName,
	})

	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.ErrInvalidUserName
	}

	item := result[0]
	if item.Password != hash.SHA1String(password) {
		return nil, errors.ErrInvalidPassword
	}
	if item.Status != 1 {
		return nil, errors.ErrUserDisable
	}
	return item, nil
}

func (l loginApp) GenerateToken(ctx context.Context, userID string) (*auth.Auth, error) {

	log.Println("GenerateToken()")

	auth, err := l.authRepo.GenerateToken(ctx, userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return auth, nil
}

func (l loginApp) DestroyToken(ctx context.Context, tokenString string) error {

	log.Println("DestroyToken()")

	err := l.authRepo.DestroyToken(ctx, tokenString)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
