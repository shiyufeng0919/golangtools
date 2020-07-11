# GO语言并发

## #轻量级线程goroutine -根据需要随时创建的"线程"

有一种机制：使用者分配足够多的任务，系统能自动帮助使用者把任务分配到CPU上，让这些任务尽量并发运作。这种机制在Go语言中被称为goroutine。

goroutine的概念类似于线程，但goroutine由Go程序运行时的调度和管理。Go程序会智能地将goroutine中的任务合理地分配给每个CPU。

Go程序从main包的main()函数开始，在程序启动时，Go程序就会为main()函数创建一个默认的goroutine。