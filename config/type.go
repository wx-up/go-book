package config

type Config struct {
	Mysql MysqlConfig
	Redis RedisConfig
}

type MysqlConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr     string
	DB       int
	Password string
}
