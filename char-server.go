package main

import (
	"fmt"
	"net"
	"sync"
)

// sync.Map
var locker sync.Mutex
var sockes map[string]net.Conn

func main() {
	sockes = make(map[string]net.Conn)
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

		// fmt.Println("=== start")
		go connClient(conn)
		// fmt.Println("--- end")
	}

	// time.Sleep(100 * time.Minute)
}

func connClient(conn net.Conn) {
	defer conn.Close()
	// 保存记录
	locker.Lock()
	sockes[conn.RemoteAddr().String()] = conn
	locker.Unlock()

	// 退出时删除记录
	defer func() {
		locker.Lock()
		delete(sockes, conn.RemoteAddr().String())
		locker.Unlock()
	}()

	// 数据发送
	conn.Write([]byte("欢迎进入聊天室\n"))
	// 数据接收
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if nil != err {
			fmt.Println("socket close: ", conn.RemoteAddr().String())
			break
		} else {
			str := fmt.Sprintf("%s say:%s\n", conn.RemoteAddr().String(), string(buf[:n]))
			// fmt.Println("socket recv: ", string(buf[:n]))
			locker.Lock()
			for _, val := range sockes {
				if conn != val {
					fmt.Println("===========", val.RemoteAddr().String())
					val.Write([]byte(str))
				}
			}

			locker.Unlock()
		}
	}
}
