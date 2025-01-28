/*
 * @Author: gongluck
 * @Date: 2025-01-28 18:59:59
 * @Last Modified by: gongluck
 * @Last Modified time: 2025-01-29 00:21:37
 */

package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	heartbeatInterval = 30 * time.Second // 心跳间隔
	writeWait         = 10 * time.Second // 写超时
	pongWait          = 60 * time.Second // pong等待时间
)

// WebSocket 升级器配置
var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true }, // 允许所有来源
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// 房间结构
type Room struct {
	sync.Mutex
	PusherConn *websocket.Conn // 推送者连接
	PusherID   string          // 推送者ID
	PullerConn *websocket.Conn // 拉取者连接
	PullerID   string          // 拉取者ID
}

// 连接结构
type connection struct {
	ws         *websocket.Conn // WebSocket连接
	lastActive time.Time       // 最后活动时间
	sync.Mutex
}

// 使用sync.Map存储房间
var rooms sync.Map

func main() {
	// 初始化日志记录器
	initlog()

	// 启动HTTP服务
	go runHTTPServer()

	// 启动HTTPS服务
	go runHTTPSServer()

	// 保持主程序运行
	select {}
}

// 初始化日志记录器，支持文件和控制台输出
func initlog() {
	// 创建日志文件
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// 使用 io.MultiWriter 同时输出到文件和控制台
	mw := io.MultiWriter(file, os.Stdout)
	gin.DefaultWriter = mw
	logrus.SetOutput(mw)
	logrus.SetFormatter(&logrus.TextFormatter{ // 设置日志格式
		DisableColors: false,
		FullTimestamp: true,
		ForceColors:   true,
	})
}

func runHTTPServer() {
	gin.SetMode(gin.ReleaseMode) // 设置Gin为发布模式
	r := gin.Default()
	r.Use(gin.Recovery()) // 使用Gin的恢复中间件

	// 注册路由
	registerRoutes(r)

	// HTTP服务器配置
	srv := &http.Server{
		Addr:    ":8080", // 监听8080端口
		Handler: r,
	}

	logrus.Info("HTTP server listening on :8080")
	if err := srv.ListenAndServe(); err != nil {
		logrus.Fatalf("HTTP server failed: %v", err)
	}
}

func runHTTPSServer() {
	gin.SetMode(gin.ReleaseMode) // 设置Gin为发布模式
	r := gin.Default()
	r.Use(gin.Recovery()) // 使用Gin的恢复中间件

	// 注册路由
	registerRoutes(r)

	// HTTPS服务器配置
	srv := &http.Server{
		Addr:    ":8443", // 监听8443端口
		Handler: r,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12, // 设置最小TLS版本
		},
	}

	logrus.Info("HTTPS server listening on :8443")
	if err := srv.ListenAndServeTLS("./cert.pem", "./key.pem"); err != nil {
		logrus.Fatalf("HTTPS server failed: %v", err)
	}
}

// 注册路由
func registerRoutes(r *gin.Engine) {
	r.GET("/:roomid/pusher/:id", handlePusher) // 处理推送者连接
	r.GET("/:roomid/puller/:id", handlePuller) // 处理拉取者连接
	r.GET("/status", handleStatus)             // 处理状态请求
}

// 处理推送者连接
func handlePusher(c *gin.Context) {
	handleConnection(c, true) // 是推送者
}

// 处理拉取者连接
func handlePuller(c *gin.Context) {
	handleConnection(c, false) // 是拉取者
}

// 处理连接
func handleConnection(c *gin.Context, isPusher bool) {
	roomID := c.Param("roomid") // 获取房间ID
	clientID := c.Param("id")   // 获取客户端ID

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil) // 升级连接
	if err != nil {
		logrus.Errorf("Failed to upgrade connection: %v", err)
		return
	}

	// 创建带心跳的连接
	wsConn := &connection{
		ws:         conn,
		lastActive: time.Now(),
	}

	// 设置Pong处理
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		wsConn.Lock()
		wsConn.lastActive = time.Now() // 更新最后活动时间
		wsConn.Unlock()
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// 获取或创建房间
	room, _ := rooms.LoadOrStore(roomID, &Room{})
	r := room.(*Room)
	r.Lock()

	// 处理连接类型
	if isPusher {
		if r.PusherConn != nil {
			r.PusherConn.Close() // 关闭现有推送者连接
		}
		r.PusherConn = conn
		r.PusherID = clientID
		logrus.Infof("Pusher connected: %s in room: %s", clientID, roomID)
	} else {
		if r.PullerConn != nil {
			r.PullerConn.Close() // 关闭现有拉取者连接
		}
		r.PullerConn = conn
		r.PullerID = clientID
		logrus.Infof("Puller connected: %s in room: %s", clientID, roomID)
	}

	r.Unlock()

	// 启动心跳goroutine
	go heartbeatTicker(wsConn)

	// 消息处理循环
	go func() {
		defer func() {
			conn.Close() // 关闭连接
			r.Lock()
			if isPusher && r.PusherConn == conn {
				r.PusherConn = nil
				r.PusherID = ""
			} else if !isPusher && r.PullerConn == conn {
				r.PullerConn = nil
				r.PullerID = ""
			}
			r.Unlock()
			logrus.Infof("Client disconnected: %s in room: %s", clientID, roomID)
		}()

		for {
			messageType, msg, err := conn.ReadMessage() // 读取消息
			if err != nil || messageType != websocket.TextMessage {
				break
			}

			r.Lock()
			// 根据连接类型转发消息
			if isPusher && r.PullerConn != nil {
				r.PullerConn.WriteMessage(websocket.TextMessage, msg)
			} else if !isPusher && r.PusherConn != nil {
				r.PusherConn.WriteMessage(websocket.TextMessage, msg)
			}
			r.Unlock()
		}
	}()
}

// 心跳定时器
func heartbeatTicker(conn *connection) {
	ticker := time.NewTicker(heartbeatInterval) // 创建心跳定时器
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			conn.Lock()
			// 检查最后活动时间
			if time.Since(conn.lastActive) > pongWait {
				conn.ws.Close() // 关闭连接
				conn.Unlock()
				return
			}

			// 发送Ping
			conn.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				conn.Unlock()
				return
			}
			conn.Unlock()
		}
	}
}

// 处理状态请求
func handleStatus(c *gin.Context) {
	result := make(map[string]interface{})

	// 遍历房间并收集状态信息
	rooms.Range(func(key, value interface{}) bool {
		roomID := key.(string)
		r := value.(*Room)
		r.Lock()
		defer r.Unlock()

		roomInfo := map[string]interface{}{
			"pusher_id": r.PusherID,
			"puller_id": r.PullerID,
			"active":    r.PusherConn != nil && r.PullerConn != nil,
		}
		result[roomID] = roomInfo
		return true
	})

	c.JSON(http.StatusOK, result) // 返回JSON状态信息
}
