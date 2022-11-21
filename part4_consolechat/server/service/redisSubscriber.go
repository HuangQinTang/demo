package service

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"chat/library/redisPool"
	"chat/utils"
	"unsafe"
)

type (
	subscriber struct {
		client   redis.PubSubConn             //订阅对象
		cbMap    map[string]subscribeCallback //回调处理消息的函数
	}
	subscribeCallback func(channel, message string)
)

func NewSubscriber() subscriber {
	conn := redisPool.Pool.Get()
	subscriberObj := subscriber{
		client:   redis.PubSubConn{conn},
		cbMap:    make(map[string]subscribeCallback),
	}
	return subscriberObj
}

//订阅频道
func (s *subscriber) Subscribe(channel interface{}, cb subscribeCallback) {
	err := s.client.Subscribe(channel)
	if err != nil{
		utils.SDD("redis Subscribe error.")
	}

	s.cbMap[channel.(string)] = cb
	s.listenSubscriber()
}

//监听频道推送
func (s *subscriber) listenSubscriber() {
	for {
		switch res := s.client.Receive().(type) {	//阻塞
		case redis.Message:	//频道有消息
			channel := (*string)(unsafe.Pointer(&res.Channel))
			message := (*string)(unsafe.Pointer(&res.Data))
			s.cbMap[*channel](*channel, *message)
		case redis.Subscription:	//订阅，退订，会走这
			//如果res.Count = 0 ,退订; res.Channel是频道名
			fmt.Printf("订阅频道：%s: 订阅事件类型：%s 订阅频道数：%d\n", res.Channel, res.Kind, res.Count)
		case error:
			utils.SDD(res.Error())
			continue
		}
	}
}

//关闭redis连接
func (this *subscriber) Close() {
	this.client.Close()
}