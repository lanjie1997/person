package config

import (
	"wholeally.com/common/logs"
	"wholeally.com/share/v4/basiconf"
	"wholeally.com/share/v4/dbapi"
)

/*
 * 初始化配置
 */
func init() {
	basiconf.ParseFile("onvifgateway.conf") // 私人配置文件

	// 初始化日志
	basiconf.InitLog("onvifgateway")

	// 初始化监听
	basiconf.InitISPAddr(3702)

	// 实始化dbserver
	dbapi.Init(basiconf.GetDBServer(), basiconf.GetDBSecretKey(), 100)
}

func Dump() {
	basiconf.Dump()
	logs.Info("DB Server: ", basiconf.GetDBServer())
	logs.Info("DB SecretKey: ", basiconf.GetDBSecretKey())
	logs.Info("State Server: ", basiconf.GetStateServer())
	logs.Info("State SecretKey: ", basiconf.GetStateSecretKey())

}
