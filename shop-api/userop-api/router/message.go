package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/userop-api/api/message"
	"shop-api/userop-api/middlewares"
)

func InitMessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message").Use(middlewares.JWTAuth())
	{
		MessageRouter.GET("", message.List)          // 留言列表页
		MessageRouter.POST("", message.New)       //新建留言
	}
}