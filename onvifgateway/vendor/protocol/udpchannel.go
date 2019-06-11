package protocol

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
	"mxj"
	"wholeally.com/common/logs"
	"wholeally.com/common/utils"
)


type UdpChannel struct {
	conn      *net.UDPConn
	udpLocker sync.Mutex
	request   Request
}

// 初始化
func (this *UdpChannel) Init() {
	this.request.Init()
}

// 开始监听
func (this *UdpChannel) Listen() error {
	// 获取监听地址
	addr, err := net.ResolveUDPAddr("udp4", ":0")
	if nil != err {
		return err
	}
	logs.Debugf("Listen addr : %s", addr.String())

	// 开始监听
	this.conn, err = net.ListenUDP("udp", addr)
	if nil != err {
		return err
	}
	
	// 启动udp接收协程
	go this.runUdpRecver()

	return nil
}

// udp接收协程
func (this *UdpChannel) runUdpRecver() {
	defer logs.Trace()

	buf := make([]byte, 20*1024)

	for {
		// 读取数据
		n, addr, err := this.conn.ReadFromUDP(buf[:])
		if nil != err {
			logs.Error(err)
			// 出错时需要退出
			os.Exit(-1)
			return
		}

		id, info, err := this.readDisResponse(addr.String(), buf[:n])
		if nil != err {
			logs.Error(err)
			continue
		}

		this.request.Notify(id, info)
	}
}

func (this *UdpChannel) readDisResponse(addr string, buf []byte) (string, *OnvifFindInfo, error) {
	// Parse XML to map
	mapXML, err := mxj.NewMapXml(buf)
	if nil != err {
		return "", nil, err
	}

	// Check if this response is for our request
	id, err := mapXML.ValueForPathString("Envelope.Header.RelatesTo")
	if nil != err {
		return "", nil, err
	}

	// Get device's ID and clean it
	devID, _ := mapXML.ValueForPathString("Envelope.Body.ProbeMatches.ProbeMatch.EndpointReference.Address")
	devID = strings.Replace(devID, "urn:uuid:", "", 1)

	// Get device's name
	devName := ""
	scopes, _ := mapXML.ValueForPathString("Envelope.Body.ProbeMatches.ProbeMatch.Scopes")
	for _, scope := range strings.Split(scopes, " ") {
		if strings.HasPrefix(scope, "onvif://www.onvif.org/name/") {
			devName = strings.Replace(scope, "onvif://www.onvif.org/name/", "", 1)
			devName = strings.Replace(devName, "_", " ", -1)
			break
		}
	}

	// Get device's xAddrs
	xAddrs, _ := mapXML.ValueForPathString("Envelope.Body.ProbeMatches.ProbeMatch.XAddrs")
	listXAddr := strings.Split(xAddrs, " ")
	if 0 >= len(listXAddr) {
		return "", nil, fmt.Errorf("Device does not have any xAddr")
	}

	//Get device's Type
	devType,_ := mapXML.ValueForPathString("Envelope.Body.ProbeMatches.ProbeMatch.Types")

	info := &OnvifFindInfo{
		Addr: addr,
		Info: DeviceInfo{
			ID:    devID,
			Name:  devName,
			Type:  devType,
			XAddr: listXAddr[0],
		},
	}

	return id, info, nil
}

// udp发送
func (this *UdpChannel) Send(ip string,n int) (*OnvifFindInfo, error) {
	// 组装数据包
	id := fmt.Sprint("uuid:" + utils.UUID())
	reqBody := `		
		<?xml version="1.0" encoding="UTF-8"?>
		<e:Envelope
		    xmlns:e="http://www.w3.org/2003/05/soap-envelope"
		    xmlns:w="http://schemas.xmlsoap.org/ws/2004/08/addressing"
		    xmlns:d="http://schemas.xmlsoap.org/ws/2005/04/discovery"
		    xmlns:dn="http://www.onvif.org/ver10/network/wsdl">
		    <e:Header>
		        <w:MessageID>` + id + `</w:MessageID>
		        <w:To e:mustUnderstand="true">urn:schemas-xmlsoap-org:ws:2005:04:discovery</w:To>
		        <w:Action a:mustUnderstand="true">http://schemas.xmlsoap.org/ws/2005/04/discovery/Probe
		        </w:Action>
		    </e:Header>
		    <e:Body>
		        <d:Probe>
		            <d:Types>dn:NetworkVideoTransmitter</d:Types>
		        </d:Probe>
		    </e:Body>
		</e:Envelope>`
	// Clean WS-Discovery message
	reqBody = regexp.MustCompile(`\>\s+\<`).ReplaceAllString(reqBody, "><")
	reqBody = regexp.MustCompile(`\s+`).ReplaceAllString(reqBody, " ")

	//消息对应
	c, err := this.request.NewWaiter(id)
	if nil != err {
		logs.Error(err)
		return nil, err
	}
	
	defer this.request.CloseWaiter(id, c)

	err = this.writeUdp([]byte(reqBody), ip)

	if nil != err {
		return nil, err
	}

	// 设置超时请求
	after := time.NewTimer(time.Second * time.Duration(n))
	defer after.Stop()
	// 等待返回
	select {
	case <-after.C:
		return nil, fmt.Errorf("time out")
	case info := <-c:
		return info, nil
	}
}

// 发送数据
func (this *UdpChannel) writeUdp(body []byte, ip string) error {
	// 接收端地址
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:3702", ip))
	if nil != err {
		return err
	}

	this.udpLocker.Lock()
	defer this.udpLocker.Unlock()

	_, err = this.conn.WriteToUDP(body, addr)

	return err
}
