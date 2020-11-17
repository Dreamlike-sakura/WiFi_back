package config

import (
	"sync"
	"gopkg.in/redis.v4"
)

var (
	client    *redis.Client
	redisOnce sync.Once
)

// 创建 redis 客户端
func createClient() {
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		PoolSize: 5,
		Network:  "tcp",
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic("redis 连接失败")
	}
}

// GetRedis 得到 redis 客户端
func GetRedis() *redis.Client {
	redisOnce.Do(func() {
		createClient()
	})

	return client
}

