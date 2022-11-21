package main

import (
	"logagent/conf"
	"logagent/etcdservice"
	"logagent/kafkaservice"
	"logagent/task"
)

// 日志收集客户端(收集指定目录下的日志文件，发送到kafka中)
// 类似开源的项目有filebeat
func main() {
	//获取配置
	cfg := conf.GetCfg()

	//连接etcd
	etcdCli := etcdservice.InitEtcdService(cfg.Etcd)
	defer etcdCli.Close()

	//创建kafka生产者客户端
	kafkaCli := kafkaservice.InitProducer(cfg.Kafka.ProducerAddr)
	defer kafkaCli.Close()

	//执行收集日志任务
	t := task.Task{
		EtcdClient:  etcdCli,
		KafkaClient: kafkaCli,
	}
	t.Run()
}
