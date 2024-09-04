package middleware

import (
	"time"
	//"github.com/andriykutsevol/DDDCasbinExample/configs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	//cfg := configs.C.CORS

	// return cors.New(cors.Config{
	// 	AllowOrigins:     cfg.AllowOrigins,
	// 	AllowMethods:     cfg.AllowMethods,
	// 	AllowHeaders:     cfg.AllowHeaders,
	// 	AllowCredentials: cfg.AllowCredentials,
	// 	MaxAge:           time.Second * time.Duration(cfg.MaxAge),
	// })

	//TODO: read from cfg
	return cors.New(cors.Config{
		AllowAllOrigins:  true,                                                         // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, // Allow all methods
		AllowHeaders:     []string{"*"},                                                // Allow all headers
		ExposeHeaders:    []string{"*"},                                                // Expose all headers
		AllowCredentials: true,                                                         // Allow credentials
		MaxAge:           12 * time.Hour,                                               // Preflight request cache duration
	})

}
