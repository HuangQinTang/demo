package gredis

import (
	"blog/pkg/setting"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os"
	"sync"
	"time"
)

var RedisConn *redis.Pool
var once sync.Once

func Setup() error {
	once.Do(func() {
		host := os.Getenv(setting.RedisSetting.Host) + ":" + setting.RedisSetting.Port
		psw := os.Getenv(setting.RedisSetting.Password)
		RedisConn = &redis.Pool{
			MaxIdle:     setting.RedisSetting.MaxIdle,
			MaxActive:   setting.RedisSetting.MaxActive,
			IdleTimeout: time.Duration(setting.RedisSetting.IdleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", host)
				if err != nil {
					fmt.Println(psw)
					fmt.Println(err.Error())
					return nil, err
				}
				if psw != "" {
					if _, err := c.Do("AUTH", psw); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
	})

	return nil
}

func Set(key string, data interface{}, time int) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	reply, err := redis.Bool(conn.Do("SET", key, value))
	conn.Do("EXPIRE", key, time)

	return reply, err
}

func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
