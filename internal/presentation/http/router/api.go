package router

import (
	"github.com/gin-gonic/gin"

	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http/middleware"
)

func (a *router) RegisterAPI(app *gin.Engine) {

	g := app.Group("/api")

	g.Use(middleware.UserAuthMiddleware(a.auth,
		middleware.AllowPathPrefixSkipper("/api/v1/pub/login"),
	))

	g.Use(middleware.CasbinMiddleware(a.casbinEnforcer,
		middleware.AllowPathPrefixSkipper("/api/v1/pub"),
	))

	v1 := g.Group("/v1")
	{
		pub := v1.Group("/pub")
		{
			gLogin := pub.Group("login")
			{
				gLogin.POST("", a.loginHandler.Login)
			}
		}
		gWeather := v1.Group("weather")
		{
			gWeather.GET(":id", a.weatherHandler.WeatherGet)
		}
	}
}
