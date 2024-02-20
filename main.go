package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wx-up/go-book/internal/web"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	web.RegisterRoutes(engine)

	addr := ":8080"
	server := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	startErrChain := make(chan error, 2)
	go func() {
		startErrChain <- server.ListenAndServe()
	}()

	go func() {
		// 2秒后没有报错，则输出服务启动成功
		ticker := time.NewTicker(time.Second * 2)
		select {
		case err := <-startErrChain:
			ticker.Stop()
			log.Fatalf("Server Start Fail: %s\n", err)
		case <-ticker.C:
			ticker.Stop()
			log.Printf("Server Start Success，Listen On %s\n", addr)
			if err := <-startErrChain; err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server Start Fail: %s\n", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit

	log.Println("Trying Server Shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Fail: ", err)
	}

	log.Println("Server Shutdown Success")
}
