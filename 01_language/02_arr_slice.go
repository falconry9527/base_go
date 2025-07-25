package main

import "fmt"

func main() {
	// slice
	s := []int{1, 2, 3, 4, 5}
	// 方法1：使用索引
	for i := 0; i < len(s); i++ {
		fmt.Printf("索引: %d, 值: %d\n", i, s[i])
	}
	print("---------------")
	var b []int = []int{4, 5, 6, 7, 78, 8}
	for i := 0; i < len(b); i++ {
		fmt.Printf("索引: %d, 值: %d\n", i, b[i])
	}

}
