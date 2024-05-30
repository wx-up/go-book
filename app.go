package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wx-up/go-book/pkg/saramax"
)

// App  wire.Struct(new(App),"*")  自动填充结构体中的字段，如果是在一个包的话字段可以小写
type App struct {
	engine *gin.Engine
	// 消费者本身也是类似一个 web 服务的东西是需要启动的，因此引入 app 这个结构体
	cs []saramax.Consumer
}
