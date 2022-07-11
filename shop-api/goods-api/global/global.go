package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop-api/goods-api/config"
	"shop-api/goods-api/proto"
)

var (
	Trans ut.Translator
	ServerConfig config.ServerConfig
	NacosConfig config.NacosConfig
	GoodsSrvClient proto.GoodsClient
)