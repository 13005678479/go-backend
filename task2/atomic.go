// 使用原子操作（ sync/atomic 包）实现一个无锁的计数器。
// 启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 定义共享计数器和互斥锁
var counter1 int64

func doAtomic(wg *sync.WaitGroup) {
	// 用defer确保所有递增完成后再通知WaitGroup
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		// 原子递增操作：直接对counter的内存地址进行操作，保证操作的原子性
		atomic.AddInt64(&counter1, 1)
	}
}

func main() {
	// 必须使用int64类型，因为atomic包的递增函数针对int64
	var wg sync.WaitGroup
	// 启动10个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go doAtomic(&wg)
	}
	// 等待所有协程完成
	wg.Wait()
	fmt.Printf("最终计数器值: %d\n", counter1)
}
