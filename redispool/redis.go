package redispool

import (
	"Chess/lib"
	"bufio"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var RedisPool *redis.Pool

type RedisConf struct {
	Host string `cfg:"host"`
	Port string `cfg:"port"`
	Pwd  string `cfg:"pwd"`
	User string `cfg:"user"`
}

var Config *RedisConf
var filename = "redis.conf"

// StartRedis 启动单线程redis
func StartRedis() (redis.Conn, error) {
	conn, err := redis.Dial("tcp",
		fmt.Sprintf("%s:%s", Config.Host, Config.Port),
		redis.DialUsername(Config.User),
		redis.DialPassword(Config.Pwd)) // 启动redis

	if err != nil {
		log.Println("redis连接失败：", err)
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

// GetConfig 获取redis配置
func GetConfig(fileName string) (*RedisConf, error) {
	fmt.Println(lib.GetRoot())
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	read := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] == '#' {
			continue
		}
		i := strings.IndexAny(line, " ")
		if i > 0 && i < len(line)-1 {
			key := strings.ToLower(line[:i])
			value := strings.TrimSpace(line[i:])
			read[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	config := &RedisConf{}
	val := reflect.ValueOf(config)
	typ := reflect.TypeOf(config)
	for i := 0; i < typ.Elem().NumField(); i++ {
		field := typ.Elem().Field(i)
		fieldVal := val.Elem().Field(i)
		key, ok := field.Tag.Lookup("cfg")
		if !ok {
			key = field.Name
		}
		value, ok := read[strings.ToLower(key)]
		if ok {
			switch field.Type.Kind() {
			case reflect.String:
				fieldVal.SetString(value)
			case reflect.Int:
				parseInt, err := strconv.ParseInt(value, 10, 64)
				if err == nil {
					fieldVal.SetInt(parseInt)
				}
			case reflect.Bool:
				v := "yes" == value
				fieldVal.SetBool(v)
			case reflect.Slice:
				if field.Type.Elem().Kind() == reflect.String {
					split := strings.Split(value, ",")
					fieldVal.Set(reflect.ValueOf(split))
				}
			}
		}
	}
	return config, nil
}

// GetSetRoom 创建房间
func GetSetRoom(conn redis.Conn, room string) {
	defer conn.Close()
	reply, _ := redis.Int(conn.Do("exist", room))
	if reply == 0 {
		_, _ = conn.Do("set", room, 1)
	} else {
		reply, _ = redis.Int(conn.Do("get", room))
		_, _ = conn.Do("set", room, reply+1)
	}
}

func Start() {
	var err error
	Config, err = GetConfig(filename)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("RedisConfig:%+v\n", Config)
	RedisPool = StartRedisPool()
}
