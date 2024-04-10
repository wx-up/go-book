package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/event"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

func Test_JsonMarshal(t *testing.T) {
	type Dog struct {
		Name string
	}
	var d *Dog
	_ = json.Unmarshal([]byte(`{"name":"wx"}`), &d)
	fmt.Println(d)
}

func Test_MongoDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// monitor 设置监控，用于调试
	monitor := &event.CommandMonitor{
		// 每个命令执行前调用
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		// 执行成功
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
		},
		// 执行失败
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
		},
	}
	opts := options.Client().ApplyURI("mongodb://root:root@localhost:27017").SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	// 不用预先创建数据库和集合，直接插入数据即可
	col := client.Database("go_book").Collection("articles")

	// 默认文档的字段是结构体字段的小写，AuthorID --> authorid
	// 可以使用 bson 这个标签来自定义字段名
	res, err := col.InsertOne(ctx, Article{
		Title:   "我是标题",
		Content: "我是内容",
	})
	assert.NoError(t, err)

	// mongoDB是没有自增主键的，这个是文档ID，mongoDB会为每个文档生成一个 _id 字段
	fmt.Println(res.InsertedID)

	// 使用结构体来构建查询条件，需要注意零值问题
	// 可以使用 `bson:"author_id,omitempty"` omitempty 标签来忽略零值
	findRes := col.FindOne(ctx, Article{Title: "我是标题"})
	fmt.Println(findRes.Raw())

	// 是否有找到
	if findRes.Err() == mongo.ErrNoDocuments {
		fmt.Println("没有找到")
	}

	// 使用 bson 来构建查询条件
	filter := bson.D{{"title", "我是标题"}}
	findRes = col.FindOne(ctx, filter)
	fmt.Println(findRes.Raw())
	var art Article
	assert.NoError(t, findRes.Decode(&art))
	fmt.Println(art)

	// 更新数据
	// 构建 filter （ 更新条件 ）
	// 构建 update （ 更新内容 ）
	// 同查询，如果使用结构体的话，需要注意零值问题
	filter = bson.D{{Key: "title", Value: "我是标题"}}
	// $set 表示操作符，表示更新内容（ mongodb 的规范 ）
	sets := bson.D{{Key: "$set", Value: bson.E{Key: "author_id", Value: 999}}}
	updateRes, err := col.UpdateMany(ctx, filter, sets)
	assert.NoError(t, err)
	fmt.Println(updateRes.ModifiedCount)
	// 使用结构体更新
	updateRes, err = col.UpdateMany(ctx, filter, bson.M{"$set": Article{Title: "哈哈哈哈哈"}})
	assert.NoError(t, err)
	fmt.Println(updateRes.ModifiedCount)

	// 删除数据
	// delRes, err := col.DeleteMany(ctx, bson.M{"title": "哈哈哈哈哈"})
	// assert.NoError(t, err)
	// fmt.Println(delRes.DeletedCount)

	// or 查询
	findAllRes, err := col.Find(ctx, bson.M{"$or": bson.A{bson.M{"title": "哈哈哈哈哈"}, bson.M{"value": 888}}})
	assert.NoError(t, err)
	arts := make([]Article, 0, 2)
	err = findAllRes.All(ctx, &arts)
	assert.NoError(t, err)
	println(len(arts))

	// 创建索引
	indexRes, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: 2}, {Key: "value", Value: 2}},
		Options: options.Index(),
		// Options: options.Index().SetUnique(true),
	})
	assert.NoError(t, err)
	fmt.Println(indexRes)
}

type Article struct {
	Title    string
	Content  string
	AuthorId int64 `bson:"author_id,omitempty"`
}

func Test_Number(t *testing.T) {
	fmt.Println(time.Now().UnixMilli())
	fmt.Println(int64(0b0000000000000000000000011111111111111111111111111111111111111111))
	n := int64(-1)
	fmt.Println(strconv.FormatInt(n, 2))
}
