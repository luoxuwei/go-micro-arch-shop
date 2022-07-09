package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop-api/user-api/config"
	"shop-api/user-api/proto"
)

var (
	Trans ut.Translator
	ServerConfig config.ServerConfig
	NacosConfig config.NacosConfig
	UserSrvClient proto.UserClient
)