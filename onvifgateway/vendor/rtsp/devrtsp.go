package rtsp

import (
	"net"
	"time"
	"sync"
	"wholeally.com/common/logs"
	"fmt"
	"strings"
	"github.com/pixelbender/go-sdp/sdp"
)
type Rtsp struct {
	method 		     string  //  方法名
	Url              string  //   地址
	header           string  //   头部
	body    	     string  //  一般不使用
	Tcpchannel       TcpChannel //TCP套接字
	username         string     //用户名
	password         string		//密码
}

type TcpChannel struct {
	conn      *net.TCPConn
	tcpLocker sync.Mutex
}
const version  = "RTSP/1.0"
// IPC 请求实时流 与录像机大概一致
func (this *Rtsp)DevSent(ip string) {
	localAddress, err := net.ResolveTCPAddr("tcp4",  "172.168.0.200:")
	remoteAddress, err := net.ResolveTCPAddr("tcp4",ip+":554")
	this.Tcpchannel.conn, err = net.DialTCP("tcp4", localAddress, remoteAddress)
	if err != nil {
		logs.Debug(err)
	}
	describe:=this.Describe()
	bend := 1
	go this.Tcpchannel.conn.Write([]byte(describe))
	logs.Debug(describe)
	func() {
		for {
			buf := make([]byte, 1500)
			n, err := this.Tcpchannel.conn.Read([]byte(buf))
			if err != nil {
				break
			}
			resp:=string(buf[:n])
			logs.Debug(resp)
			if bend == 1{
				if this.GetStateCode(resp) != "200"{
					break
				}
				body:=GetBody(resp)
				trackID:=GetTrackID(body)
				this.SETUP()
				this.Url+=fmt.Sprintf("/%s",trackID)
				setup:=this.CreateRtspRequest()
				go this.Tcpchannel.conn.Write([]byte(setup))
				logs.Debug(setup)
			}
			if bend == 2 {
				this.PLAY()
				play:=this.CreateRtspRequest()
				go this.Tcpchannel.conn.Write([]byte(play))
				logs.Debug(play)
			}
			if bend == 3 {
				this.TEARDOWN()
				teardown:=this.CreateRtspRequest()
				after := time.NewTimer(time.Second * 20)
				defer after.Stop()
				<-after.C
				go this.Tcpchannel.conn.Write([]byte(teardown))
			}
			bend += 1
		}
	}()
	defer this.Tcpchannel.conn.Close()
}

func (this * Rtsp)CreateRtspRequest()string{
	request:=this.method+" "+this.Url+" "+version +"\r\n"
	request+=this.header
	request+="\r\n\r\n"
	return request
}
func (this * Rtsp)Describe()string{
	this.method="DESCRIBE"
	this.header="CSeq:1"+"\r\n"+"Accept:application/sdp"+"\r\n"+"User-Agent: NKPlayer-VSPlayer1.0"
	return this.CreateRtspRequest()
}
func (this *Rtsp)SETUP(){
	this.method="SETUP"
	this.header="CSeq:2"+"\r\n"+"User-Agent: NKPlayer_VSPlayer1.0"+"\r\n"+"Transport:RTP/AVP/TCP;unicast;interleaved=0-1"
}
func (this *Rtsp)PLAY(){
	this.method="PLAY"
	this.header="CSeq:3"+"\r\n"+"Range:npt=0.000-"+"\r\n"+"User-Agent:  NKPlayer_VSPlayer1.0"
}
func (this *Rtsp)TEARDOWN(){
	this.method="TEARDOWN"
	this.header="CSeq:4"+"\r\n"+"User-Agent:  NKPlayer_VSPlayer1.0"
}

func (this *Rtsp)GetStateCode(s string)string{
	a:=strings.SplitAfterN(s,"\r\n",-1)
	a1:=strings.SplitAfterN(a[0]," ",-1)
	a2:=strings.ReplaceAll(a1[1]," ","")
	return a2
}

func GetBody(resp string)string{
	b:=strings.Split(resp,"\r\n\r\n")
	if len(b)<2{
		return ""
	}
	return b[1]
}
func GetTrackID(body string)string{
	var val string
	session,err := sdp.ParseString(body)
	if nil != err || nil == session {
		logs.Error(err.Error())
		return ""
	}
	if 0 < len(session.Media) {
		if session.Media[0].Type=="video"{
			val=session.Media[0].Attributes.Get("control").Value
		}
	}
	return val
}

