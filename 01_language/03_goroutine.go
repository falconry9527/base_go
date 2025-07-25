package main

import (
	"fmt"
	"time"
)

type UnsafeCounter struct {
	count int
}

// 增加计数
func (c *UnsafeCounter) Increment() {
	c.count += 1
}

// 获取当前计数
func (c *UnsafeCounter) GetCount() int {
	return c.count
}

func main() {
	counter := UnsafeCounter{}

	// 启动100个goroutine同时增加计数
	for i := 0; i < 1000; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				counter.Increment()
			}
		}()
	}
	// 等待一段时间确保所有goroutine完成
	time.Sleep(time.Second)
	// 输出最终计数
	fmt.Printf("Final count: %d\n", counter.GetCount())
}
