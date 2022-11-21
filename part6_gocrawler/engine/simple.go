package engine

import (
	"log"
)

type SimpleEngine struct{}

// Run 引擎,相当于调度中心
// 接收种子sessds,存放到任务队列requests,循环送到fetcher
// fetch 请求种子指定的资源, 放入解析器parser解析，得到数据
func (e SimpleEngine) Run(sessds ...Request) {
	var requests []Request
	for _, r := range sessds {
		requests = append(requests, r)
	}
	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]

		parseResult, err := Worker(r)
		if err != nil {
			continue
		}
		requests = append(requests, parseResult.Requests...)

		for _, item := range parseResult.Items {
			log.Printf("Got item %v", item)
		}
	}
}


