package nsq

import (
	"github.com/nsqio/go-nsq"
	"chat/library/config"
	"time"
)

//消费者配置
type NsqConsumerConfig struct {
	Topic   string
	Channel string
}

//消费者对象
type NsqConsumer struct {
	consumer *nsq.Consumer
}

//返回一个消费者
func (n NsqConsumerConfig) NewConsumer(f nsq.Handler) (*NsqConsumer, error) {
	nsgConfig := nsq.NewConfig()
	nsgConfig.LookupdPollInterval = time.Second*30 //集群时，表示每30秒从nsqlookupd获取最新的nsqd地址；非集群时表示与nsqd连接中断后隔多久尝试重连
	consumer, err := nsq.NewConsumer(n.Topic, n.Channel, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	//根据nsqds数量来配置，集群情况下，如果值为1，会出现其他nsqd单点的消息需要很长时间被消费
	consumer.ChangeMaxInFlight(config.GetConfig().Nsq.MaxInFlight)
	consumer.AddHandler(f)
	return &NsqConsumer{consumer: consumer}, nil
}

//消费者集群调用
func (c *NsqConsumer) ConnLookupd() error {
	return c.consumer.ConnectToNSQLookupd(config.GetConfig().Nsq.LookupAddr)
}

//消费者直连nsqd
func (c *NsqConsumer) ConnNsqd(nsqAddr string) error {
	return c.consumer.ConnectToNSQD(nsqAddr)
}

//终止连接
func (c *NsqConsumer) Stop() {
	c.consumer.Stop()
}
