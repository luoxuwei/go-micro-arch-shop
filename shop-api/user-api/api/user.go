package api

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop-api/user-api/forms"
	"shop-api/user-api/global"
	"shop-api/user-api/middlewares"
	"shop-api/user-api/models"
	"shop-api/user-api/proto"
	"strings"
	"time"
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

	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, false){
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha":"验证码错误",
		})
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
			if passRsp.Success {
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:             uint(rsp.Id),
					NickName:       rsp.NickName,
					AuthorityId:    uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(), //签名的生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
						Issuer: "imooc",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg":"生成token失败",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"id": rsp.Id,
					"nick_name": rsp.NickName,
					"token": token,
					"expired_at": (time.Now().Unix() + 60*60*24*30)*1000,
				})

			} else {
				c.JSON(http.StatusInternalServerError, map[string]string{
					"msg":"登录失败",
				})
			}
		}
	}
}
