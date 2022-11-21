package rediservice

import (
	"github.com/garyburd/redigo/redis"
)

type redisService struct {
	pool *redis.Pool
}

var (
	pool        *redis.Pool
	RedisServer = redisService{}
)

const DuplicateKey = "duplicate:zhenai:userid:"

func init() {
	//初始化redis连接池
	pool = &redis.Pool{
		MaxIdle:     500,  //最大空闲连接数
		MaxActive:   2000, //最大连接数，0表示没有限制（redis默认最大连接数10000）
		IdleTimeout: 120,  //最大空闲时间(秒)
		Dial: func() (redis.Conn, error) { //初始化连接代码
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	RedisServer.pool = pool
}

// ZhenaiDuplicate 珍爱网userid去重
// return true已设置 false未设置
func (p *redisService) ZhenaiDuplicate(userId string) (bool, error) {
	conn := p.pool.Get()
	defer conn.Close()

	//md5 := md5.Sum([]byte(uniqueKey))
	//md5Str := fmt.Sprintf("%x", md5)

	result, err := redis.Int(conn.Do("SETNX", DuplicateKey+userId, 1))
	if err != nil {
		return false, err
	}
	if result == 1 {
		return true, nil
	}
	return false, nil
}
