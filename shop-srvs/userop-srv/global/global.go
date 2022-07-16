package global

import (
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"shop-srvs/userop-srv/config"
)

var (
	DB *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig config.NacosConfig
	RedisSync   *redsync.Redsync
)