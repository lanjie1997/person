package proxy
import (
	"sync"
	"global"
)
//设备管理类
type DeviceManage struct{
	locker sync.Mutex
	devs   map[string]*Device
}
func (this *DeviceManage) Init() {
	this.devs = make(map[string]*Device)
}

//  设备实例列表  初始化设备状态
func (this *DeviceManage)GetDiscoverList(){
	infoList := global.DbGetOnvifDeviceList()
	for i:=0;i<len(infoList);i++{
		this.AddDev(infoList[i])
	}
}
// 添加设备
func (this *DeviceManage)AddDev(ip string){
	dev := &Device{
		ip:ip,
	}
	dev.Init(18)

	this.locker.Lock()
	defer this.locker.Unlock()

	// 如果设备存在   ，则返回
	if _,ok := this.devs[ip];ok {
		return
	}
	this.devs[ip] = dev
	go dev.SentWhat()
}

//删除设备(此处还需要加关闭tcp套接字的操作和关流的操作)
func (this *DeviceManage)DeletaDev(ip string){
	dev := this.GetDev(ip)
	dev.Close()
	this.locker.Lock()
	defer this.locker.Unlock()
	delete(this.devs,ip)
}

func (this *DeviceManage)GetDev(ip string)*Device{
	return this.devs[ip]
}
