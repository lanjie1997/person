package main

import (
	"fmt"
	"io"
	"net"
	"os"
)
func main() {
	if len(os.Args) < 3 {
		fmt.Println("参数错误")
		return
	}
	IP := os.Args[1]
	PORT := os.Args[2]
	conn, err := net.Dial("tcp", IP+":"+PORT)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("连接成功！！")
	for {
		go io.Copy(io.Writer(os.Stdin),conn)
		io.Copy(conn,io.Reader(os.Stdout))

	}
}

////接收数据
		//buf := make([]byte, 1024)
		//n, err := conn.Read(buf)
		//if nil != err {
		//	break
		//} else {
		//	str := fmt.Sprintf("%s", string(buf[:n]))
		//	fmt.Println(str)
		//}
		//var input string
		////从键盘和标准输入os.Stdin读取输入
		//fmt.Scanln(&input)
		//if err != nil {
		//	fmt.Println(err.Error())
		//	break
		//}
		//if input == "Q" {
		//	fmt.Println("退出...")
		//	return
		//}
		//_, err = conn.Write([]byte(input)) // 第一个参数是发送的字符数
		//if err != nil {
		//	fmt.Println("发送数据ERROR:", err)
		//	return

