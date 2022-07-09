package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop-srvs/user-srv/global"
	"shop-srvs/user-srv/initialize"
	"shop-srvs/user-srv/proto"
)

var userSrvClient proto.UserClient

func TestCreateUser(){
	for i := 0; i<10; i++ {
		rsp, err := userSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: fmt.Sprintf("Xuwei%d",i),
			Mobile: fmt.Sprintf("1878222222%d",i),
			PassWord: "admin123",
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Id)
	}
}

func TestGetUserList(){
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 5,
	})
	if err != nil {
		panic(err)
	}
	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.NickName, user.PassWord)
		checkRsp, err := userSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.PassWord,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkRsp.Success)
	}
}

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

	//TestCreateUser()
	TestGetUserList()
}
