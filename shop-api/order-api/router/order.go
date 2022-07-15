package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/order-api/api/shop_cart"
	"shop-api/order-api/middlewares"
)

func InitOrderRouter(Router *gin.RouterGroup){
	//OrderRouter := Router.Group("orders")
	//{
	//}
}

func InitShopCartRouter(Router *gin.RouterGroup){
	ShopCartsRouter := Router.Group("shopcarts").Use(middlewares.JWTAuth())
	{
		ShopCartsRouter.GET("", shop_cart.List) //购物车列表
		ShopCartsRouter.POST("", shop_cart.New) //添加商品到购物车
	}
}
