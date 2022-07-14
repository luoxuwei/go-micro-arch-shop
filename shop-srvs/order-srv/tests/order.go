package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop-srvs/order-srv/global"
	"shop-srvs/order-srv/initialize"
	"shop-srvs/order-srv/proto"
	"sync"
)

var inventoryClient proto.InventoryClient

func main() {
    initialize.InitConfig()

	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	Host := ""
	Port := 0
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.Name))

	if err != nil {
		panic(err)
	}
	for _, value := range data{
		Host = value.Address
		Port = value.Port
		break
	}
	if Host == ""{
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}

	Conn, err := grpc.Dial(fmt.Sprintf("%s:%d", Host, Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
	}
	inventoryClient = proto.NewInventoryClient(Conn)

}
