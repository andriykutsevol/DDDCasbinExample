package router

import (
	casbin "github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/auth"
	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http/handler"
	//"github.com/go-swagger/go-swagger/fixtures/goparsing/petstore/rest/handlers"
)

type Router interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

type router struct {
	auth           auth.Repository
	casbinEnforcer *casbin.SyncedEnforcer
	loginHandler   handler.Login
	menuHandler    handler.Menu
	roleHandler    handler.Role
	userHandler    handler.User
	demosHandler   handler.Demos
	weatherHandler handler.Weather
}

func NewRouter(
	auth auth.Repository,
	casbinEnforcer *casbin.SyncedEnforcer,
	loginHandler handler.Login,
	menuHandler handler.Menu,
	roleHandler handler.Role,
	userHandler handler.User,
	demosHandler handler.Demos,
	weatherHandler handler.Weather,
) Router {
	return &router{
		auth:           auth,
		casbinEnforcer: casbinEnforcer,
		loginHandler:   loginHandler,
		menuHandler:    menuHandler,
		roleHandler:    roleHandler,
		userHandler:    userHandler,
		demosHandler:   demosHandler,
		weatherHandler: weatherHandler,
	}
}

func (a *router) Register(app *gin.Engine) error {
	a.RegisterAPI(app)
	return nil
}

func (a *router) Prefixes() []string {
	return []string{
		"/api/",
	}
}
