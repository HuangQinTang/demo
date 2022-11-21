package serverBoot

import (
	"chat/dao"
	"chat/library/config"
	"chat/library/redisPool"
	"chat/server/service"
)

//全局初始化
func init() {
	redisPool.InitPool(config.GetConfig().Redis) //初始化redis连接池，写代码要注意不要关了缓存池，因为所有的服务基于缓存池，且木有重新建立机制，崩了就凉凉
	daoInit()                                    //注入redis连接池进dao
	go service.NsqService.InitMesConsumerNsq()   //nsq 监听消息
}

//初始化dao，redis连接池必须先初始化(待优化)
func daoInit() {
	dao.NewUserDao()
	dao.NewMesDao()
}
