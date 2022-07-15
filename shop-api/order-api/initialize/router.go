package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-api/order-api/middlewares"
	"shop-api/order-api/router"
)

func InitRouters() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"code":http.StatusOK,
			"success":true,
		})
	})

	ApiGroup := Router.Group("/g/v1")
	ApiGroup.Use(middlewares.Cors())
	router.InitOrderRouter(ApiGroup)
	return Router
}
