package main

import (
	"api"
	"config"
	"serve"
	"wholeally.com/common.v4/checking"
	"wholeally.com/share/v4/regclient"
)

func main() {
	config.Dump()

	// 启动性参监测
	checking.RegisterChecking(regclient.CHECK_ONVIFGATEWAY, nil)

	// 初始化api路由
	api.InitApiRouter()

	serve.Main()
}
