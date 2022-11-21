package main

import (
	"crawler/engine"
	"crawler/persist"
	"crawler/scheduler"
	"crawler/zhenai/parser"
	"github.com/olivere/elastic/v7"
)

func main() {
	client, err := elastic.NewClient(
		//es服务不是跑在本地的，swtSniff=false, 不维护集群状态
		elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	//简单版本(仅打印拉取信息)
	//engine.SimpleEngine{}.Run(engine.Request{
	//	Url: "http://www.zhenai.com/zhenghun",
	//	Parser: parser.NewParseCityList(10),
	//})

	//并发版
	e := engine.ConcurrentEngine{
		Scheduler: &scheduler.SimpleScheduler{},
		WorkerCount: 2,	//任务数
		ItemChan: persist.SimpleInfoSeaver(client),
		RequestProcessor: engine.Worker,
	}
	e.Run(engine.Request{
		Url: "http://www.zhenai.com/zhenghun",
		Parser: parser.NewParseCityList(1),
	})

	//队列版
	//e := engine.ConcurrentEngine{
	//	Scheduler:   &scheduler.QueuedScheduler{},
	//	WorkerCount: 1,                                //任务数
	//	ItemChan:    persist.SimpleInfoSeaver(client), //处理响应结果的管道
	//	RequestProcessor: engine.Worker,
	//}
	//e.Run(engine.Request{
	//	Url:    "http://www.zhenai.com/zhenghun",
	//	Parser: parser.NewParseCityList(10),
	//})
}
