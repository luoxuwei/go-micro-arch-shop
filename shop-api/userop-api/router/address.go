package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/userop-api/api/address"
	"shop-api/userop-api/middlewares"
)

func InitAddressRouter(Router *gin.RouterGroup) {
	AddressRouter := Router.Group("address")
	{
		AddressRouter.GET("", middlewares.JWTAuth(), address.List)          // 收货地址列表页
		AddressRouter.DELETE("/:id", middlewares.JWTAuth(), address.Delete) // 删除收货地址
		AddressRouter.POST("", middlewares.JWTAuth(), address.New)          //新建收货地址
		AddressRouter.PUT("/:id",middlewares.JWTAuth(), address.Update)     //修改收货地址
	}
}