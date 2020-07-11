package 并发

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
)

//示例:Telnet回音服务器 -TCP服务器的基本结构(应用Go语言Socket/goroutine和channel实现简单的telnet协议)
func TestTelnet(t *testing.T) {
	//创建一个程序结束码的channel,退出信号
	exitChan := make(chan int)
	//服务器并发执行
	go Server("127.0.0.1:7001", exitChan)
	//channel阻塞，等待接收返回值
	code := <-exitChan
	//标记程序返回值并退出
	os.Exit(code)
}

//1，接受连接：服务逻辑，传入地址和退出的channel
func Server(address string, exitChan chan int) {
	//根据给定地址进行侦听
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("net.Listen() fail,", err)
		exitChan <- 1
	}
	fmt.Println("侦听器地址:", address)
	//延迟关闭侦听器
	defer l.Close()
	//侦听循环
	for {
		//新连接未到来时,Accept是阻塞的
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("l.Accept fail,", err)
			continue
		}
		//根据连接开启会话，此过程需并发执行
		go handleSession(conn, exitChan)
	}
}

//2,会话处理
func handleSession(conn net.Conn, exitChan chan int) {
	fmt.Println("session started...")
	//创建一个网络连接数据的读取器
	reader := bufio.NewReader(conn)
	//循环接收数据
	for {
		//读取字符串，直到碰到回车返回
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("session closed")
			conn.Close()
			break
		}
		//去掉字符串尾部回车
		str = strings.TrimSpace(str)
		//处理telnet指令
		if !processTelnetCommand(str, exitChan) {
			conn.Close()
			break
		}
		//回音Echo逻辑，发什么数据，则原样返回
		conn.Write([]byte(str + "\r\n"))
	}
}

//3,Telnet命令处理 (Telnet是一种协议,command端可使用telnet命令发起tcp连接)
func processTelnetCommand(str string, exitChan chan int) bool {
	//@close指令表示终止本次session
	if strings.Contains(str, "@close") {
		fmt.Println("session closed")
		return false
	} else if strings.Contains(str, "@shutdown") { //@shutdown:终止服务process
		fmt.Println("server shutdown")
		exitChan <- 0
		return false
	}
	fmt.Println("processTelnetCommand receive input:", str)
	return true
}
