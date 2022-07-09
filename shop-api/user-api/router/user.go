package router

import (
	"github.com/gin-gonic/gin"

	"shop-api/user-api/api"
)

func InitUserRouter(Router *gin.RouterGroup){
	UserRouter := Router.Group("user")
	{
		UserRouter.POST("pwd_login", api.PassWordLogin)
	}
}

func InitBaseRouter(Router *gin.RouterGroup){
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha)
	}

}