package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/user-api/middlewares"

	"shop-api/user-api/api"
)

func InitUserRouter(Router *gin.RouterGroup){
	UserRouter := Router.Group("user")
	{
		UserRouter.POST("register", api.Register)
		UserRouter.POST("pwd_login", api.PassWordLogin)
		UserRouter.GET("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
	}
}

func InitBaseRouter(Router *gin.RouterGroup){
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha)
		BaseRouter.POST("send_sms", api.SendSms)
	}

}