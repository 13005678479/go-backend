// 编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，
// 另一个协程从通道中接收这些整数并打印出来

package main

import (
	"fmt"
	"sync"
)

// sendNumbers 向通道发送1到10的整数
func sendNumbers(ch chan<- int, wg *sync.WaitGroup) {
	// 协程完成后通知WaitGroup
	defer wg.Done()
	for i := 1; i < 10; i++ {
		// 发送数据到通道
		ch <- i
	}
	// 发送完毕关闭通道
	close(ch)
}

// receiveAndPrint 从通道接收整数并打印
func receiveAndPrint(ch <-chan int, wg *sync.WaitGroup) {
	// 协程完成后通知WaitGroup
	defer wg.Done()
	for num := range ch {
		fmt.Printf("收到：%d\n", num)
	}
}

func mainf() {
	// 创建通道和WaitGroup
	numChan := make(chan int)
	var wg sync.WaitGroup

	// 启动发送协程
	wg.Add(1)
	// wg 是 sync.WaitGroup 类型的变量（用于等待多个协程完成）。
	// 当我们调用 sendNumbers 和 receiveAndPrint 函数时，参数传递的是 &wg，即 wg 的指针
	go sendNumbers(numChan, &wg)

	// 启动接收协程
	wg.Add(1)
	go receiveAndPrint(numChan, &wg)

	// 等待所有协程完成
	wg.Wait()
	fmt.Println("所有数据处理完毕，程序结束")
}
