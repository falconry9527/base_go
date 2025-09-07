package client

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketClient 封装的WebSocket客户端
type WebSocketClient struct {
	url                  string
	conn                 *websocket.Conn
	done                 chan struct{}
	sendChan             chan []byte
	mu                   sync.RWMutex
	isConnected          bool
	reconnectAttempts    int
	maxReconnectAttempts int
}

// Message 消息结构体
type Message struct {
	Type      string      `json:"type"`
	Content   interface{} `json:"content"`
	Timestamp int64       `json:"timestamp"`
}

// NewWebSocketClient 创建新的WebSocket客户端实例
func NewWebSocketClient(url string) *WebSocketClient {
	return &WebSocketClient{
		url:                  url,
		done:                 make(chan struct{}),
		sendChan:             make(chan []byte, 100),
		isConnected:          false,
		maxReconnectAttempts: 5,
	}
}

// ==================== 连接管理方法 ====================

// Connect 连接到WebSocket服务器
func (c *WebSocketClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected {
		return fmt.Errorf("已经连接到服务器")
	}

	fmt.Printf("正在连接到 %s...\n", c.url)

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second

	conn, _, err := dialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}

	c.conn = conn
	c.isConnected = true
	c.reconnectAttempts = 0

	fmt.Println("✅ 连接成功!")

	// 启动读写goroutine
	go c.readPump()
	go c.writePump()

	return nil
}

// Disconnect 断开连接
func (c *WebSocketClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		return fmt.Errorf("未连接到服务器")
	}

	close(c.done)

	// 发送关闭消息
	err := c.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	if err != nil {
		log.Printf("发送关闭消息失败: %v", err)
	}

	err = c.conn.Close()
	c.isConnected = false

	fmt.Println("连接已关闭")
	return err
}

// Reconnect 重新连接
func (c *WebSocketClient) Reconnect() error {
	if c.isConnected {
		c.Disconnect()
	}

	for attempt := 1; attempt <= c.maxReconnectAttempts; attempt++ {
		fmt.Printf("尝试第 %d 次重连...\n", attempt)

		err := c.Connect()
		if err == nil {
			fmt.Println("✅ 重连成功!")
			return nil
		}

		log.Printf("重连失败 (%d/%d): %v", attempt, c.maxReconnectAttempts, err)
		time.Sleep(time.Duration(attempt*2) * time.Second) // 指数退避
	}

	return fmt.Errorf("重连失败，已达到最大重试次数")
}

// IsConnected 检查是否已连接
func (c *WebSocketClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// ==================== 消息发送方法 ====================

// SendText 发送文本消息
func (c *WebSocketClient) SendText(message string) error {
	if !c.IsConnected() {
		return fmt.Errorf("未连接到服务器，无法发送消息")
	}
	c.sendChan <- []byte(message)
	return nil
}

// SendJSON 发送JSON消息
func (c *WebSocketClient) SendJSON(message interface{}) error {
	if !c.IsConnected() {
		return fmt.Errorf("未连接到服务器，无法发送消息")
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("JSON编码失败: %v", err)
	}

	c.sendChan <- jsonData
	return nil
}

// SendMessage 发送结构化消息
func (c *WebSocketClient) SendMessage(msgType string, content interface{}) error {
	message := Message{
		Type:      msgType,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}
	return c.SendJSON(message)
}

// ==================== 消息接收方法 ====================

// SetMessageHandler 设置消息处理回调
func (c *WebSocketClient) SetMessageHandler(handler func([]byte)) {
	// 这里可以实现消息处理回调机制
	// 简化版本，在readPump中直接处理
}

// ==================== 内部方法 ====================

func (c *WebSocketClient) readPump() {
	defer func() {
		c.mu.Lock()
		c.isConnected = false
		c.mu.Unlock()

		if r := recover(); r != nil {
			log.Printf("readPump发生panic: %v", r)
		}

		// 尝试重连
		if !c.isConnected {
			go c.Reconnect()
		}
	}()

	for {
		select {
		case <-c.done:
			return
		default:
			c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("读取错误: %v", err)
				}
				return
			}

			c.handleIncomingMessage(message)
		}
	}
}

func (c *WebSocketClient) writePump() {
	ticker := time.NewTicker(30 * time.Second) // 心跳间隔
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			return

		case message, ok := <-c.sendChan:
			if !ok {
				// channel关闭
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("消息发送失败: %v", err)
				// 发送失败，尝试重新放入channel（可选）
				go func() {
					time.Sleep(100 * time.Millisecond)
					select {
					case c.sendChan <- message:
					default:
						log.Println("发送队列已满，消息丢弃")
					}
				}()
			}

		case <-ticker.C:
			// 发送心跳包
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("心跳发送失败: %v", err)
				return
			}
		}
	}
}

func (c *WebSocketClient) handleIncomingMessage(message []byte) {
	fmt.Printf("📨 收到消息: %s\n", string(message))

	// 尝试解析为JSON
	var msg Message
	if err := json.Unmarshal(message, &msg); err == nil {
		fmt.Printf("📦 解析消息: 类型=%s, 内容=%v, 时间=%s\n",
			msg.Type,
			msg.Content,
			time.Unix(msg.Timestamp, 0).Format("15:04:05"))
	}
}

// ==================== 工具方法 ====================

// GetConnectionInfo 获取连接信息
func (c *WebSocketClient) GetConnectionInfo() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	info := make(map[string]interface{})
	info["connected"] = c.isConnected
	info["url"] = c.url
	if c.conn != nil {
		info["local_addr"] = c.conn.LocalAddr().String()
		info["remote_addr"] = c.conn.RemoteAddr().String()
	}
	return info
}

// SetMaxReconnectAttempts 设置最大重连次数
func (c *WebSocketClient) SetMaxReconnectAttempts(max int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.maxReconnectAttempts = max
}
