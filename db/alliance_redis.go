package db

import (
	"github.com/gomodule/redigo/redis"
)

var redispool *redis.Pool

func InitRedisConnPool() {
	redispool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
}

func GetRedisPool() *redis.Pool {
	return redispool
}
