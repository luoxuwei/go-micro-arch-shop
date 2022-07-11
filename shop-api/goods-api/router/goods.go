package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/goods-api/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup){
	UserRouter := Router.Group("goods")
	{
		UserRouter.GET("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), )
	}
}
