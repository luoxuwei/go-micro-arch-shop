package global

import (
	"shop-api/user-api/config"
	"shop-api/user-api/proto"
)

var (
	ServerConfig config.ServerConfig
	NacosConfig config.NacosConfig
	UserSrvClient proto.UserClient
)