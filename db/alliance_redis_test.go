package db

import (
	"testing"

	"github.com/gomodule/redigo/redis"
)

func TestInitRedisConnPool(t *testing.T) {
	InitRedisConnPool()
	pool := GetRedisPool()
	conn := pool.Get()
	redis.String(conn.Do("set", "kkk", 0))
}
