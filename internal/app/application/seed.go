package application

import (
	"context"
	"log"
	"os"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuaction"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu/menuactionresource"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/role"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/rolemenu"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user/userrole"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/weather"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/auth"

	//"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
	

	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/hash"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/uuid"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/yaml"

	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	// //"github.com/casbin/casbin/v2"

)


type Seed interface {
	Execute(ctx context.Context) error
}

type SeedMenus []struct {
	Name     string `yaml:"name"`
	Icon     string `yaml:"icon"`
	Router   string `yaml:"router,omitempty"`
	Sequence int    `yaml:"sequence"`
	Actions  []struct {
		Code      string `yaml:"code"`
		Name      string `yaml:"name"`
		Resources []struct {
			Method string `yaml:"method"`
			Path   string `yaml:"path"`
		} `yaml:"resources"`
	} `yaml:"actions,omitempty"`
	Children SeedMenus
}

type SeedRoles []struct {
	Name      string `yaml:"name"`
	Sequence  int    `yaml:"sequence"`
	Memo      string `yaml:"memo"`
	Rolemenus []struct {
		MenuID   string `yaml:"menuid"`
		ActionID string `yaml:"actionid"`
	} `yaml:"rolemenus,omitempty"`
}

type SeedUsers []struct {
	Name     string `yaml:"name"`
	RealName string `yaml:"realname"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
	Phone    string `yaml:"phone"`
	RoleIDs  []struct {
		RoleID string `yaml:"roleid"`
	} `yaml:"roleids"`
}


type seedApp struct {
	menuSvc        menu.Service
	menuactionrepo menuaction.Repository
	roleRepo       role.Repository
	roleMenuRepo   rolemenu.Repository
	userRepo       user.Repository
	userRoleRepo   userrole.Repository
	weatherRepo    weather.Repository
	authRepo       auth.Repository
}

func NewSeed(
	menuSvc menu.Service,
	menuactionrepo menuaction.Repository,
	roleRepo role.Repository,
	roleMenuRepo rolemenu.Repository,
	userRepo user.Repository,
	userRoleRepo userrole.Repository,
	weatherRepo weather.Repository,
	authRepo auth.Repository,
) Seed {
	return &seedApp{
		menuSvc:        menuSvc,
		menuactionrepo: menuactionrepo,
		roleRepo:       roleRepo,
		roleMenuRepo:   roleMenuRepo,
		userRepo:       userRepo,
		userRoleRepo:   userRoleRepo,
		weatherRepo:    weatherRepo,
		authRepo: authRepo,
	}
}



//================================================================================================
//================================================================================================

var (
    userPoolID = "your_user_pool_id"
    clientID   = "your_client_id"
    region     = "your_region"
)


// If we run this, we suppose that db schema is empty
// In the development mode we just purge all tables.
func (s seedApp) Execute(ctx context.Context) error {

	// We have to reach login repository directly here.
	// For initial testing of Cognito

	s.authRepo.GenerateToken(ctx, "userID")




	// if err := s.menuSeed(ctx, "../../configs/menu.yaml"); err != nil {
	// 	return err
	// }
	// if err := s.roleSeed(ctx, "../../configs/role.yaml"); err != nil {
	// 	return err
	// }
	// if err := s.userSeed(ctx, "../../configs/user.yaml"); err != nil {
	// 	return err
	// }
	// if err := s.weatherSeed(ctx); err != nil {
	// 	return err
	// }


	log.Println("Seed database has been recreated")

	return nil
}

//================================================================================================
//================================================================================================

func (s seedApp) userSeed(ctx context.Context, menuRolePath string) error {

	err := s.userRepo.Purge(ctx)
	if err != nil {
		return err
	}

	err = s.userRoleRepo.Purge(ctx)
	if err != nil {
		return err
	}

	data, err := s.readUserData(menuRolePath)
	if err != nil {
		return err
	}
	err = s.createUsers(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (s seedApp) createUsers(ctx context.Context, usersSeed SeedUsers) error {

	for _, userSeed := range usersSeed {

		user := &user.User{
			ID:       uuid.MustString(),
			UserName: userSeed.Name,
			RealName: userSeed.RealName,
			Password: hash.SHA1String(userSeed.Password),
			Email:    &userSeed.Email,
			Phone:    &userSeed.Phone,
			Status:   1,
			IDString: &userSeed.Name,
		}

		var userRoles userrole.UserRoles

		for _, userRoleID := range userSeed.RoleIDs {

			roleIDString := userRoleID.RoleID
			idString := *user.IDString + "::" + roleIDString
			userIDString := user.IDString

			userRoles = append(userRoles, &userrole.UserRole{
				ID:           uuid.MustString(),
				UserID:       user.ID,
				RoleID:       userRoleID.RoleID,
				IDString:     &idString,
				UserIDString: userIDString,
				RoleIDString: &roleIDString,
			})
		}

		err := s.createUser(ctx, user, &userRoles)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s seedApp) createUser(ctx context.Context, user *user.User, userRoles *userrole.UserRoles) error {

	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	for _, userrole := range *userRoles {

		role, err := s.roleRepo.GetByIdString(ctx, *userrole.RoleIDString)
		if err != nil {
			return err
		}

		userrole.RoleID = role.ID

		err = s.userRoleRepo.Create(ctx, userrole)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s seedApp) readUserData(name string) (SeedUsers, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data SeedUsers
	d := yaml.NewDecoder(file)
	d.SetStrict(true)
	err = d.Decode(&data)
	return data, err
}

//================================================================================================
//================================================================================================

func (s seedApp) roleSeed(ctx context.Context, menuRolePath string) error {
	_ = s.roleRepo.Purge(ctx)
	_ = s.roleMenuRepo.Purge(ctx)

	data, err := s.readRoleData(menuRolePath)
	if err != nil {
		return err
	}

	err = s.createRoles(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (s seedApp) createRoles(ctx context.Context, list SeedRoles) error {
	for _, item := range list {

		roleid := uuid.MustString()
		var rms rolemenu.RoleMenus

		for _, rmenu := range item.Rolemenus {

			roleIDString := item.Name
			menuIDString := rmenu.MenuID
			actionIDString := rmenu.ActionID

			rm := &rolemenu.RoleMenu{
				RoleID:         roleid,
				MenuID:         rmenu.MenuID,
				ActionID:       rmenu.ActionID,
				RoleIDString:   &roleIDString,
				MenuIDString:   &menuIDString,
				ActionIDString: &actionIDString,
			}
			rms = append(rms, rm)
		}

		ritem := &role.Role{
			ID:        roleid,
			Name:      item.Name,
			Sequence:  item.Sequence,
			Memo:      &item.Memo,
			Status:    1,
			RoleMenus: rms,
			IDString:  &item.Name,
		}

		err := s.createRole(ctx, ritem)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s seedApp) createRole(ctx context.Context, item *role.Role) error {
	err := s.roleRepo.Create(ctx, item)
	if err != nil {
		return err
	}

	for _, rmItem := range item.RoleMenus {

		rmItem.ID = uuid.MustString()
		rmItem.RoleID = item.ID

		menu, err := s.menuSvc.GetByIdString(ctx, *rmItem.MenuIDString)
		if err != nil {
			return err
		}
		rmItem.MenuID = menu.ID

		menuAction, err := s.menuactionrepo.GetByIdString(ctx, *rmItem.ActionIDString)
		if err != nil {
			return err
		}
		rmItem.ActionID = menuAction.ID

		rmItem.IDString = new(string)
		*rmItem.IDString = *item.IDString + "::" + *rmItem.ActionIDString

		rmItem.RoleIDString = item.IDString

		err = s.roleMenuRepo.Create(ctx, rmItem)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s seedApp) readRoleData(name string) (SeedRoles, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data SeedRoles
	d := yaml.NewDecoder(file)
	d.SetStrict(true)
	err = d.Decode(&data)
	return data, err

}

//================================================================================================
//================================================================================================

func (s seedApp) menuSeed(ctx context.Context, menuSeedPath string) error {
	// TODO. We are not doing any checks at this time. We assume that menu will be created from scratch.

	_ = s.roleMenuRepo.Purge(ctx)
	_ = s.userRoleRepo.Purge(ctx)
	_ = s.menuSvc.PurgeMmenu(ctx)

	data, err := s.readMenuData(menuSeedPath)
	if err != nil {
		return err
	}
	err = s.createMenus(ctx, "", data)
	if err != nil {
		return err
	}
	return nil

}

func (s seedApp) readMenuData(name string) (SeedMenus, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data SeedMenus
	d := yaml.NewDecoder(file)
	d.SetStrict(true)
	err = d.Decode(&data)
	return data, err
}

func (s seedApp) createMenus(ctx context.Context, parentID string, list SeedMenus) error {
	for _, item := range list {
		var as menuaction.MenuActions
		for _, action := range item.Actions {
			var ars menuactionresource.MenuActionResources
			for _, r := range action.Resources {
				ars = append(ars, &menuactionresource.MenuActionResource{
					Method: r.Method,
					Path:   r.Path,
				})
			}
			as = append(as, &menuaction.MenuAction{
				Code:      action.Code,
				Name:      action.Name,
				Resources: ars,
			})
		}
		sitem := &menu.Menu{
			Name:       item.Name,
			Sequence:   item.Sequence,
			Icon:       item.Icon,
			Router:     item.Router,
			ParentID:   parentID,
			Status:     1,
			ShowStatus: 1,
			Actions:    as,
		}

		menuID, err := s.menuSvc.Create(ctx, sitem)
		if err != nil {
			return err
		}

		if item.Children != nil && len(item.Children) > 0 {
			err := s.createMenus(ctx, menuID, item.Children)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//================================================================================================
//================================================================================================

func (s seedApp) weatherSeed(ctx context.Context) error {

	seedParams := map[string]string{
		"id": "dnipro",
	}
	err := s.weatherRepo.Seed(ctx, seedParams)
	if err != nil {
		return err
	}
	return nil
}
