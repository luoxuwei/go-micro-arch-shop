package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop-srvs/goods-srv/global"
	"shop-srvs/goods-srv/initialize"
	"shop-srvs/goods-srv/proto"
)

var goodsClient proto.GoodsClient

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
	goodsClient = proto.NewGoodsClient(Conn)
	//TestGetBrandList()
	//TestGetCategoryBrandList()
	//TestGetCategoryList()
	//TestGetSubCategoryList()
	TestGetGoodsList()
}

func TestGetBrandList(){
	rsp, err := goodsClient.BrandList(context.Background(), &proto.BrandFilterRequest{
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)
	}
}

func TestGetCategoryBrandList(){
	rsp, err := goodsClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.Data)
}

func TestGetCategoryList(){
	rsp, err := goodsClient.GetAllCategorysList(context.Background(), &empty.Empty{
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.JsonData)
}

func TestGetSubCategoryList(){
	rsp, err := goodsClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id:       135487,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.SubCategorys)
}

func TestGetGoodsList(){
	rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		TopCategory: 130361,
		PriceMin: 90,
		//KeyWords: "深海速冻",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, good := range rsp.Data {
		fmt.Println(good.Name, good.ShopPrice)
	}
}

func TestBatchGetGoods(){
	rsp, err := goodsClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: []int32{421, 422, 423},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, good := range rsp.Data {
		fmt.Println(good.Name, good.ShopPrice)
	}
}

func TestGetGoodsDetail(){
	rsp, err := goodsClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: 421,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Name)
	fmt.Println(rsp.DescImages)
}