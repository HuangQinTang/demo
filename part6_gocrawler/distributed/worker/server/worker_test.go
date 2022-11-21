package main

import (
	"crawler/distributed/config"
	"crawler/distributed/rpcsupport"
	"crawler/distributed/worker"
	"fmt"
	"testing"
	"time"
)

func TestCrawlService(t *testing.T) {
	host := fmt.Sprintf(":%d", config.WorkerPort0)
	go rpcsupport.ServerRpc(host, worker.CrawlService{})
	time.Sleep(time.Second)

	client, err := rpcsupport.NewClient(host)
	if err != nil {
		panic(err)
	}

	req := worker.Request{
		Url:"https://www.zhenai.com/zhenghun/wuhan",
		Parser: worker.SerializedParser{
			Name: config.ParseSimpleInfo,
		},
	}
	var result worker.ParseResult
	err = client.Call(config.CrawlServiceRpc, req, &result)
	if err != nil {
		t.Errorf(err.Error())
	} else {
		fmt.Println(result)
	}
}
