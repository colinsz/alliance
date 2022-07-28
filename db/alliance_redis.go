package db

import (
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
)

var redispool *redis.Pool

func InitRedisConnPool() {
	addr := os.Getenv("REDIS_URL")
	fmt.Println("init redis: ", addr)
	redispool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}
}

func GetRedisPool() *redis.Pool {
	return redispool
}
