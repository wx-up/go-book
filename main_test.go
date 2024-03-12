package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/redis/go-redis/v9"
)

type Phone string

func (p Phone) MarshalJSON() ([]byte, error) {
	if len(p) != 11 {
		return json.Marshal(string(p))
	}

	// 注意需要强制转化一下，否则还是 Phone 类型，调用 json.Marshal 会递归
	v := string(p[0:4] + "***" + p[7:])
	return json.Marshal(v)
}

type User struct {
	Phone Phone `json:"phone"`
}

func Ha() ([]byte, error) {
	return json.Marshal(map[string]any{
		"name": "age",
	})
}

func Test(t *testing.T) {
	var u User
	u.Phone = "13800138000"
	bs, _ := json.Marshal(u)
	fmt.Println(string(bs))

	fmt.Println(rand.Intn(10000))
	// nil指针寻址的话会 panic
	// var s *string
	// fmt.Println(s == nil)
	// fmt.Println(*s)
}

var luaTest = `
return true
`

func Test_Redis(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	val, err := client.Eval(context.Background(), luaTest, []string{}).Bool()
	fmt.Println(err == redis.Nil)
	fmt.Println(val)
}
