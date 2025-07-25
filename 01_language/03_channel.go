package main

import (
	"fmt"
	"time"
)

// 只接收channel的函数
func receiveOnly(ch <-chan int) {
	for v := range ch {
		fmt.Printf("接收到: %d\n", v)
	}
}

// 只发送channel的函数
func sendOnly(ch chan<- int) {
	for i := 0; i < 5; i++ {
		ch <- i
		fmt.Printf("发送: %d\n", i)
	}
	close(ch)
}

func main() {
	// 基本操作
	// 3是channel的缓冲区,不是长度
	ch2 := make(chan int, 3)
	// 发送数据
	ch2 <- 2333
	// 接收数据
	message, ok := <-ch2
	fmt.Println(message, ok)
	// 关闭channel
	close(ch2)

	// 案例
	// 创建一个带缓冲的channel
	ch := make(chan int, 3)
	// 启动发送 goroutine
	go sendOnly(ch)
	// 启动接收 goroutine
	go receiveOnly(ch)

	// 使用select进行多路复用
	// timeout := time.After(2 * time.Second)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				fmt.Println("Channel已关闭")
				return
			}
			fmt.Printf("主goroutine接收到: %d\n", v)
		//case <-timeout:
		//	fmt.Println("操作超时")
		//	return
		default:
			fmt.Println("没有数据，等待中...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
