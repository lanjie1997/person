package main

import (
	"fmt"
	"net"
)

//var locker sync.Mutex
//var sockes map[string]net.Conn

func main() {
	//sockes = make(map[string]net.Conn)
	// 监听端口
	listener, err := net.Listen("tcp", ":9090")
	if nil != err {
		panic(err)
	}
	defer listener.Close()
	for {
		// 等待连接
		conn, err := listener.Accept()
		if nil != err {
			panic(err)
		}
		go connClient(conn)
	}
}
func connClient(conn net.Conn) {
	en := make(map[string]string)
	//翻译的数据
	en["like\r\n"] = "喜欢"
	en["no\r\n"] = "不"
	en["hi\r\n"]= "你好"
	en["world\r\n"] = "世界"
	en["database\r\n"] = "数据库"
	en["great\r\n"] = "棒"
	en["computer\r\n"] = "电脑"
	defer conn.Close()
	//// 保存记录
	//locker.Lock()
	//sockes[conn.RemoteAddr().String()] = conn
	//locker.Unlock()
	//// 退出时删除记录
	//defer func() {
	//	locker.Lock()
	//	delete(sockes, conn.RemoteAddr().String())
	//	locker.Unlock()
	//}()
	// 数据发送
	conn.Write([]byte("翻译 \n"))
	// 数据接收
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if nil != err {
			fmt.Println("socket close: ", conn.RemoteAddr().String())
			break
		} else {
			str := fmt.Sprintf("中文：%s\n", (en[string(buf[:n])]))
			conn.Write([]byte(str))
			//locker.Lock()
			//for _, val := range sockes {
			//	if conn == val {
			//		fmt.Println("===========", val.RemoteAddr().String())
			//		val.Write([]byte(str))
			//	}
			//}
			//locker.Unlock()
		}
	}
}
