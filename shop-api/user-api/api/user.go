package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop-api/user-api/forms"
	"shop-api/user-api/global"
	"shop-api/user-api/proto"
	"strings"
)

//去掉struct名称
func removeTopStruct(fileds map[string]string) map[string]string{
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleValidatorError(c *gin.Context, err error){
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg":err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func PassWordLogin(c *gin.Context)  {
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	if rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg":"用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"msg":"登录失败",
				})
			}
			return
		}

		c.JSON(http.StatusInternalServerError, map[string]string{
			"msg":"登录失败",
		})
	} else {
		if passRsp, pasErr := global.UserSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password: passwordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		}); pasErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"msg":"登录失败",
			})
		} else {
			if (passRsp.Success) {
				c.JSON(http.StatusOK, gin.H{
					"msg":"登录成功",
				})
			} else {
				c.JSON(http.StatusInternalServerError, map[string]string{
					"msg":"登录失败",
				})
			}
		}
	}
}
