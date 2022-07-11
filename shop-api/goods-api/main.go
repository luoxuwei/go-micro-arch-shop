package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"shop-api/goods-api/global"
	"shop-api/goods-api/initialize"
	"shop-api/goods-api/utils"
	"shop-api/goods-api/utils/consul"

	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitSrvConn()
	Router := initialize.InitRouters()
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	viper.AutomaticEnv()
	//为了能支持集群部署，线上环境启动获取端口号，如果是本地开发环境为了方便调试，端口号固定
	debug := viper.GetBool("SHOP_DEBUG")
	if !debug{
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	consul_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := consul_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
	zap.S().Debugf("启动服务器, 端口： %d", global.ServerConfig.Port)

	go func(){
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil{
			zap.S().Panic("启动失败:", err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := consul_client.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功:")
	}

}