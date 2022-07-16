package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/order-api/api/order"
	"shop-api/order-api/api/pay"
	"shop-api/order-api/api/shop_cart"
	"shop-api/order-api/middlewares"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("orders").Use(middlewares.JWTAuth())
	{
		OrderRouter.GET("", order.List) //订单列表
		OrderRouter.GET("/:id", order.Detail)  // 订单详情
		OrderRouter.POST("",  order.New)  // 新建订单
	}

	PayRouter := Router.Group("pay")
	{
		PayRouter.POST("alipay/notify", pay.Notify)
	}
}

func InitShopCartRouter(Router *gin.RouterGroup){
	ShopCartsRouter := Router.Group("shopcarts").Use(middlewares.JWTAuth())
	{
		ShopCartsRouter.GET("", shop_cart.List) //购物车列表
		ShopCartsRouter.POST("", shop_cart.New) //添加商品到购物车
		ShopCartsRouter.PATCH("/:id", shop_cart.Update) //修改条目
		ShopCartsRouter.DELETE("/:id", shop_cart.Delete) //删除条目
	}
}
