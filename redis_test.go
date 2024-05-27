package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func Test_Expire_Event(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 如果没有设置密码，则为空字符串
		DB:       0,  // 使用默认数据库
	})

	defer client.Close()

	// 修改配置,开启事件监听 ps: 修改配置文件,效果等同
	_, err := client.ConfigSet(context.Background(), "notify-keyspace-events", "Ex").Result()
	if err != nil {
		panic(err)
	}
	// 设置key，并设置过期时间为10秒
	err = client.Set(context.Background(), "mykey", "myvalue", 10*time.Second).Err()
	if err != nil {
		panic(err)
	}
	//订阅
	pubsub := client.Subscribe(context.Background(), "__keyevent@0__:expired")
	defer pubsub.Close()

	// 开启goroutine，接收过期事件
	go func() {
		for msg := range pubsub.Channel() {
			// 处理过期事件
			fmt.Println("Key expired:", msg.Payload)
		}
	}()
	select {} //阻塞主进程
}
