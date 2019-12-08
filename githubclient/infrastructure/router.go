package infrastructure

import (
	"githubclient/interfaces/controllers"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Router is routing API path
func Router() *gin.Engine {
	controller := controllers.NewController()

	router := gin.Default()
	v1 := router.Group("/v1")

	github := v1.Group("/github")
	// githubトークンを基にユーザ判定
	github.GET("/user", func(c *gin.Context) { controller.GetUser(c) })

	zap.S().Info("running")
	return router
}
