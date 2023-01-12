package redispool

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

var RedisPool *redis.Pool

// StartRedis 启动单线程redis
func StartRedis() (redis.Conn, error) {
	host := "r-bp1h1589rd5rbavzpppd.redis.rds.aliyuncs.com"
	pwd := "abc12345@"

	host = "127.0.0.1"
	pwd = ""
	conn, err := redis.Dial("tcp",
		fmt.Sprintf("%s:6379", host),
		redis.DialPassword(pwd)) // 启动redis
	if err != nil {
		return nil, err
	} else {
		log.Println("redis连接成功")
	}
	return conn, nil
}

// StartRedisPool 启动redis线程池
func StartRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:         100,
		MaxActive:       100,
		Wait:            true,
		MaxConnLifetime: time.Second * 60,
		Dial: func() (redis.Conn, error) {
			return StartRedis()
		},
	}
}

// GetsSetRoom 创建房间
func GetsSetRoom(conn redis.Conn, room string) {
	defer conn.Close()
	reply, _ := redis.Int(conn.Do("exist", room))
	if reply == 0 {
		_, _ = conn.Do("set", room, 1)
	} else {
		reply, _ = redis.Int(conn.Do("get", room))
		_, _ = conn.Do("set", room, reply+1)
	}
}

func init() {
	RedisPool = StartRedisPool()
}
