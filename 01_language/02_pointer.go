package main

import "fmt"

func main() {

	// 指针 pointer
	var p1 *int
	var p2 *string

	i := 1
	s1 := "Hello"
	p1 = &i
	p2 = &s1
	p3 := &p2
	fmt.Println(p1)
	fmt.Println(p2)
	fmt.Println(p3)

	// 访问指针
	fmt.Println(*p1)
	fmt.Println(*p2)
	fmt.Println(**p3)

	fmt.Println(*p1 == i)
}
