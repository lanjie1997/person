package serve

import (
	"fmt"
	"os"
	"os/signal"
	"protocol"

	"proxy"

	"wholeally.com/common/logs"
	"wholeally.com/share/v4/basiconf"
	"wholeally.com/share/v4/regclient"
	"wholeally.com/share/v4/service"
	"wholeally.com/share/v4/stateapi"
)

func Main() {
	// 连接配置中心
	service.QuickStart(regclient.CODE_ONVIFGATEWAY, nil, nil, msqNotify)

	// 初始化状态服务器(无用，下个版本去掉)
	stateapi.SetGatewayInfo(fmt.Sprint(basiconf.GetClientID()), basiconf.GetServerID())
	stateapi.Init(basiconf.GetStateServer(), basiconf.GetStateSecretKey(), 1)

	// 监听设备搜索udp
	err := protocol.UdpListen()
	if nil != err {
		logs.Error(err)
		os.Exit(-1)
		return
	}
	

	// 进入业务
	go proxy.Run()

	// 等待程序退出
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
