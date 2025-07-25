package main

func main() {
	// 常量和枚举
	const a int = 1

	const (
		h    byte = 3
		i         = "value"
		j, k      = "v", 4
		l, m      = 5, false
	)

	type Gender string
	const (
		Male   Gender = "Male"
		Female Gender = "Female"
	)
}
