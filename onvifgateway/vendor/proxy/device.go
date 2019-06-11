package proxy

import (
	"fmt"
	"protocol"
	"rtsp"
	"time"
	"wholeally.com/common/logs"
)
//创建一个SOAP类的实例
var soap protocol.SOAP

// 设备类
type Device struct {
	ip           	 string              //IP地址
	state        	 bool                //设备状态
	info         	 protocol.DeviceInfo // 设备信息
	times        	 int                 //超时计数
	closeChan    	 chan bool           //管道
	mediaUri     	 []string			// rtsp地址
	username     	 string				// 账号
	password     	 string				//密码
	mediaservice 	 []string			//媒体服务地址
	videosource		 []string			// VideoSourceToken （某个通道的VideoSourceToken存在=======》此通道在线）
	chanvideo        map[string]string  // 一个视频源匹配一个通道
	sourcefmedia     map[string]string  // 一个videoSourceToken 对应一个 rtsp地址
	rtp          	 map[string]*rtsp.Rtsp // 一个rtsp地址对应一个RTSP类的实例
}

// 关闭设备方法
func (this *Device) Close() {
	close(this.closeChan)
}
//设备初始化
func (this *Device) Init(n int) {
	this.times = 0
	this.state = false
	this.closeChan = make(chan bool)
	this.rtp=make(map[string]*rtsp.Rtsp)
	this.sourcefmedia=make(map[string]string)
	this.chanvideo=make(map[string]string)
}

// 设备方法（根据设备状态动态分配 发心跳还是设备发现）
func (this *Device) SentWhat() {
	for {
		//设备不在线发 设备发现请求
		if this.state == false {
			after := time.NewTicker(time.Second * time.Duration(10))
			defer after.Stop()
			select {
			case <-after.C:
				info, err := protocol.UdpSend(this.ip, 30)
				if nil != err {
					//设备不在线，不用发心跳包
					logs.Debugf("ip:%s ===============>%s", this.ip, err)
				} else {
					//设备在线，定时发心跳包
					this.state = true
					//保存设备信息
					this.info = info.Info
					this.GetStreamUri()
					logs.Debug(fmt.Sprintf("info =========>%+v", info))
					//请求流
					go this.SentRtspRequest()
				}
			//设备被删除时直接返回
			case _, ok := <-this.closeChan:
				if !ok {
					return
				}
			}
			//设备在线时发送心跳包
		} else {
			//每个设备定时10秒发送一次
			after := time.NewTicker(time.Second * 10)
			defer after.Stop()
			for {
				select {
				case <-after.C:
					//连续三次超时=======》离线
					_, err := protocol.UdpSend(this.ip, 1)
					if err != nil {
						this.times += 1
						if this.times == 3 {
							this.state = false
							this.times = 0
							logs.Debugf("%s 设备离线", this.ip)
						}
					} else {
						this.times = 0
						logs.Debugf("%s 设备在线", this.ip)
					}
				//设备删除时直接返回
				case _, ok := <-this.closeChan:
					if !ok {
						return
					}
				}
				//当设备为不在线时 跳出循环到设备发现
				if this.state == false {
					break
				}
			}
		}
	}
}
//获取录像机与IPC的 配置文件令牌
func (this *Device) GetProfileToken()[]string{
		//获取媒体服务地址
		soap.XMLNs = []string{
			`xmlns:tds="http://www.onvif.org/ver10/device/wsdl"`,
			`xmlns:tt="http://www.onvif.org/ver10/schema"`,
		}
		soap.Body = `  <tds:GetCapabilities>
      <tds:Category>All</tds:Category>
    </tds:GetCapabilities>`
		soap.User = "admin"
		soap.Password = "admin123"
		//获取返回内容
		resp, _ := soap.SendRequest(this.info.XAddr)
		//解析成[]interface{}
		mS, _ := resp.ValueForPath("Envelope.Body.GetCapabilitiesResponse.Capabilities.Media.XAddr")
		//遍历转换
		for i:=0 ;i<len(mS);i++{
			this.mediaservice=append(this.mediaservice,mS[i].(string))
		}
		//获取 ProfileToken
		soap.XMLNs = []string{
			`xmlns:trt="http://www.onvif.org/ver10/media/wsdl"`,
			`xmlns:tt="http://www.onvif.org/ver10/schema"`,
		}
		soap.User = "admin"
		soap.Password = "admin123"
		soap.Body = "<trt:GetProfiles />"
		resp, _ = soap.SendRequest(this.mediaservice[0])
		PF, _ := resp.ValueForPath("Envelope.Body.GetProfilesResponse.Profiles.-token")
		profiletoken:=[]string{}
		for  i:=0;i<len(PF);i++{
			profiletoken=append(profiletoken,PF[i].(string))
		}
		logs.Debug(profiletoken)
		return profiletoken
	}
//获取rtsp地址
func (this *Device) GetStreamUri() {
	//获取ProfileToken字符串切片
	profiletoken:=this.GetProfileToken()
	// 创建一个字符串空切片接收数据
	mediaUri:=[]string{}
	// 获取VideoSourceToken的切片
	videosource:=this.GetVideoSource()
	//获取该设设备全部rtsp地址（在线）
	for i:=0;i<len(videosource);i++{
		soap.XMLNs = []string{
			`xmlns:trt="http://www.onvif.org/ver10/media/wsdl"`,
			`xmlns:tt="http://www.onvif.org/ver10/schema"`,
		}
		soap.User = "admin"
		soap.Password = "admin123"
		soap.Body = `<GetStreamUri xmlns="http://www.onvif.org/ver10/media/wsdl">
      <StreamSetup>
        <!-- Attribute Wild card could not be matched. Generated XML may not be valid. -->
        <Stream xmlns="http://www.onvif.org/ver10/schema">RTP-Unicast</Stream>
        <Transport xmlns="http://www.onvif.org/ver10/schema">
          <Protocol>UDP</Protocol>
        </Transport>
      </StreamSetup>
      <ProfileToken>` + profiletoken[i] + `</ProfileToken>
    </GetStreamUri>`
		resp, _:= soap.SendRequest(this.mediaservice[0])//媒体服务地址只有一个
		MU, _ := resp.ValueForPath("Envelope.Body.GetStreamUriResponse.MediaUri.Uri")
		this.mediaUri=append(mediaUri,MU[0].(string))
	}
	for i:=0;i<len(this.mediaUri);i++{
		logs.Debugf("IP：%s====>rtsp地址:%s", this.ip,this.mediaUri[i])
		//一个视频源对应一个rtsp地址
		this.sourcefmedia[this.videosource[i]]=this.mediaUri[i]
	}
}
//获取VideoSourceToken
func (this *Device)GetVideoSource()[]string{
		soap.XMLNs = []string{
			`xmlns:trt="http://www.onvif.org/ver10/media/wsdl"`,
			`xmlns:tt="http://www.onvif.org/ver10/schema"`,
		}
		soap.User = this.info.User
		soap.Password = this.info.Password
		soap.Body = ` <trt:GetVideoSources />`
		soap.User = "admin"
		soap.Password = "admin123"
		//媒体服务地址只有一个
		resp, _ := soap.SendRequest(this.mediaservice[0])
		//IPC只有一个视频源  录像机可能有多个
		vs,_:=resp.ValueForPath("Envelope.Body.GetVideoSourcesResponse.VideoSources.-token")
		videosouce:=[]string{}
		for i:=0;i<len(vs);i++{
			videosouce=append(videosouce,vs[i].(string))
		}
		//保存到设备类
		this.videosource=videosouce
		logs.Debug(this.videosource)
	    return this.videosource
}
//测试 IPC   O K          录像机  O K
func (this *Device) SentRtspRequest() {
	this.rtp[this.mediaUri[0]]=&rtsp.Rtsp{
		Url:this.mediaUri[0],
	}
	if this.info.Type == "dn2:NetworkVideoTransmitter"{
		this.rtp[this.mediaUri[0]].DevSent(this.ip)
	}else {
		this.rtp[this.mediaUri[0]].ReceiverSent(this.ip)
	}
}
