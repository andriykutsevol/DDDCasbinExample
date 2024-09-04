package middleware

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"github.com/andriykutsevol/DDDCasbinExample/configs"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"
	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http"
)

func CasbinMiddleware(enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) gin.HandlerFunc {

	log.Println("CasbinMiddleware")

	cfg := configs.C.Casbin

	if !cfg.Enable {
		return EmptyMiddleware()
	}

	return func(c *gin.Context) {

		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		m := c.Request.Method

		//TODO debug mode
		//fmt.Println(http.GetUserID(c), p, m)

		if b, err := enforcer.Enforce(http.GetUserID(c), p, m); err != nil {
			http.ResError(c, errors.WithStack(err))
			return
		} else if !b {
			http.ResError(c, errors.ErrNoPerm)
			return
		}
		c.Next()
	}
}
