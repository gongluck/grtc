/*
 * @Author: gongluck
 * @Date: 2025-01-29 20:06:25
 * @Last Modified by: gongluck
 * @Last Modified time: 2025-01-29 22:43:24
 */

package main

import (
	"flag"
	"log"
	"pusher/transport"
	"pusher/util"
	"pusher/webrtc"
	"time"
)

func main() {
	// 读取命令行参数
	var (
		signalServiceURL string
	)
	flag.StringVar(&signalServiceURL, "signal-url", "ws://localhost:8080", "Signal service URL")
	flag.Parse()

	// 获取唯一设备 ID（MAC 地址）
	deviceID := util.GetDeviceID()

	// 定义消息处理函数
	messageHandler := func(message []byte) {
		log.Printf("Received message: %s\n", message)
		webrtc.HandleSignaling(message)
	}

	// 启动 WebSocket 客户端
	client := transport.NewWebSocketClient(signalServiceURL+"/"+deviceID+"/pusher/"+deviceID, deviceID, messageHandler)

	// 尝试连接并保持连接状态
	for {
		// 连接到 WebSocket 服务器
		if err := client.Connect(); err != nil {
			log.Printf("Could not connect to WebSocket server: %v. Retrying in 5 seconds...\n", err)
			time.Sleep(5 * time.Second) // 等待 5 秒后重试
			continue
		}

		// 连接成功后保持程序运行
		log.Println("Connected to WebSocket server. Listening for messages...")
		client.Listen() // 监听消息并保持连接
	}
}
