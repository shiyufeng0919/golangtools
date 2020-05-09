package socket_transport_files_test

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"testing"
	"time"
)

/*
 模拟golang socket断点续传大文件,上传文件为.zip压缩包 add by syf 2020.5.1
 参见博文：https://blog.csdn.net/luckytanggu/article/details/79830493
 示例描述：client_zip_test运行在本机,以本机作为客户端，向服务端(虚拟机)应用socket传输.zip文件(server_zip_test部署在虚拟机)。
 项目需求：jdchain支持动态入网及动态加入账本，需拷备数据库目录rocksdb(rocksdb0/rocksdb1...)由于账本文件较大，且可能会跨集群。
*/
var testFileName = "/tmp/shared/rocksdb0.zip" //测试读取客户端的文件名称

//测试客户端读取.zip文件，并发送到服务端
func TestClientSendZipFile(t *testing.T) {
	//设置拨号超时时间为30秒
	conn, err := net.DialTimeout("tcp", "127.0.0.1:35000", time.Second*30)
	if err != nil {
		log.Fatalf("client dial faild: %s\n", err)
	}
	clientConn(conn)
}

//获取服务端发送的消息
func clientRead(conn net.Conn) int {
	buf := make([]byte, 5)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("receive server info faild: %s\n", err)
	}
	//string conver int
	off, err := strconv.Atoi(string(buf[:n]))
	if err != nil {
		log.Fatalf("string conver int faild: %s\n", err)
	}
	return off
}

//客户端连接服务端发送文件
func clientConn(conn net.Conn) {
	defer conn.Close()
	//发送"start-->"消息通知服务端，我要开始发送文件内容了
	//你赶紧告诉我你那边已经接收了多少内容,我从你已经接收的内容处开始继续发送
	clientWrite(conn, []byte("start-->"))
	off := clientRead(conn)
	//发送文件的内容
	fp, err := os.OpenFile(testFileName, os.O_RDONLY, 0755)
	if err != nil {
		log.Fatalf("open file faild: %s\n", err)
	}
	defer fp.Close()
	//设置从哪里开始读取文件内容
	_, err = fp.Seek(int64(off), 0)
	if err != nil {
		log.Fatalf("set file seek faild: %s\n", err)
	}
	log.Printf("read file at seek: %d\n", off)
	for {
		//每次发送100个字节大小的内容
		data := make([]byte, 100)
		n, err := fp.Read(data)
		if err != nil {
			if err == io.EOF {
				//如果已经读取完文件内容
				//就发送'<--end'消息通知服务端，文件内容发送完了
				time.Sleep(time.Second * 1)
				clientWrite(conn, []byte("<--end"))
				log.Println("send all content, now quit")
				break
			}
			log.Fatalf("read file err: %s\n", err)
		}
		// 发送文件内容到服务端
		clientWrite(conn, data[:n])
	}
}

//发送内容到服务端
func clientWrite(conn net.Conn, data []byte) {
	_, err := conn.Write(data)
	if err != nil {
		log.Fatalf("send 【%s】 content faild: %s\n", string(data), err)
	}
	log.Printf("send 【%s】 content success\n", string(data))
}
