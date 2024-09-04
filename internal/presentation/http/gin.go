package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"

	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http/response"

	errorsd "github.com/andriykutsevol/DDDCasbinExample/internal/domain/errors"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/json"
)

const (
	prefix           = "ddd-gin-admin"
	UserIDKey        = prefix + "/user-id"
	ReqBodyKey       = prefix + "/req-body"
	ResBodyKey       = prefix + "/res-body"
	LoggerReqBodyKey = prefix + "/logger-req-body"
)

func GetToken(c *gin.Context) string {

	var token string
	auth := c.GetHeader("Authorization")

	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

func GetUserID(c *gin.Context) string {
	return c.GetString(UserIDKey)
}

func SetUserID(c *gin.Context, userID string) {
	c.Set(UserIDKey, userID)
}

func GetBody(c *gin.Context) []byte {
	if v, ok := c.Get(ReqBodyKey); ok {
		if b, ok := v.([]byte); ok {
			return b
		}
	}
	return nil
}

func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errorsd.Wrap400Response(err, fmt.Sprintf("400 Bad Request - %s", err.Error()))
	}
	return nil
}

func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errorsd.Wrap400Response(err, fmt.Sprintf("400 Bad Request - %s", err.Error()))
	}
	return nil
}

func ParseForm(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		return errorsd.Wrap400Response(err, fmt.Sprintf("400 Bad Request - %s", err.Error()))
	}
	return nil
}

func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	c.Set(ResBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}

func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, v)
}

func ResOK(c *gin.Context) {
	ResSuccess(c, response.StatusResult{Status: response.OKStatus})
}

func ResList(c *gin.Context, v interface{}) {
	ResSuccess(c, response.ListResult{List: v})
}

func ResPage(c *gin.Context, v interface{}, pr *pagination.Pagination) {
	list := response.ListResult{
		List:       v,
		Pagination: pr,
	}
	ResSuccess(c, list)
}

func ResError(c *gin.Context, err error, status ...int) {
	//ctx := c.Request.Context()
	var res *errorsd.ResponseError

	if err != nil {
		if e, ok := err.(*errorsd.ResponseError); ok {
			res = e
		} else {
			res = errorsd.UnWrapResponse(errorsd.ErrInternalServer)
			res.ERR = err
		}
	} else {
		res = errorsd.UnWrapResponse(errorsd.ErrInternalServer)
	}

	if len(status) > 0 {
		res.StatusCode = status[0]
	}

	if err := res.ERR; err != nil {
		if res.Message == "" {
			res.Message = err.Error()
		}

		if status := res.StatusCode; status >= 400 && status < 500 {
			//logger.WithContext(ctx).Warnf(err.Error())
			log.Println(err.Error())
		} else if status >= 500 {
			//logger.WithContext(logger.NewStackContext(ctx, err)).Errorf(err.Error())
			log.Println(err.Error())
		}
	}

	eitem := response.ErrorItem{
		Code:    res.Code,
		Message: res.Message,
	}
	ResJSON(c, res.StatusCode, response.ErrorResult{Error: eitem})
}
