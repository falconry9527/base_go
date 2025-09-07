package client

import (
	"testing"
	"time"
)

func TestWebSocketClient(t *testing.T) {
	client := NewWebSocketClient("ws://localhost:8080/ws")

	// 测试连接
	if err := client.Connect(); err != nil {
		t.Skip("测试服务器未启动，跳过测试")
	}
	defer client.Disconnect()

	// 测试连接状态
	if !client.IsConnected() {
		t.Error("客户端应该处于连接状态")
	}

	// 测试发送消息
	if err := client.SendText("test message"); err != nil {
		t.Errorf("发送消息失败: %v", err)
	}

	// 测试发送JSON
	testData := map[string]interface{}{
		"action": "test",
		"data":   "hello",
	}
	if err := client.SendJSON(testData); err != nil {
		t.Errorf("发送JSON失败: %v", err)
	}

	// 等待一下让消息发送完成
	time.Sleep(100 * time.Millisecond)

	// 测试获取连接信息
	info := client.GetConnectionInfo()
	if info["connected"] != true {
		t.Error("连接信息显示未连接")
	}
}
