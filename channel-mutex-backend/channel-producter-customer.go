// 实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
package main

import (
	"fmt"
	"sync"
)

// 生产者：向缓冲通道发送100个整数
// chan<- 是 Go 语言中对通道的 “方向限定”，表示这个通道参数在函数内部只能用于发送数据
// wg *sync.WaitGroup 是一个指向 sync.WaitGroup 类型的指针参数。
func producer(ch chan<- int, wg *sync.WaitGroup) {
	// 通知WaitGroup当前协程完成
	defer wg.Done()
	for i := 1; i <= 100; i++ {
		// 发送数据到缓冲通道
		ch <- i
		// fmt.Printf("发送: %d\n", i)
	}
	// 发送完毕后关闭通道
	close(ch)
}

// 消费者：从通道接收数据并打印
func consumer(ch <-chan int, wg *sync.WaitGroup) {
	// 通知WaitGroup当前协程完成
	defer wg.Done()
	for num := range ch {
		fmt.Printf("接收: %d\n", num)
	}
}

func mainT() {
	// 创建带缓冲的通道，缓冲区大小设置为10（可根据需要调整）
	bufferedChan := make(chan int, 10)
	var wg sync.WaitGroup

	// 启动生产者协程
	wg.Add(1)
	go producer(bufferedChan, &wg)

	// 启动消费者协程
	wg.Add(1)
	go consumer(bufferedChan, &wg)

	// 等待生产者和消费者完成
	wg.Wait()
	fmt.Println("所有数据处理完毕")
}
