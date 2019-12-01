package infrastructure

import (
	"backend/interfaces/controllers"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Router is gin router
func Router() *gin.Engine {
	conf := NewConf()
	controller := controllers.NewController(NewDB(conf.DBConf.Database, conf.DBConf.DSN))

	router := gin.Default()
	v1 := router.Group("/v1")

	auth := v1.Group("/auth")
	github := auth.Group("/github")
	github.GET("/login", func(c *gin.Context) { controller.Login(c, conf) })
	github.GET("/callback", func(c *gin.Context) { controller.Callback(c, conf) })
	github.GET("/token", func(c *gin.Context) { controller.Auth(c) })

	zap.S().Info("running")
	return router
}
