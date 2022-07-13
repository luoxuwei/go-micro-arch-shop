package middlewares

import (
	"github.com/gin-gonic/gin"
	"shop-api/oss-api/models"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc{
	return func(ctx *gin.Context){
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg":"无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}

}
