package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop-api/oss-api/config"
)

var (
	Trans ut.Translator

	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

)
