package handler

import (
	//"github.com/LyricTian/captcha"

	"github.com/gin-gonic/gin"

	"github.com/andriykutsevol/DDDCasbinExample/internal/app/application"

	// "github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"
	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http"
	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http/request"
	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http/response"

	// "github.com/andriykutsevol/DDDCasbinExample/configs"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/structure"
	//"github.com/linzhengen/ddd-gin-admin/pkg/logger"
)

type Login interface {
	GetCaptcha(c *gin.Context)
	ResCaptcha(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	RefreshToken(c *gin.Context)
	GetUserInfo(c *gin.Context)
	QueryUserMenuTree(c *gin.Context)
	UpdatePassword(c *gin.Context)
}

type login struct {
	loginApp application.Login
}

func NewLogin(loginApp application.Login) Login {
	return &login{
		loginApp: loginApp,
	}
}

func (a *login) GetCaptcha(c *gin.Context) {

}

func (a *login) ResCaptcha(c *gin.Context) {

}

func (a *login) Login(c *gin.Context) {

	ctx := c.Request.Context()
	var item request.LoginParam
	if err := http.ParseJSON(c, &item); err != nil {
		http.ResError(c, err)
		return
	}

	// if !captcha.VerifyString(item.CaptchaID, item.CaptchaCode) {
	// 	api.ResError(c, errors.New400Response("Invalid Captcha"))
	// 	return
	// }

	user, err := a.loginApp.Verify(ctx, item.UserName, item.Password)
	if err != nil {
		http.ResError(c, err)
		return
	}

	userID := user.ID
	http.SetUserID(c, userID)

	tokenInfo, err := a.loginApp.GenerateToken(ctx, userID)
	if err != nil {
		http.ResError(c, err)
		return
	}

	respTokenInfo := new(response.LoginTokenInfo)
	structure.Copy(tokenInfo, respTokenInfo)
	// ctx = logger.NewUserIDContext(ctx, userID)
	// ctx = logger.NewTagContext(ctx, "__login__")
	// logger.WithContext(ctx).Infof("logged in")
	http.ResSuccess(c, respTokenInfo)

}

func (a *login) Logout(c *gin.Context) {

}

func (a *login) RefreshToken(c *gin.Context) {

}

func (a *login) GetUserInfo(c *gin.Context) {

}

func (a *login) QueryUserMenuTree(c *gin.Context) {

}

func (a *login) UpdatePassword(c *gin.Context) {

}
