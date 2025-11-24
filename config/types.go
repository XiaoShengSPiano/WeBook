package config

// 数据库配置
type DBConfig struct {
	DNS string
}

// Redis配置
type RedisConfig struct {
	Addr string
}

// 全局配置
type config struct {
	DB    DBConfig
	Redis RedisConfig
}
