package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
)

func main() {
	initConfig()
	initLogger()
	//initPrometheus()
	//initOTLP()
	app := InitApp()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			log.Println(err)
		}
	}
	err := app.server.Serve()
	log.Println(err)
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
	cFile := pflag.String("config", "config.yaml", "指定配置文件")
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
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println(viper.UnmarshalKey("db.mysql", &c))
		fmt.Println(c)
	})
}
