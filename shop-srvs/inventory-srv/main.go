package main

import (
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"net"
	"os"
	"os/signal"
	"shop-srvs/inventory-srv/utils/consul"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"shop-srvs/inventory-srv/global"
	"shop-srvs/inventory-srv/handler"
	"shop-srvs/inventory-srv/initialize"
	"shop-srvs/inventory-srv/proto"
	"shop-srvs/inventory-srv/utils"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号")

	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitRedsync()

	zap.S().Info(global.ServerConfig)

	flag.Parse()
	zap.S().Info("ip: ", *IP)
	//为了支持部署多个实例，端口需要动态获取当前可用接口不能写死，但为了调试方便支持设置固定端口
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterInventoryServer(server, &handler.InventoryServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err = register_client.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
	zap.S().Debugf("启动服务器, 端口： %d", *Port)

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//监听库存归还topic
	rocketmqInfo := global.ServerConfig.RocketMqInfo
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{fmt.Sprintf("%s:%d", rocketmqInfo.Host, rocketmqInfo.Port)}),
		consumer.WithGroupName("shop-inventory"),
	)

	if err := c.Subscribe("order_reback", consumer.MessageSelector{},handler.AutoReback); err != nil {
		fmt.Println("读取消息失败")
	}
	_ = c.Start()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = c.Shutdown()
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败")
	} else {
		zap.S().Info("注销成功")
	}

}
