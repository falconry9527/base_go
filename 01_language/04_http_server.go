package main

import (
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	println("========")
	// 设置 Content-Type，根据需求可改为 text/plain、application/octet-stream 等
	w.Header().Set("Content-Type", "application/octet-stream")
	// 直接写入原始数据到响应体
	rawData := []byte("Hello, this is raw data!\n")
	w.Write(rawData)
}

func main() {
	http.HandleFunc("/raw", helloHandler)
	http.ListenAndServe(":8080", nil)
}
