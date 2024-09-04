package handler

import (
	nethttp "net/http"
	//"encoding/json"
	"github.com/gin-gonic/gin"

	"github.com/andriykutsevol/DDDCasbinExample/internal/app/application"
	//"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http"
	//httprequest "github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http/request"
	//httpresponse "github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http/response"
)

type Weather interface {
	WeatherGet(c *gin.Context)
}

type weatherHandler struct {
	weatherApp application.Weather
}

func NewWeather(weatherApp application.Weather) Weather {
	return &weatherHandler{
		weatherApp: weatherApp,
	}
}

//-----------------------------------------------------------

func (w *weatherHandler) WeatherGet(c *gin.Context) {
	city := c.Param("id")
	_, app_response := w.weatherApp.HandleGet(c, city)
	c.JSON(nethttp.StatusOK, gin.H{"message": app_response})
}
