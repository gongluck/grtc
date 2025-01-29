/*
 * @Author: gongluck
 * @Date: 2025-01-29 00:29:27
 * @Last Modified by: gongluck
 * @Last Modified time: 2025-01-29 13:58:33
 */

package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// 页面数据结构
type PageData struct {
	SignalServiceURL string // 用于存储信令服务的 URL
}

func main() {
	// 读取命令行参数
	var (
		signalServiceURL string
		httpPort         string
		httpsPort        string
	)
	flag.StringVar(&signalServiceURL, "signal-url", "http://localhost:8080", "Signal service URL")
	flag.StringVar(&httpPort, "http-port", "8081", "HTTP server port")
	flag.StringVar(&httpsPort, "https-port", "8082", "HTTPS server port") // 使用 8082 作为 HTTPS 端口
	flag.Parse()

	r := gin.Default()

	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 首页路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", PageData{
			SignalServiceURL: signalServiceURL, // 将信令服务的 URL 传递给模板
		})
	})

	// 动态路由，处理其他模板文件
	r.GET("/:templateName", func(c *gin.Context) {
		templateName := c.Param("templateName")
		c.HTML(http.StatusOK, templateName, PageData{
			SignalServiceURL: signalServiceURL, // 将信令服务的 URL 传递给模板
		})
	})

	var wg sync.WaitGroup

	// 启动HTTP服务器
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("HTTP server listening on :%s\n", httpPort)
		if err := r.Run(":" + httpPort); err != nil {
			log.Fatalf("Failed to run HTTP server: %v", err)
		}
	}()

	// 启动HTTPS服务器
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("HTTPS server listening on :%s\n", httpsPort)
		if err := r.RunTLS(":"+httpsPort, "./cert.pem", "./key.pem"); err != nil {
			log.Fatalf("Failed to run HTTPS server: %v", err)
		}
	}()

	// 等待所有 goroutine 完成
	wg.Wait()
}
