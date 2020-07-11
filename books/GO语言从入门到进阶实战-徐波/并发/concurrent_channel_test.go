package 并发

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

/*
channel是一个引用对象，和map类似。map在没有任何外部引用时，Go程序在运行时（runtime）会自动对内存进行垃圾回收（Garbage Collection,GC）
无缓冲channel和有缓冲channel:
无缓冲channel:(同步)
	无缓冲通道保证收发过程同步
有缓冲channel:(异步)
	带缓冲通道在发送时无需等待接收方接收即可完成发送过程，并且不会发生阻塞，只有当存储空间满时才会发生阻塞
*/
type TestChannel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

//示例1:创建普通类型的channel
func TestOneCreateChannel(t *testing.T) {
	ch := make(chan int)
	ch1 := make(chan string)
	go func() {
		fmt.Println("start goroutine...")
		ch <- 1
		ch1 <- "yufeng"
		fmt.Println("exit goroutine...")
	}()
	fmt.Println("wait goroutine...")
	//等待匿名goroutine
	<-ch //此种方式匆略从通道接收到的数据
	data := <-ch1
	fmt.Println("all done...receive data:", data)
}

//示例2:创建interface类型的channel
func TestTwoCreateChannel(t *testing.T) {
	//make创建一个数据类型为interface的channel
	ch1 := make(chan interface{})
	go func() {
		fmt.Println("start goroutine...")
		//向channel中写入数据，此时发送将持续阻塞状态，直到数据被接收
		ch1 <- 0
		ch1 <- "kaixin"
		ch1 <- "yufeng"
		fmt.Println("exit goroutine...")
	}()
	fmt.Println("wait goroutine...")

	data1 := <-ch1 //此时仅会从channel读取0,不会执行上述exit goroutine...直到所有数据均从ch1中被取出
	fmt.Println("from channel read data1:", data1)

	for data := range ch1 { //此时将所有数据均从channel中取出，会执行上述exit goroutine...结束阻塞
		fmt.Println("from channel read data:", data)
	}
	fmt.Println("all done...")
}

//示例3: 创建指针类型channel
func TestThreeCreateChannel(t *testing.T) {
	//make创建一个指针类型的channel
	ch := make(chan *TestChannel)
	go func() {
		fmt.Println("start goroutine...")
		channl := &TestChannel{
			Id:   1,
			Name: "yufeng",
		}
		ch <- channl
		fmt.Println("exit goroutine...")
	}()
	fmt.Println("wait goroutine...")
	data := <-ch
	fmt.Printf("from channel read message,[id=%v,name=%s]", data.Id, data.Name)
}

//示例4: 并发打印-模拟生产者，消费者模式
func TestGoroutinePrinter(t *testing.T) {
	ch := make(chan int)
	//并发执行printer
	go printer(ch)
	//生产数据
	for i := 1; i <= 10; i++ {
		//将数据通过channel传递给printer
		ch <- i
	}
	//通知并发的printer结束循环(已没有数据)
	ch <- 0
	//等待printer结束
	<-ch
}

func printer(ch chan int) {
	for {
		//消费数据
		data := <-ch
		if data == 0 {
			break
		}
		fmt.Println("from channel read data:", data)
	}
	ch <- 0
}

//示例5: 单向通道
func TestSingleChannel(t *testing.T) {
	//只能读取的channel，不能写入。此channel无意义(即：一个不能发送数据只能读取的channel是无意义的)
	ch1 := make(<-chan int)
	//只能写入数据的channel,不能读取
	ch2 := make(chan<- string)
	go func() {
		ch2 <- "yufeng"
		//Timer结构体即为一个只能接收的单向通道
		time.NewTimer(time.Second)
	}()

	<-ch1
}

/*为什么Go语言对通道要限制长度而不提供无限长度的通道？
  通道（channel）是在两个goroutine间通信的桥梁。使用goroutine的代码必然有一方提供数据，一方消费数据。
  当提供数据一方的数据供给速度大于消费方的数据处理速度时，如果通道不限制长度，那么内存将不断膨胀直到应用崩溃。
  因此，限制通道的长度有利于约束数据提供方的供给速度，供给数据量必须在消费方处理量+通道长度的范围内，才能正常地处理数据。
*/
//示例6: 带缓冲的通道
func TestBufferChannel(t *testing.T) {
	//创建一个带有3个元素缓冲大小的整型类型channel
	ch := make(chan int, 3)
	fmt.Println("向channel添加数据前,当前通道大小:", len(ch)) //0
	//发送3个int型元素到channel，因使用了缓冲channel,即使没有goroutine接收，发送者也不会发生阻塞
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println("向channel添加数据后,当前通道大小:", len(ch)) //3
	ch <- 4
	fmt.Println("超过缓冲大小再向channel添加数据后,当前通道大小:", len(ch)) //异常，此时缓冲channel已被填满3,再发送数据时发生阻塞
}

//示例7: 通道的多路复用-同时处理接收和发送多个channel的数据
//实现：模拟远程过程调用RPC(Remote Procedure Call,远程过程调用),利用channel代替socket实现RPC过程
func TestMultiplexingChannel(t *testing.T) {
	//创建一个无缓冲字符串channel
	ch := make(chan string)
	//并发执行服务器逻辑
	go RPCServer(ch)
	//客户端请求数据和接收数据
	response, err := RPCClient(ch, "hello")
	if err != nil {
		fmt.Println("err")
		return
	}
	fmt.Println("client received:", response)
}

//模拟RPC客户端的请求和接收消息封装
func RPCClient(ch chan string, req string) (string, error) {
	//向服务器发送请求
	ch <- req
	//select实现多路复用
	select {
	//等待服务器返回
	case ack := <-ch: //接收到服务器返回数据
		return ack, nil
	case <-time.After(time.Second): //超时
		return "", errors.New("Time out")
	}
}

//模拟RPC服务器端接收客户端请求和回应
func RPCServer(ch chan string) {
	for {
		//接收client端请求
		data := <-ch
		fmt.Println("RPC server receive data:", data)
		//向client端反馈已收到
		ch <- "RPC Server Have Receive Client Data:" + data
	}
}

//示例8: 使用channel响应计时器的事件
func TestTimerInChannel(t *testing.T) {
	//声明一个退出用的channel,向里写数据代表"退出"
	exit := make(chan int)
	fmt.Println("start...")
	//1秒后调用匿名函数,此匿名函数会在另外一个goroutine中被调用
	time.AfterFunc(time.Second, func() {
		fmt.Println("one second after...")
		//通知goroutine已经结束
		exit <- 0
	})
	<-exit
	fmt.Println("all done...")
}

/*
计时器Timer: 原理和闹钟类似，即给定多少时间后触发
打点器Ticker: 原理和钟表类似，每整点触发
*/
//示例9:定时计时Timer和打点器Ticker比较
func TestTimerPKTickerInMultiplexingChannel(t *testing.T) {
	//创建一个打点器，每500ms触发一次
	ticker := time.NewTicker(time.Millisecond * 500)
	//创建一个计时器，每2s触发一次
	stopper := time.NewTimer(time.Second * 2)
	//声明计数变量
	var i int
	//不断检查channel
	for {
		//多路复用channel
		select {
		case <-stopper.C: //计时器到了
			fmt.Println("stop...")
			goto StopHere
		case <-ticker.C: //打点器触发了
			i++ //记录触发了多少次
			fmt.Println("ticker:", i)
		}
	}
StopHere:
	fmt.Println("all done...")
}

//示例10:关闭channel后继续使用channel
//demo1: 关闭channel后再继续向channel中发送data会panic
func TestOneCloseChannel(t *testing.T) {
	ch := make(chan int)
	close(ch)                                                  //关闭channel
	fmt.Printf("ptr:%p cap:%d len:%d\n", ch, cap(ch), len(ch)) //ptr:0xc0000742a0 cap:0 len:0
	ch <- 1                                                    //再向channel中发送数据,此时panic
}

//demo2:从已关闭的channel接收data时将不会发生阻塞
func TestTwoCloseChannel(t *testing.T) {
	//创建一个缓冲channel,缓冲区大小为2
	ch := make(chan int, 2)
	//向channel中写入2条数据
	ch <- 1
	ch <- 2
	//关闭缓冲channel
	close(ch)
	//cap(ch)：获取一个对象的容量,即channel的缓冲大小,此处多遍历一个元素，目的：channel越界访问
	for i := 0; i < cap(ch)+1; i++ {
		//非阻塞接收数据,data:接收的数据，未接收到数据则data为0值;ok:是否接收到数据
		data, ok := <-ch
		//data := <- ch //阻塞模式接收数据,直到接收到数据并赋值给data
		//此处可正常打印channel中写入的2条数据，第三条数据为0 false,即从已关闭的channel依然能够访问数据,即使channel没有数据，在获取时也不会发生阻塞
		fmt.Printf("data:%v,ok:%v \n", data, ok)
	}
}
