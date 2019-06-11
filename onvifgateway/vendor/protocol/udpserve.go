package protocol


var udpChannel UdpChannel

type OnvifFindInfo struct {
	Addr string     // udp来源-ip
	Info DeviceInfo // 设备信息
}



type DeviceInfo struct {
	ID       string
	Name     string
	XAddr    string
	Type     string 
	User     string
	Password string
}

// 开始监听
func UdpListen() error {
	udpChannel.Init()
	return udpChannel.Listen()
}

// udp发送
func UdpSend(ip string,n int) (*OnvifFindInfo, error) {
	return udpChannel.Send(ip,n)
}

//
