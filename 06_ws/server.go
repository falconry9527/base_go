package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求，生产环境应该限制
	},
}

// Client 代表一个WebSocket客户端
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	log.Println("WebSocket 服务器启动在 :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级HTTP连接到WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket升级失败:", err)
		return
	}
	defer conn.Close()

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	log.Printf("客户端连接成功: %s", conn.RemoteAddr())

	// 启动读写goroutine
	go client.writePump()
	client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
		log.Printf("客户端断开连接: %s", c.conn.RemoteAddr())
	}()

	c.conn.SetReadLimit(512) // 限制消息大小
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读取错误: %v", err)
			}
			break
		}

		log.Printf("收到消息: %s", string(message))

		// 回复客户端
		response := []byte("服务器收到: " + string(message))
		c.send <- response
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(10 * time.Second) // 心跳包
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// channel关闭
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("写入错误: %v", err)
				return
			}

		case <-ticker.C:
			// 发送心跳包
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
