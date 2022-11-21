package main

import (
	"crawler/distributed/config"
	"crawler/distributed/persist"
	"crawler/distributed/rpcsupport"
	"flag"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
)

var port = flag.Int("port", 0, "the port for me to listen on")

// 存储服务
func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}
	log.Fatal(serverRpc(fmt.Sprintf(":%d", *port), config.Zhenai))
}

// serverRpc 启动itemsaver服务
// host 服务端口 esIndex 数据存储至es的索引名
func serverRpc(host string, esIndex string) error {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	fmt.Println("itemsaver start", host)
	return rpcsupport.ServerRpc(host, &persist.ItemSaverService{Client: client, EsIndex: esIndex})
}
