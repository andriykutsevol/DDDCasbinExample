package middleware

import (
	//"fmt"
	"github.com/gin-gonic/gin"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/auth"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/contextx"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"
	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http"
)

func wrapUserAuthContext(c *gin.Context, userID string) {
	http.SetUserID(c, userID)
	ctx := contextx.NewUserID(c.Request.Context(), userID)
	// TODO. Setpu logging
	//ctx = logger.NewUserIDContext(ctx, userID)
	c.Request = c.Request.WithContext(ctx)
}

func UserAuthMiddleware(a auth.Repository, skippers ...SkipperFunc) gin.HandlerFunc {
	// TODO. Look at the original files. Handle Authentication.

	return func(c *gin.Context) {

		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		userID, err := a.ParseUserID(c.Request.Context(), http.GetToken(c))
		if err != nil {
			if err == auth.ErrInvalidToken {
				http.ResError(c, errors.ErrInvalidToken)
				return
			}
			http.ResError(c, errors.WithStack(err))
			return
		}

		wrapUserAuthContext(c, userID)
		c.Next()
	}

}
