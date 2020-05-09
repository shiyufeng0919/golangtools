package scp_copy_files_test

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"
)

/*
 模拟scp在两台主机间传送文件，上传文件为.zip格式压缩包 add by syf 2020.5.3
 file_scp类似于scp（安全复制，远程文件复制程序），但是使用golang在网络上的主机之间复制文件。
 参考：https://github.com/gerryyang/goinaction/tree/master/src/file_scp
 客户端：以本机作为客户端，向服务端(部署在虚拟机)发送.zip文件
 项目需求：jdchain支持动态入网及动态加入账本，需拷备数据库目录rocksdb(rocksdb0/rocksdb1...)由于账本文件较大，且可能会跨集群。
*/
const (
	VERSION    = "1.0.0"
	REQ_HEADER = "10001" // 0X2711
)

//测试客户端所发送的文件名称
var testFileName = "/tmp/shared/rocksdb0.zip"

//模拟客户端向服务端发送文件
func TestScpSendFile(t *testing.T) {
	//GOMAXPROCS设置可以执行的最大CPU数
	runtime.GOMAXPROCS(runtime.NumCPU())
	//操作命令行
	file, blk, job, service := handleCommandLine()
	fin, err := os.Open(file)
	checkError(err)
	defer fin.Close()
	buf := make([]byte, blk)
	fin_info, _ := fin.Stat()
	fmt.Printf("file[%s] size[%d]\n", fin_info.Name(), fin_info.Size())
	lock_chan_len := fin_info.Size()/blk + 1
	lock_chan := make(chan bool, lock_chan_len)
	lock_job_chan := make(chan bool, job)
	start := time.Now()
	var cid int = 0
	work(fin, buf, service, lock_chan, lock_job_chan, &cid)
	// wait for all goroutine completion
	for i := 0; i < cid; i++ {
		<-lock_chan
	}
	elapsed := 1000000 * time.Since(start).Seconds()
	fmt.Println("time elapsed: ", elapsed, "us")
}

func printVersion() {
	fmt.Println("file send" + VERSION + " by syf")
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

//操作命令行
func handleCommandLine() (file string, blk int64, job int64, service string) {
	var f = flag.String("f", testFileName, "file")
	var b = flag.Int64("b", 1024, "blk")
	var j = flag.Int64("j", 1, "job")
	var s = flag.String("s", "127.0.0.1:9001", "host:ip")
	var printVer bool
	flag.BoolVar(&printVer, "v", false, "print version")
	flag.Parse()
	if printVer {
		printVersion()
		os.Exit(0)
	}
	fmt.Println("file:", *f)
	fmt.Println("blk:", *b)
	fmt.Println("job:", *j)
	fmt.Println("host:ip:", *s)
	return *f, *b, *j, *s
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8) // int64 is 8 byte
	binary.LittleEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.LittleEndian.Uint64(buf))
}

func proc(req *string, reqlen int, cid int, offset int64, service string, lock_chan chan bool, lock_job_chan chan bool) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	var cid_name string
	cid_name = fmt.Sprintf("%d", cid)
	send(tcpAddr, req, reqlen, cid_name, offset, lock_chan)
	// flow control, refer to job
	<-lock_job_chan
}

func send(tcpAddr *net.TCPAddr, req *string, reqlen int, cid_name string, offset int64, lock_chan chan bool) {
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DialTCP: cid_name[%s] err[%s]\n", cid_name, err.Error())
		lock_chan <- true
		return
	}
	var cid int
	fmt.Sscanf(cid_name, "%d", &cid)
	req_header, _ := strconv.ParseInt(REQ_HEADER, 10, 64)
	var reqlen_total int64 = int64(3*8 + reqlen)
	var bytes_buf bytes.Buffer
	bytes_buf.Write(Int64ToBytes(reqlen_total))
	bytes_buf.Write(Int64ToBytes(req_header))
	bytes_buf.Write(Int64ToBytes(offset))
	bytes_buf.Write([]byte(*req))
	// TODO here may be fail, and need to retry
	wcnt, werr := conn.Write(bytes_buf.Bytes())
	if werr != nil {
		fmt.Fprintf(os.Stderr, "Write: cid_name[%s] err[%s]\n", cid_name, err.Error())
		lock_chan <- true
		return
	}
	fmt.Printf("Write: ok cid_name[%s] wcnt[%d]\n", cid_name, wcnt)
	var ans_buf [4]byte
	_, err = conn.Read(ans_buf[0:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read: cid_name[%s] err[%s]\n", cid_name, err.Error())
		lock_chan <- true
		return
	}
	fmt.Fprintf(os.Stdout, "Read: cid_name[%s] ans_buf[%s]\n", cid_name, ans_buf)
	lock_chan <- true
	conn.Close()
}

func work(fin *os.File, buf []byte, service string, lock_chan chan bool, lock_job_chan chan bool, cid *int) {
	var offset int64 = 0
	for {
		//time.Sleep(time.Duration(1) * time.Second)
		cnt, err := fin.ReadAt(buf, offset)
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "ReadAt: fatal error[%s]\n", err.Error())
			os.Exit(1)
		} else if err == io.EOF {
			if cnt == 0 {
				fmt.Println("At end of file, that error is io.EOF and cnt is 0")
				break
			}
			// flow control, refer to job
			lock_job_chan <- true
			//fmt.Printf("cid[%d] offset[%d] read bytes[%d] buf[%q]\n", *cid, offset, cnt, buf[:cnt])
			fmt.Printf("cid[%d] offset[%d] read bytes[%d] left\n", *cid, offset, cnt)
			var req string = string(buf[:cnt])
			go proc(&req, cnt, *cid, offset, service, lock_chan, lock_job_chan)
			offset += int64(cnt)
			*cid++
			break
		} else {
			// flow control, refer to job
			lock_job_chan <- true
			//fmt.Printf("cid[%d] offset[%d] read bytes[%d] buf[%q]\n", *cid, offset, cnt, buf[:cnt])
			fmt.Printf("cid[%d] offset[%d] read bytes[%d] all\n", *cid, offset, cnt)
			var req string = string(buf[:cnt])
			go proc(&req, cnt, *cid, offset, service, lock_chan, lock_job_chan)
			offset += int64(len(buf))
			*cid++
		}
	}
}
