package forms

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Type uint `form:"type" json:"type" binding:"required,oneof=1 2"`
	//注册发送短信验证码和动态验证码登录发送验证码
}

type PassWordLoginForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Captcha string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}