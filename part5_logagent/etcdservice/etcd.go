package etcdservice

import (
	"go.etcd.io/etcd/clientv3"
	"log"
	"logagent/conf"
	"time"
)

func InitEtcdService(config conf.Etcd) *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Endpoints,
		DialTimeout: time.Second * config.DialTimeout,
	})
	if err != nil {
		log.Fatalf("connect to etcd failed, err:%v", err)
	}
	return cli
}
