package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop-srvs/user-srv/global"
	"shop-srvs/user-srv/initialize"
	"shop-srvs/user-srv/proto"
)

var userSrvClient proto.UserClient
func main() {
    initialize.InitConfig()

	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	userSrvHost := ""
	userSrvPort := 0
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.Name))

	if err != nil {
		panic(err)
	}
	for _, value := range data{
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == ""{
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
	}

	userSrvClient = proto.NewUserClient(userConn)

}
