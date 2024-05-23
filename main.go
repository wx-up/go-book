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

	"github.com/wx-up/go-book/pkg/otelx"

	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"go.opentelemetry.io/otel"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.uber.org/zap"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/pflag"

	"github.com/spf13/viper"

	_ "github.com/spf13/viper/remote"
)

func main() {
	// InitConfigByRemote()
	initConfig()
	initLogger()
	initPrometheus()
	initOTLP()

	app := InitWebService()
	app.engine.ContextWithFallback = true

	// 启动消费者
	// 这里不够优雅，如果第一个消费者启动成功，后面的消费者失败，理论上应该让第一个消费则退出，go-zero有类似的实现
	// 这里粗暴一点，直接panic退出
	for _, c := range app.cs {
		if err := c.Start(); err != nil {
			panic(err)
		}
	}

	// 启动服务
	addr := ":8080"
	server := &http.Server{
		Addr:    addr,
		Handler: app.engine,
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

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
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

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}

func initOTLP() func(ctx context.Context) {
	res, err := newResource("demo", "v0.0.1")
	if err != nil {
		panic(err)
	}

	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// 初始化 trace provider
	// 这个 provider 就是用来在打点的时候构建 trace 的
	tp, err := newTraceProvider(res)
	if err != nil {
		panic(err)
	}
	// 用完需要关闭
	newTp := otelx.NewMyTraceProvider(tp)
	otel.SetTracerProvider(newTp)
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 检测配置文件的变化
		newTp.Enabled.Store(viper.GetBool("otel.enabled"))
	})
	return func(ctx context.Context) {
		tp.Shutdown(ctx)
	}
}

// resource 代表系统或者模块的抽象
func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		// SchemaURL 主要用于指定什么版本，照着抄就行不用管
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
}

// newPropagator 传播器，用于跨服务传播spanContext
func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
	// propagation.TraceContext{}, // 用于传递分布式追踪的上下文信息，包括 Trace ID 和 Span ID
	// propagation.Baggage{}, // 用于传递用户自定义的 baggage 数据
	)
}

// newTraceProvider 用于指定 trace 真正的实现者
func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	// 这里使用 zipkin 的 exporter
	// 还可以使用：jeager、skywalking
	// 可以简单认为：zipkin 适配了 otel 的 API 所以代码层面我们使用的是 otel 的 API 实际上我们数据是上传到 zipkin
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans")
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}
