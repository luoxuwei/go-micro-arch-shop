package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/goods-api/api/banners"
	"shop-api/goods-api/middlewares"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners")
	{
		BannerRouter.GET("", banners.List)          // 轮播图列表页
		BannerRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.Delete) // 删除轮播图
		BannerRouter.POST("",  middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.New)       //新建轮播图
		BannerRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.Update) //修改轮播图信息
	}
}