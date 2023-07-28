package dao

import (
	"github.com/go-redis/redis"
)

var Client *redis.Client

func RedisInit() {
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := Client.Ping().Result()
	if err != nil {
		panic("Redis连接失败")
	}
}
