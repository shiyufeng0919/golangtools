package 并发

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
)

//同步 - 保证并发环境下数据访问的正确性(原子操作，互斥锁，互斥读写锁，线程组)
/*
Go程序可以使用通道进行多个goroutine间的数据交换，但这仅仅是数据同步中的一种方法。通道内部的实现依然使用了各种锁，因此优雅代码的代价是性能。
在某些轻量级的场合，原子访问（atomic包）、互斥锁（sync.Mutex）以及等待组（sync.Wait Group）能最大程度满足需求。
*/

var (
	seq1 int64
	seq2 int64

	count      int        //逻辑中使用的某个变量
	countMutex sync.Mutex //与变量对应的使用互斥锁

	num        int
	numRWMutex sync.RWMutex //读写互斥锁，读比写多情况下更高效
)

//示例1:原子操作保证变量不发生竞态。竞态检测-检测代码在并发环境下可能出现的问题，可以使用互斥锁（sync.Mutex）解决竞态问题，但是对性能消耗较大。在这种情况下，推荐使用原子操作（atomic）进行变量操作
func TestAtomicOperation(t *testing.T) {
	//生成10个并发序列号
	for i := 0; i < 10; i++ {
		go genIdOne()
		go genIdTwo()
	}
	fmt.Println("genIdOne():", genIdOne())
	fmt.Println("genIdTwo():", genIdTwo())
}

//序列号生成器
func genIdOne() int64 {
	//尝试原子的增加序列号
	atomic.AddInt64(&seq1, 1)
	//此seq会产生竞态问题
	return seq1
}

//序列号生成器(推荐)
func genIdTwo() int64 {
	//尝试原子的增加序列号
	return atomic.AddInt64(&seq2, 1)
}

//示例2:互斥锁(sync.Mutex)，保证同时只有一个goroutine可以访问共享资源
func TestSyncMutext(t *testing.T) {
	for i := 0; i < 10; i++ {
		go setCount(i)
	}
	fmt.Println(getCount())
}

func getCount() int {
	countMutex.Lock()
	defer countMutex.Unlock()
	return count
}

func setCount(c int) {
	countMutex.Lock()
	count = c
	fmt.Println("set count value:", count)
	countMutex.Unlock()
}

//示例3:读写互斥锁(sync.RWMutex)-在读比写多的环境下，比互斥锁效率高
func TestSyncRWMutex(t *testing.T) {
	for i := 0; i < 10; i++ {
		go setNumber(i)
	}
}

func setNumber(c int) {
	numRWMutex.RLock()
	num = c
	fmt.Println("set num value:", num)
	numRWMutex.RUnlock()
}

//示例4:等待组(sync.WaitGroup) -保证在并发环境中完成指定数量的任务
//除了可以使用通道（channel）和互斥锁进行两个并发程序间的同步外，还可以使用等待组进行多个任务的同步
func TestSyncWaitGroup(t *testing.T) {
	var wgr sync.WaitGroup
	//此一组等待任务只需一个"等待组"
	var urls = []string{
		"http://www.baidu.com/",
		"http://www.github.com/",
		"https://www.golangtc.com/",
	}
	for _, url := range urls {
		//每一个任务开始时，将等待组加1
		wgr.Add(1)
		//开启一个并发(将url通过goroutine的参数进行传递，是为了避免url变量通过闭包放入匿名函数后又被修改的问题)
		go func(url string) {
			//表示函数执行完成时将等待组减1
			defer wgr.Done() //<=> wg.Add(-1)
			//使用http访问提供的网址
			_, err := http.Get(url)
			fmt.Printf("访问url:%s,发生错误:%v \n", url, err)
		}(url)
	}
	//等待所有任务完成,停止阻塞
	wgr.Wait()
	fmt.Println("all done...")
}
