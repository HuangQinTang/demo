package conf

import "time"

type Config struct {
	Kafka Kafka `ini:"kafka"`
	Etcd  Etcd  `ini:"etcd"`
}

// Kafka 配置
type Kafka struct {
	ProducerAddr []string `ini:"producerAddr"`
}

// Etcd 配置
type Etcd struct {
	Endpoints   []string      `ini:"endpoints"`
	DialTimeout time.Duration `ini:"dialTimeout"`
}
