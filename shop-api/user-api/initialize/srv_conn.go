package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"shop-api/user-api/global"
	"shop-api/user-api/proto"
)

func InitSrvConn(){
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}