//go:build !k8s

package config

var C = Config{
	Redis: RedisConfig{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	},
	Mysql: MysqlConfig{
		DSN: "root:root@tcp(localhost:3306)/go_book?charset=utf8mb4&parseTime=True&multiStatements=true&loc=Local",
	},
}
