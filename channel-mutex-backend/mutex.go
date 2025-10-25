// 编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
// 启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
package main

import (
	"fmt"
	"sync"
)

// 定义共享计数器和互斥锁
var (
	counter int
	// 用于保护counter的互斥锁
	mu sync.Mutex
)

// 递增函数：使用互斥锁保护计数器的递增操作
func increment(wg *sync.WaitGroup) {
	// 用defer确保所有递增完成后再通知WaitGroup
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		// 加锁：确保同一时间只有一个协程能操作counter
		mu.Lock()
		// 临界区：对共享变量的修改
		counter++
		// 解锁：允许其他协程访问counter
		mu.Unlock()
	}
}

func mains() {
	var wg sync.WaitGroup
	// 启动10个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go increment(&wg)
	}
	// 等待所有协程完成
	wg.Wait()
	fmt.Printf("最终计数器值: %d\n", counter)
}
