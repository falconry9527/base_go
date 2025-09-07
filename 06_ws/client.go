package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// 连接WebSocket服务器
	url := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer conn.Close()

	fmt.Println("成功连接到WebSocket服务器!")
	fmt.Println("输入消息发送给服务器，输入 'exit' 退出")

	// 启动goroutine接收消息
	done := make(chan struct{})
	go receiveMessages(conn, done)

	// 读取用户输入并发送
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()

		if text == "exit" {
			fmt.Println("正在退出...")
			break
		}

		// 发送消息到服务器
		err := conn.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			log.Println("发送失败:", err)
			break
		}

		fmt.Printf("已发送: %s\n", text)
	}

	// 优雅关闭
	close(done)
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(100 * time.Millisecond)
}

func receiveMessages(conn *websocket.Conn, done chan struct{}) {
	for {
		select {
		case <-done:
			return
		default:
			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("接收错误: %v", err)
				}
				return
			}

			fmt.Printf("收到服务器回复: %s\n", string(message))
		}
	}
}
