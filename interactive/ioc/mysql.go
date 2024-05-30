package ioc

import (
	"fmt"
	"gorm.io/plugin/opentelemetry/tracing"
	"time"

	prometheusClient "github.com/prometheus/client_golang/prometheus"

	"go.uber.org/zap"

	"go.uber.org/zap/zapcore"

	glogger "gorm.io/gorm/logger"

	"github.com/spf13/viper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateMysql() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var c Config
	err := viper.UnmarshalKey("db.mysql", &c)
	if err != nil {
		panic(fmt.Errorf("初始化数据库配置失败：%w", err))
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: c.DSN,
	}), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(zap.L().Info), glogger.Config{
			SlowThreshold:             time.Millisecond * 100, // 慢 SQL 阈值（ insert、select 等语句 ）
			Colorful:                  true,                   // 彩色日志
			IgnoreRecordNotFoundError: false,                  // 是否忽略找不到记录的错误
			ParameterizedQueries:      false,                  // 参数值是否填充到 sql 语句中，true 不填充，线上建议是用 true，填充操作是有性能损耗的
			LogLevel:                  glogger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}
	//err = db.Use(prometheus.New(prometheus.Config{
	//	DBName:          "go_book",
	//	RefreshInterval: 15,    // gorm 插件本身采集的时间间隔
	//	StartServer:     false, // 是否启用一个 http 服务来暴露指标
	//
	//	// prometheus 的话一般不需要配置，因此 prometheus 主动会拉取
	//	PushUser:     "",
	//	PushAddr:     "",
	//	PushPassword: "",
	//
	//	MetricsCollector: []prometheus.MetricsCollector{
	//		&prometheus.MySQL{
	//			VariableNames: []string{"Threads_running"},
	//		},
	//	},
	//}))
	//if err != nil {
	//	panic(fmt.Errorf("gorm 接入 proemtheus 插件失败：%w", err))
	//}

	// 使用callback机制检测SQL的执行
	// db.Use(newExecTimeCallback())

	db.Use(tracing.NewPlugin(
		tracing.WithDBName("go_book"),
		// 不记录 metric
		tracing.WithoutMetrics(),
		// 不记录查询参数值，即 SELECT * FROM  users WHERE name = ？ 语句中 ？的值
		// 因为参数值有可能包含敏感信息
		tracing.WithoutQueryVariables(),
	))

	//_ = model.InitTables(db)
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(200)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)
	sqlDB.SetMaxOpenConns(1200)
	return db
}

// ExecTimeCallback 统计执行时间的 callback
type ExecTimeCallback struct {
	vectorSummary *prometheusClient.SummaryVec
}

func (ec *ExecTimeCallback) Name() string {
	return "gorm:prometheus-execTime"
}

func (ec *ExecTimeCallback) Initialize(db *gorm.DB) error {
	ec.registerAll(db)
	return nil
}

func newExecTimeCallback() *ExecTimeCallback {
	vectorSummary := prometheusClient.NewSummaryVec(prometheusClient.SummaryOpts{
		Namespace: "wx",
		Subsystem: "go_book",
		Name:      "gorm",
		Help:      "sql执行指标",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.99:  0.005,
			0.999: 0.0001,
		},
	}, []string{"type", "table"})

	// 这一句千万不能漏
	prometheusClient.MustRegister(vectorSummary)

	return &ExecTimeCallback{
		vectorSummary: vectorSummary,
	}
}

func (ec *ExecTimeCallback) registerAll(db *gorm.DB) {
	// insert 语句
	err := db.Callback().Create().Before("*").Register("prometheus:before_create", ec.before("create"))
	if err != nil {
		panic(fmt.Errorf("gorm 注册 prometheus:before_create 回调失败：%w", err))
	}
	err = db.Callback().Create().After("*").Register("prometheus:after_create", ec.after("create"))
	if err != nil {
		panic(fmt.Errorf("gorm 注册 prometheus:after_create 回调失败：%w", err))
	}
	// update 语句
	err = db.Callback().Update().Before("*").Register("prometheus:before_update", ec.before("update"))
	if err != nil {
		panic(fmt.Errorf("gorm 注册 prometheus:before_update 回调失败：%w", err))
	}
	err = db.Callback().Update().After("*").Register("prometheus:after_update", ec.after("update"))
	if err != nil {
		panic(fmt.Errorf("gorm 注册 prometheus:after_update 回调失败：%w", err))
	}
}

func (ec *ExecTimeCallback) before(typ string) func(tx *gorm.DB) {
	return func(tx *gorm.DB) {
		startTime := time.Now()
		tx.Set("start_time", startTime)
	}
}

func (ec *ExecTimeCallback) after(typ string) func(tx *gorm.DB) {
	return func(tx *gorm.DB) {
		val, _ := tx.Get("start_time")
		startTime, ok := val.(time.Time)
		if !ok {
			return
		}

		table := tx.Statement.Table
		// 绕过orm的查询比如exec或者raw时 db.Statement.Table 有可能会空
		if table == "" {
			table = "unknown"
		}
		duration := time.Since(startTime).Milliseconds()

		ec.vectorSummary.WithLabelValues(typ, table).Observe(float64(duration))
	}
}

type gormLoggerFunc func(msg string, fields ...zapcore.Field)

func (f gormLoggerFunc) Printf(m string, args ...interface{}) {
	f(fmt.Sprintf(m, args...))
}

//func CreateDBProvider(db *gorm.DB) dao.DBProvider {
//	viper.OnConfigChange(func(in fsnotify.Event) {
//		type Config struct {
//			DSN string `yaml:"dsn"`
//		}
//		var c Config
//		err := viper.UnmarshalKey("db", &c)
//		if err != nil {
//			return
//		}
//		newDB, err := gorm.Open(mysql.New(mysql.Config{
//			DSN: c.DSN,
//		}), &gorm.Config{})
//		if err != nil {
//			return
//		}
//		// 原子操作
//		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&db)), unsafe.Pointer(newDB))
//	})
//	return func() *gorm.DB {
//		return db
//	}
//}
