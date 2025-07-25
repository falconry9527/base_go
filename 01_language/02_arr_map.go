package main

import "fmt"

func main() {
	// 初始化方式1
	var m1 map[string]int
	m1 = map[string]int{}
	m1["222"] = int(1)

	// 初始化方式2
	m2 := map[string]int{}
	m2["222"] = int(2)

	// 初始化方式3
	m3 := make(map[string]int)
	m3["222"] = 2

	m := make(map[string]int, 10)
	m["1"] = int(1)
	m["2"] = int(2)
	m["3"] = int(3)
	m["4"] = int(4)
	m["5"] = int(5)
	m["6"] = int(6)

	// 元素操作
	// 获取元素
	value1 := m["1"]
	fmt.Println("m[\"1\"] =", value1)

	// 判断元素是否存在
	value1, exist := m["1"]
	fmt.Println("m[\"1\"] =", value1, ", exist =", exist)

	// 修改值
	fmt.Println("before modify, m[\"2\"] =", m["2"])
	m["2"] = 20
	fmt.Println("after modify, m[\"2\"] =", m["2"])

	// 删除值
	delete(m, "1")
	value11, exist1 := m["1"]
	fmt.Println("m[\"1\"] =", value11, ", exist =", exist1)

	// 获取map的长度
	m["10"] = 10
	fmt.Println("after add, len(m) =", len(m))

	// 遍历map集合main
	for key, value := range m {
		fmt.Println("iterate map, m[", key, "] =", value)
	}

}
