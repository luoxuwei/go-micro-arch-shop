package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop-api/userop-api/config"
	"shop-api/userop-api/proto"
)

var (
	Trans ut.Translator
	ServerConfig config.ServerConfig
	NacosConfig config.NacosConfig
	AddressSrvClient proto.AddressClient
	UserFavSrvClient proto.UserFavClient
	MessageSrvClient proto.MessageClient
	GoodsSrvClient   proto.GoodsClient
)