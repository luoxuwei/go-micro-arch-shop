package initialize

import (
	"github.com/gin-gonic/gin"
	"shop-api/user-api/router"
)

func InitRouters() *gin.Engine {
	Router := gin.Default()

	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup)

	return Router
}