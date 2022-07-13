package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop-srvs/inventory-srv/global"
	"shop-srvs/inventory-srv/initialize"
	"shop-srvs/inventory-srv/proto"
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
	//TestSetInv(422, 10)
	//TestInvDetail(421)
	//TestSell()
	//TestReback()
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i<10; i++ {
		go TestSell(&wg)
	}
	wg.Wait()
}

func TestSetInv(goodsId, Num int32){
	_, err := inventoryClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
		Num: Num,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("设置库存成功")
}

func TestInvDetail(goodsId int32) {
	rsp, err := inventoryClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Num)
}

func TestSell(wg *sync.WaitGroup) {
	defer wg.Done()
	/*
		要测试事务的效果，
	*/
	_, err := inventoryClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 1},
			//{GoodsId: 422, Num: 1},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("库存扣减成功")
}

func TestReback() {
	_, err := inventoryClient.Reback(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 20},
			{GoodsId: 422, Num: 20},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("归还成功")
}
