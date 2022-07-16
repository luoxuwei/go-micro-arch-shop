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
	Conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UseropSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	global.AddressSrvClient = proto.NewAddressClient(Conn)
    global.MessageSrvClient = proto.NewMessageClient(Conn)
    global.UserFavSrvClient = proto.NewUserFavClient(Conn)
}