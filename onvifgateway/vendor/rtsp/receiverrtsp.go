package rtsp

import (
	"net"
	"strings"
	"fmt"
	"wholeally.com/common/logs"
	"crypto/md5"
	"encoding/hex"
	"time"
)
// 接收鉴权时nonce的值
var Nonce string
//接收鉴权时realm的值
var Realm string
//接收Session的值
var Session string

// 获取录像机实时流的方法
func (this *Rtsp)ReceiverSent(ip string) {
	//创建tcp连接
	localAddress, err := net.ResolveTCPAddr("tcp4", "172.168.0.200:")
	remoteAddress, err := net.ResolveTCPAddr("tcp4", ip+":554")
	this.Tcpchannel.conn, err = net.DialTCP("tcp4", localAddress, remoteAddress)
	if err != nil {
		return
	}
	//创建Describe请求 发送
	describe:=this.Describe()
	go this.Tcpchannel.conn.Write([]byte(describe))
	logs.Debug(describe)
	//请求顺序控制
	bend:=1
	//读取
	func (){
		for {
			buf := make([]byte ,1500)
			n,err:=this.Tcpchannel.conn.Read(buf)
			if nil != err{
				logs.Debug(err)
				break
			}
			resp:=string(buf[:n])
			logs.Debug(resp)
			if bend == 1{
				// 获取并解析返回的鉴权头部信息     map[realm:***  nonoce:***  其他]
				authdigest:=make(map[string]string)
				authdigest=this.GetNonceRealm(this.Getheader(resp))
				//获取nonce与realm的值
				Nonce = authdigest["nonce"]
				Realm = authdigest["realm"]
				//创建带鉴权信息的Describe请求  发送
				describeagain:=this.RDescribe(Nonce,Realm)
				go this.Tcpchannel.conn.Write([]byte(describeagain))
				logs.Debug(describeagain)
			}
			if bend == 2{
				//构建Setup请求 并加入鉴权信息
				this.SETUP()
				this.header+="\r\n"+this.Auth(Nonce,Realm)
				//获取video 轨道ID 并加入URL
				body:=GetBody(resp)
				trackID:=GetTrackID(body)
				this.Url+=fmt.Sprintf("/%s",trackID)
				//创建setup请求 发送
				setup:=this.CreateRtspRequest()
				go this.Tcpchannel.conn.Write([]byte(setup))
			}
			if bend == 3{
				//解析并获取session的值
				auth:=make(map[string]string)
				auth = this.Getheader(resp)
				Session=auth["Session"]
				//构建play请求并加入鉴权信息  发送
				this.PLAY()
				this.header+="\r\n"+"Session:"+Session+"\r\n"+this.Auth(Nonce,Realm)
				play:=this.CreateRtspRequest()
				go this.Tcpchannel.conn.Write([]byte(play))
			}
			if bend ==4{
				//构建teardown请求
				this.TEARDOWN()
				this.header+="\r\n"+this.Auth(Nonce,Realm)
				teardown:=this.CreateRtspRequest()
				// 创建20秒的定时器 （延迟20秒关流）
				after := time.NewTimer(time.Second * 20)
				defer after.Stop()
				<-after.C
				go this.Tcpchannel.conn.Write([]byte(teardown))
			}
			bend+=1
		}
	}()
	defer this.Tcpchannel.conn.Close()
}
//带鉴权信息的Describe请求
func (this *Rtsp)RDescribe(nonce,realm string)string{
	this.method="DESCRIBE"
	this.header="CSeq:1"+"\r\n"+"Accept:application/sdp"+"\r\n"+"User-Agent: NKPlayer-VSPlayer1.0"+"\r\n"+this.Auth(nonce ,realm)
	return this.CreateRtspRequest()
}
// response算法     （为rtsp密码为明文的算法）       此处账号密码需改成从数据库获取
func (this *Rtsp)GetResponse(nonce ,realm string)string{
	return Md5(Md5("admin:"+realm+":admin123")+":"+nonce+":"+Md5(this.method+":"+this.Url))
}
//鉴权头信息的方法
func (this *Rtsp)Auth(nonce ,realm string)string{
	response:=this.GetResponse(nonce,realm)
	return fmt.Sprintf(`Authorization: Digest username="admin",uri="%s",nonce="%s",response="%s",realm="%s"`,this.Url,nonce,response,realm)
}
// 当返回401时   对返回内容的头部进行解析====》map
func (this *Rtsp)Getheader(resp string)map[string]string{
		x:=strings.Split(resp,"\r\n")
		z:=make(map[string]string)
		for i:=0;i<len(x[1:len(x)-3]);i++{
			y:=strings.Split(x[1:][i],":")
			z[y[0]]=y[1]
		}
		return z
}
// 对鉴权行的解析  =======》map
func (this *Rtsp)GetNonceRealm(z map[string]string)map[string]string{
	var x string
	for k,v:=range z {
		if k == "WWW-Authenticate" {
			x=v
		}
	}
	x1 := strings.Split(x,",")
	x1[0] = strings.ReplaceAll(x1[0],"Digest ","")
	a := make(map[string]string)
	for i:=0;i<len(x1);i++{
		x1[i]=strings.ReplaceAll(x1[i],"\"","")
		y:=strings.Split(x1[i],"=")
		y[0]=strings.ReplaceAll(y[0]," ","")
		a[y[0]]=y[1]
	}
	return a
}
// md5加密算法
func Md5(src string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(src))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}


