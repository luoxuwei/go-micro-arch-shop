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
)

var orderClient proto.OrderClient

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
	orderClient = proto.NewOrderClient(Conn)

	//TestCreateCartItem(1,1,422)
	TestCartItemList(1)
	//TestUpdateCartItem(1, 422)
	//TestCreateOrder()
}

//向购物车中添加商品
func TestCreateCartItem(userId, nums, goodsId int32){
	rsp, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId: userId,
		Nums: nums,
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

func TestCartItemList(userId int32) {
	rsp, err := orderClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: userId,
	})
	if err != nil {
		panic(err)
	}
	for _, item := range rsp.Data {
		fmt.Println(item.Id, item.GoodsId, item.Nums)
	}
}

func TestUpdateCartItem(user, goods int32) {
	_, err := orderClient.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		UserId: user,
		GoodsId: goods,
		Checked: true,
	})
	if err != nil {
		panic(err)
	}
}

func TestCreateOrder() {
	_, err := orderClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  1,
		Address: "北京市",
		Name:    "xuwei",
		Mobile:  "18787878787",
		Post:    "请尽快发货",
	})
	if err != nil {
		panic(err)
	}
}

