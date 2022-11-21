package redisPool

import (
	"github.com/garyburd/redigo/redis"
	"chat/library/config"
	"sync"
)

var (
	Pool *redis.Pool
	once sync.Once
)

func InitPool(redisConfig config.Redis) {
	once.Do(func() {
		Pool = &redis.Pool{
			MaxIdle:     redisConfig.MaxIdle,                  //最大空闲连接数
			MaxActive:   redisConfig.MaxActive,                //最大连接数，0表示没有限制
			IdleTimeout: redisConfig.IdleTimeout * 1000000000, //最大空闲时间
			Dial: func() (redis.Conn, error) { //初始化连接代码
				return redis.Dial("tcp", redisConfig.Address)
			},
		}
	})
}
