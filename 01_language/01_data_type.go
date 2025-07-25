package main

import "fmt"

func main() {
	// 一.基础数据类型: 整型，浮点数，string，byte
	// 整型：int，int8，int16，int32，int64，uint，uint8，uint16，uint32，uint64，uintptr。
	// 十进制
	var h uint8 = 15
	print(h)

	// 浮点数：float32，float64。
	var float1 float32 = 10
	float2 := 10.0
	print(float1, float2)
	// string
	var s string = "Hello, world!"
	print(s)

	// byte 类型
	var bytes []byte = []byte(s)
	fmt.Println("convert \"Hello, world!\" to bytes: ", bytes)

}
