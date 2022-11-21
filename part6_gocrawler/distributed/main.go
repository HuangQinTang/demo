package main

import (
	itemServer "crawler/distributed/persist/client"
	"crawler/distributed/rpcsupport"
	worker "crawler/distributed/worker/client"
	"crawler/engine"
	"crawler/scheduler"
	"crawler/zhenai/parser"
	"flag"
	"log"
	"net/rpc"
	"strings"
)

var (
	itemSaveHost = flag.String("itemsaver_host", "", "itemsaver host")
	workerHosts  = flag.String("worker_hosts", "", "worker hosts (comma separated)")
)

func main() {
	flag.Parse()
	//存储服务客户端
	itemChan, err := itemServer.ItemSeaver(*itemSaveHost)
	if err != nil {
		panic(err)
	}

	//爬取服务客户端连接池
	pool := createClientPool(strings.Split(*workerHosts, ","))
	processor := worker.CreateProcessor(pool)

	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      2, //任务数
		ItemChan:         itemChan,
		RequestProcessor: processor,
	}

	//启动调度器，调度器负责分发任务给客户端执行对应的服务
	e.Run(engine.Request{
		Url:    "http://www.zhenai.com/zhenghun",
		Parser: parser.NewParseCityList(470),
	})
}

func createClientPool(hosts []string) chan *rpc.Client {
	var clients []*rpc.Client
	for _, h := range hosts {
		client, err := rpcsupport.NewClient(h)
		if err == nil {
			clients = append(clients, client)
			log.Printf("connected to %s", h)
		} else {
			log.Printf("errr connctingto %s: %v", h, err)
		}
	}
	out := make(chan *rpc.Client)
	go func() {
		for { //循环分发
			for _, client := range clients {
				out <- client
			}
		}
	}()
	return out
}
