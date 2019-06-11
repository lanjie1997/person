package proxy

var devManage DeviceManage


// 进入业务
func Run() {
	//设备管理类初始化
	devManage.Init()
	// 设备管理类获取设备   一获取设备马上发现设备与心跳
	devManage.GetDiscoverList()
	
	


}







