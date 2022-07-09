package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"shop-api/user-api/global"
	"shop-api/user-api/initialize"
	"shop-api/user-api/utils"
	"syscall"
)

func main() {

	initialize.InitLogger()
	initialize.InitConfig()
	Router := initialize.InitRouters()

	viper.AutomaticEnv()
	//为了能支持集群部署，线上环境启动获取端口号，如果是本地开发环境为了方便调试，端口号固定
	debug := viper.GetBool("SHOP_DEBUG")
	if !debug{
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	go func(){
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil{
			zap.S().Panic("启动失败:", err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}