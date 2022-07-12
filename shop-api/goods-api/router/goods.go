package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/goods-api/middlewares"
	"shop-api/goods-api/api/goods"
)

func InitGoodsRouter(Router *gin.RouterGroup){
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("", middlewares.JWTAuth(), goods.List)
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)
		GoodsRouter.GET("/:id", goods.Detail)
		GoodsRouter.PUT("/:id",middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)
		GoodsRouter.PATCH("/:id",middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)
		GoodsRouter.DELETE("/:id",middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete)
	}
}
