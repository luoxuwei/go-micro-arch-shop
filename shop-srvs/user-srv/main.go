package main

import "shop-srvs/user-srv/initialize"

func main() {
	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()


}
