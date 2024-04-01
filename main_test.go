package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"go.uber.org/zap"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/require"

	"github.com/redis/go-redis/v9"
)

func Test_Fail(t *testing.T) {
	a := 1
	require.Equal(t, a, 2)
	fmt.Println(666)
}

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
		Addr:     "localhost:7379",
		Password: "",
		DB:       0,
	})
	// val, err := client.Eval(context.Background(), luaTest, []string{}).Bool()
	// fmt.Println(err == redis.Nil)
	// fmt.Println(val)
	fmt.Println(client.Exists(context.Background(), "test-test").Result())
}

func Test_Viper(t *testing.T) {
	// 读取的文件名叫做 dev ，不包括文件扩展名，比如 .go、.yaml 等
	// 扩展名由 SetConfigType 指定
	viper.SetConfigName("dev")
	// 读取的文件类型是 yaml
	viper.SetConfigType("yaml")
	// 在当前目录的 config 目录下查找
	// 从函数的命名中可以知道：使用 add 而不是 set 是指可以添加多个路径
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	viper.SetConfigFile("config/dev.yaml")
	if err := viper.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	viper.GetInt64("db.port")
}

func Test_Viper2(t *testing.T) {
	cfg := `
db:
  dsn: "test"

redis:
  addr: "localhost:6379"
  password: ""
  db: 0
`
	viper.SetConfigType("yaml")

	_ = viper.ReadConfig(bytes.NewReader([]byte(cfg)))

	fmt.Println(viper.GetString("redis.addr"))
}

func Test2(t *testing.T) {
	type a struct {
		name string
	}

	type b struct {
		a *a
	}
	v := &b{
		a: &a{
			name: "test",
		},
	}
	newA := &a{
		name: "haha",
	}
	p := (*unsafe.Pointer)(unsafe.Pointer(&v.a))
	atomic.StorePointer(p, unsafe.Pointer(newA))
	fmt.Println(v.a.name)
}

func Test_Log(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	zap.L().Debug("test")

	// zap.L().Error(fmt.Sprintf("出错了 %s", uuid.New().String()), zap.Error(errors.New("数据库错误")))

	zap.L().Info("这是一条日志", zap.String("name", "张三"), zap.Int("age", 20))
}

func TestInitConfigByRemote(t *testing.T) {
	ch := make(chan string, 2)
	ch <- "123"
	ch <- "456"
	fmt.Println("哈哈好")
	close(ch)
	for v := range ch {
		fmt.Println(v)
	}
}

func Test_Time_Sub(t *testing.T) {
	ti, _ := time.ParseInLocation("2006-01-02 15:04:05", "2024-04-01 12:00:00", time.Local)
	fmt.Println(time.Now().Sub(ti).String())
}
