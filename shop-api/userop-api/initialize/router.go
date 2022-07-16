package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-api/userop-api/middlewares"
	"shop-api/userop-api/router"
)

func InitRouters() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"code":http.StatusOK,
			"success":true,
		})
	})

	ApiGroup := Router.Group("/up/v1")
	ApiGroup.Use(middlewares.Cors())
	router.InitMessageRouter(ApiGroup)
	router.InitAddressRouter(ApiGroup)
	return Router
}
