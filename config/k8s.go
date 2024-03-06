//go:build k8s

package config

var C = Config{
	Redis: RedisConfig{
		Addr:     "redis:6379",
		DB:       0,
		Password: "",
	},
	Mysql: MysqlConfig{
		DSN: "root:root@tcp(mysql:3306)/go_book?charset=utf8mb4&parseTime=True&multiStatements=true&loc=Local",
	},
}
