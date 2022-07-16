package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop-api/userop-api/global"
	"shop-api/userop-api/proto"
)

func InitSrvConn(){
	consulInfo := global.ServerConfig.ConsulInfo
	GoodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【商品服务失败】")
	}

	UseropConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UseropSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户操作服务失败】")
	}

	global.AddressSrvClient = proto.NewAddressClient(UseropConn)
    global.MessageSrvClient = proto.NewMessageClient(UseropConn)
    global.UserFavSrvClient = proto.NewUserFavClient(UseropConn)
    global.GoodsSrvClient   = proto.NewGoodsClient(GoodsConn)
}