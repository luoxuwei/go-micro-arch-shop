package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-api/user-api/middlewares"
	"shop-api/user-api/router"
)

func InitRouters() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"code":http.StatusOK,
			"success":true,
		})
	})

	ApiGroup := Router.Group("/u/v1")
	ApiGroup.Use(middlewares.Cors())
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	return Router
}
