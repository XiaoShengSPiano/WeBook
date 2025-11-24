//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DNS: "root:root@tcp(webook-mysql:13309)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:11479",
	},
}
