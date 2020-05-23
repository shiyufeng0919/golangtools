package lock_test

/*
  示例：golang 互斥锁 sync.Mutex add by syf 2020.5.13
  sync.Mutex为互斥锁（也叫全局锁），Lock()加锁，Unlock()解锁
  适用于场景：读写次数不确定的场景（读写次数没有明显区别），同时只能一个读或者写
*/
