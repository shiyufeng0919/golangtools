package 并发

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

//利用两个goroutine并发运行，模拟生产者和消费者模式
func TestProducerAndConsumer(t *testing.T) {
	//创建一个字符串类型的通道
	channel := make(chan string)
	//创建生产者函数的并发goroutine，模拟3个并发，向通道写入消息
	go producer("kaixin", channel)
	go producer("yufeng", channel)
	go producer("fighting", channel)
	//数据消费函数,模拟从通道读取消息
	customer(channel)
}

//数据生产者
func producer(header string, channel chan<- string) {
	//无限循环，不停地生产数据
	for {
		//将随机数和字符串格式化为字符串发送给channel
		channel <- fmt.Sprintf("生产者向通道写消息:%s: %v", header, rand.Int31())
		//等待2秒
		time.Sleep(2 * time.Second)
	}
}

//数据消费者
func customer(channel <-chan string) {
	//不停地获取数据
	for {
		//从通道中取出数据，此处会阻塞直到通道中返回数据
		message := <-channel
		fmt.Println("消费者从通道读消息:", message)
	}
}
