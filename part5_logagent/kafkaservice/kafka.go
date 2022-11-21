package kafkaservice

import (
	"github.com/Shopify/sarama"
	"log"
)

// InitProducer 初始化kafka生产者对象
func InitProducer(addrs []string) sarama.SyncProducer {
	//1.生产者配置
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll          //客户端,主分片,副本都分别ack
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner //随机分片（每个分片消息是按顺序存储的,随机的话我理解消息不能保证顺序存储）
	cfg.Producer.Return.Successes = true                   //成功交付的信息

	//2.链接kafka
	client, err := sarama.NewSyncProducer(addrs, cfg)
	if err != nil {
		log.Fatalf("connect kafka failed err: %s", err.Error())
	}
	return client
}
