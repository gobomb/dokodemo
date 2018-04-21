package util

import (
	"github.com/go-redis/redis"
	"github.com/qiniu/log"
)

var client *redis.Client

// NewRedisClient 创建一个 Redis 实例
func newRedisClient() {
	client = redis.NewClient(&redis.Options{Addr: Configuration.RedisDB.Address, Password: Configuration.RedisDB.Password})
	log.Info(client)
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("[NewRedisClient] failed %v", err)
		return
	}
}

// RedisClient 获取 redis 实例
func RedisClient() *redis.Client {
	return client
}
