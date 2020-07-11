package 并发

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

/*
 1,所有goroutine在main()函数结束时会一同结束
*/

//示例1:使用普通函数创建goroutine,模拟两个任务互不影响
func TestCommonFunction(t *testing.T) {
	go runningA()
	go runningB()
}
func runningA() {
	var times int
	for {
		times++
		fmt.Println("yufeng", times)
		time.Sleep(time.Second)
	}
}
func runningB() {
	var input string
	fmt.Scanln(&input)
}

//示例2:使用匿名函数创建goroutine (注：go关键字后也可以为匿名函数或闭包启动goroutine)
func TestAnonymousFunction(t *testing.T) {
	//设置使用CPU核数，未设置：默认使用所有CPU核；设置1,2...来修改CPU核数。此参数用于设置程序使用的最多CPU核数
	//另:GOMAXPROCS同时也是一个环境变量，在应用程序启动前设置环境变量也可以起到相同的作用。GO在GOMAXPROCS数量与任务数量相等时，可以做到并行执行，但一般情况下都是并发执行。
	runtime.GOMAXPROCS(runtime.NumCPU())

	//匿名函数由一个不带函数名的函数声明和函数体组成；Go的匿名函数是一个闭包
	//demo1:不带参数的匿名函数
	go func() {
		fmt.Println("具体执行业务逻辑")
	}()

	//demo2:带有参数的匿名函数
	i := 1
	go func(i int) {
		fmt.Println("执行具体业务逻辑中:", i) //1
	}(i)
	i++
	fmt.Println("执行具体业务逻辑后:", i) //2

	//demo3:带有参数的匿名函数
	channel := make(chan string)
	go func(channel chan<- string) {
		channel <- "yufeng"
	}(channel)

	message := <-channel
	fmt.Println("message:", message) //yufeng
}
