/*
 * @Author: gongluck
 * @Date: 2025-01-29 20:07:49
 * @Last Modified by: gongluck
 * @Last Modified time: 2025-01-29 20:20:34
 */

package transport

import (
	"github.com/gorilla/websocket"
	"log"
)

type WebSocketClient struct {
	address  string
	deviceID string
	conn     *websocket.Conn
}

// 创建新的 WebSocket 客户端
func NewWebSocketClient(address, deviceID string) *WebSocketClient {
	return &WebSocketClient{
		address:  address,
		deviceID: deviceID,
	}
}

// 连接到 WebSocket 服务器
func (client *WebSocketClient) Connect() error {
	var err error
	client.conn, _, err = websocket.DefaultDialer.Dial(client.address, nil)
	if err != nil {
		return err
	}
	log.Println("Connected to WebSocket server:", client.address)

	return nil
}

// 监听 WebSocket 消息
func (client *WebSocketClient) Listen() {
	defer client.conn.Close()
	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("Connection closed, attempting to reconnect:", err)
			return // 连接关闭，返回以便重连
		}
		log.Printf("Received message: %s\n", msg)
	}
}

// 发送消息
func (client *WebSocketClient) Send(message string) error {
	return client.conn.WriteMessage(websocket.TextMessage, []byte(message))
}
