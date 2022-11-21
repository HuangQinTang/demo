package client

import (
	"crawler/distributed/config"
	"crawler/distributed/worker"
	"crawler/engine"
	"net/rpc"
)

func CreateProcessor(clientChan chan *rpc.Client) engine.Processor {
	return func(req engine.Request) (engine.ParseResult, error) {
		//入参序列化
		sReq := worker.SerializeRequest(req)

		//序列化后的对象作入参，调取对应rpc服务
		var sResult worker.ParseResult
		client := <-clientChan
		err := client.Call(config.CrawlServiceRpc, sReq, &sResult)
		if err != nil {
			return engine.ParseResult{}, err
		}

		return worker.DeserializeResult(sResult), nil
	}
}
