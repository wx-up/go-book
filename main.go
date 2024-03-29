package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/pflag"

	"github.com/spf13/viper"

	_ "github.com/spf13/viper/remote"
)

func main() {
	// InitConfigByRemote()
	initConfig()

	engine := InitWebService()

	// 启动服务
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

// initConfig 读取配置文件
func initConfig() {
	cFile := pflag.String("config", "config/config.yaml", "指定配置文件")
	pflag.Parse() // 解析命令行参数
	viper.SetConfigFile(*cFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	fmt.Println(viper.AllKeys())
	type C struct {
		DSN string
	}
	var c C
	fmt.Println(viper.UnmarshalKey("db.mysql", &c))
	fmt.Println(viper.Get("db.mysql"))
	fmt.Println(c)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println(viper.UnmarshalKey("db.mysql", &c))
		fmt.Println(c)
	})
}

// InitConfigByRemote etcd 远程配置中心读取
func InitConfigByRemote() {
	// etcd 代表 etcd2.x 以及之前的版本
	// etcd3 代表 etcd3.x 以及之后的版本，一般都是用 etcd3
	// 使用 path 来做隔离，类似命名空间的概念
	err := viper.AddRemoteProvider("etcd3", "http://127.0.0.1:12379", "/go_book")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
	type C struct {
		Dsn string
	}
	var c C
	viper.UnmarshalKey("db.mysql", &c)
	fmt.Println(c)
	go func() {
		for {
			time.Sleep(time.Second * 5)
			viper.WatchRemoteConfig()
			viper.UnmarshalKey("db.mysql", &c)
			fmt.Println(c)
		}
	}()
}
