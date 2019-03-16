package main

import (
	"fmt"
	"net"
	"strings"
)

func (self response) add()string{
	x:=self.state+self.header+self.nothing
	x="HTTP/1.1 200 OK\r\nContent-Type: text/html;charset=utf8\nConnection: Keep-Alive\n\n"
	return x
}
//response结构体
type response struct {
	state string
	header string
	nothing string
	body string
}
var res response
var res1 response
var res2 response
var refuse response
func main() {
	ln,err :=net.Listen("tcp","127.0.0.1:8080")//监听
	if nil != err{
		fmt.Println("error",err)
		return
	}
	for{
		conn,err :=ln.Accept()//建立连接
		if nil != err{
			return
		}
		go handle1(conn)//处理请求和返回响应
	}
}
func handle1(conn net.Conn){
	refuse.header="HTTP/1.1 404 错误\r\n"
	res.body="<!DOCTYPE html><html lang='en'><head><meta charset='UTF-8'><title>hello world</title></head><body><h1 style='text-align: center'> <span style='color: yellow'>h</span><span style='color: red'>e</span><span style='color: blue'>l</span><span style='color: skyblue'>l</span><span style='color: pink'>o</span>     <span style='color: yellow'>w</span><span style='color: red'>o</span><span style='color: blue'>r</span><span style='color: skyblue'>l</span><span style='color: pink'>d</span></h1></body></html>"
	res1.body="<!DOCTYPE html><html><head><title> index</title></head><body><h1>这是主页</h1></body></html>"
	res2.body="<!DOCTYPE html><head><title>登陆</title></head><body><form action='/index ' method='post'>账号：<input type='text'><br>密码：<input type='password' name='username'><br><button type='submit' name='password'>登陆</button></form></body></html>"
	//读取request
	var req []byte
	req = make([]byte, 1024)
	n,err :=conn.Read(req)
	if nil != err{
		return
	}
	request:=string(req[:n])
	x:=strings.Replace(request,"\r\n",",",-1)//\r\n替换成空
	y:=strings.SplitAfterN(x," ",-1)//以" "分割字符串
	if y[0]=="GET "&& y[1]=="/favicon.ico "{
		conn.Write([]byte(refuse.header))

	}
	if y[0]=="GET "&& y[1]=="/ "{
		fmt.Printf("%s  访问hello world页面 \n",conn.RemoteAddr().String())
		conn.Write([]byte(res.add()+res.body))

	}
	if y[0]=="GET "&& y[1]=="/index " {
		fmt.Printf("%s  访问 index 页面  \n",conn.RemoteAddr().String())
		conn.Write([]byte(res1.add()+res1.body))

	}
	if y[0]=="POST "&& y[1]=="/index "{
		fmt.Printf("%s  访问 index 页面  \n",conn.RemoteAddr().String())
		conn.Write([]byte(res1.add()+res1.body))
		fmt.Println(y)


	}
	if y[0]=="GET "&& y[1]=="/login " {
		fmt.Printf("%s  访问 login 页面  \n",conn.RemoteAddr().String())
		conn.Write([]byte(res2.add()+res2.body))

	}
	defer conn.Close()//完成式断开连接
}

