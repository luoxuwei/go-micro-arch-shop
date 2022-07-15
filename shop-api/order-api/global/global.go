package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop-api/order-api/config"
	"shop-api/order-api/proto"
)

var (
	Trans ut.Translator
	ServerConfig config.ServerConfig
	NacosConfig config.NacosConfig
	OrderSrvClient proto.OrderClient
	GoodsSrvClient proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)