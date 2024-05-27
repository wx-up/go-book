package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/robfig/cron/v3"

	"golang.org/x/sync/errgroup"

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
		Keys:    bson.D{{Key: "title", Value: 1}, {Key: "value", Value: -1}},
		Options: options.Index(),
		// Options: options.Index().SetUnique(true),
	})
	assert.NoError(t, err)
	fmt.Println(indexRes)

	// upsert
	filter = bson.D{{"species", "Ledebouria socialis"}, {"plant_id", 3}}
	update := bson.D{{"$set", bson.D{{"species", "Ledebouria socialis"}, {"plant_id", 3}, {"height", 8.3}}}}
	upsertResult, err := col.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	fmt.Println(upsertResult, err)
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

	fmt.Println(math.Pow(2, 12))
}

func Test_Err(t *testing.T) {
	eg := errgroup.Group{}
	eg.SetLimit(3)
	eg.TryGo(func() error {
		return errors.New("error")
	})
	eg.TryGo(func() error {
		fmt.Println("2")
		return nil
	})
	eg.TryGo(func() error {
		fmt.Println("3")
		return nil
	})
	fmt.Println(eg.Wait())
}

func Test_Channel1(t *testing.T) {
	c1 := make(chan int, 1)
	c2 := make(chan int, 1)
	c1 <- 1
	c2 <- 2
	select {
	case v1 := <-c1:
		fmt.Println(v1)
		fmt.Println(<-c2)
	case v2 := <-c2:
		fmt.Println(v2)
		fmt.Println(<-c1)
	default:
		fmt.Println("default")
	}
}

func Test_Channel(t *testing.T) {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 3
	close(ch)
	val, ok := <-ch
	fmt.Println(val, ok)
	val, ok = <-ch
	fmt.Println(val, ok)
	val, ok = <-ch
	fmt.Println(val, ok)
	val = <-ch
	fmt.Println(val)

	//select {
	//case ch <- 1:
	//default:
	//	fmt.Println("channel is full")
	//}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cancel()
	select {
	case <-ctx.Done():
		fmt.Println("context is done")
	default:
	}
}

func Test_Context(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println(ctx)
	cancel()
	fmt.Println("关闭一次")
	cancel()
	fmt.Println("关闭两次")
}

func name() string {
	return ""
}

func Test_Timer(t *testing.T) {
	tm := time.NewTimer(time.Second)
	defer tm.Stop()
	for now := range tm.C {
		t.Log(now.Unix())
	}
}

func Test_Ticker(t *testing.T) {
	tm := time.NewTicker(time.Second)
	defer tm.Stop()
	for now := range tm.C {
		t.Log(now.Unix())
	}
}

func TestTicker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ticker := time.NewTicker(time.Second)

	// 不要忘记关闭，避免潜在的goroutine泄漏
	defer ticker.Stop()

	done := false
	for !done {
		select {
		case now := <-ticker.C:
			t.Log(now.Unix())
		case <-ctx.Done():
			fmt.Println("退出了")
			done = true
		}
	}
}

func Test_Cron(t *testing.T) {
	expr := cron.New(cron.WithSeconds())
	// 线程安全
	id, err := expr.AddFunc("@every 1s", func() {
		t.Log("hello world")
	})
	require.NoError(t, err)
	t.Log(id)

	// 调用 start 之后开始调度任务
	expr.Start()
	time.Sleep(time.Second * 5)

	// 调用 Stop 之后只是暂停了后续任务的调度
	// 正在运行的任务会等待执行完成
	ctx := expr.Stop()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	select {
	case <-ctx.Done():
	case <-timeoutCtx.Done():
	}
}

func Test_Multi_Channel(t *testing.T) {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)
	go func() {
		for {
			num := 10
			select {
			case ch1 <- num:
			default:
				ch2 <- num
			}
		}
	}()

	go func() {
		for {
			num := 10
			select {
			case ch2 <- num:
			default:
				ch1 <- num
			}
		}
	}()
}
