package main

import "fmt"

func main() {
	var i int32 = 17
	var b byte = 5
	var f float32

	// 1. 数字类型可以直接强转
	f = float32(i) / float32(b)
	fmt.Printf("f 的值为: %f\n", f)
	// 当int32类型强转成byte时，高位被直接舍弃
	var i2 int32 = 256
	var b2 byte = byte(i2)
	fmt.Printf("b2 的值为: %d\n", b2)

	// 2. string 和  []byte 的转换
	str := "hello, 123, 你好"
	var bytes []byte = []byte(str)
	fmt.Printf("bytes 的值为: %v \n", bytes)
	fmt.Printf("bytes 的值为: %v \n", string(bytes))

}
