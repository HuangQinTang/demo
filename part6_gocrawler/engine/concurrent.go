package engine

import (
	"crawler/model"
)

// ConcurrentEngine 并发版引擎
type ConcurrentEngine struct {
	Scheduler        Scheduler //调度器
	WorkerCount      int       //同时进行的任务数(干活的工人数）
	ItemChan         chan model.SimpleInfo
	RequestProcessor Processor
}

//
type Processor func(r Request) (ParseResult, error)

// Scheduler 调度器
type Scheduler interface {
	Submit(request Request)   //提交任务
	WorkerChan() chan Request //获取工作的管道
	ReadyNotifier             //准备干活
	Run()                     //开始调度
}

// ReadyNotifier 接收一个工作管道，告诉调度器我可以干活了
type ReadyNotifier interface {
	WorkerReady(chan Request)
}

// Run 引擎启动
func (e *ConcurrentEngine) Run(seeds ...Request) {
	//创建输出管道
	out := make(chan ParseResult)

	//启动调度器
	e.Scheduler.Run()

	//根据设置的任务数启动对应工作协程(创建workerCount个工作管道)
	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)
	}

	//把要干的活提交到任务管道
	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	//输出结果
	for {
		result := <-out
		for _, item := range result.Items {
			temp := item
			go func() { e.ItemChan <- temp }()
		}

		//如果还有要做的任务继续提交
		for _, request := range result.Requests {
			e.Scheduler.Submit(request)
		}
	}
}

// createWorker 创建工作goroute
func (e *ConcurrentEngine) createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			//把当前工作管道放入调度器，等待分发任务
			ready.WorkerReady(in)

			//任务分发到我(当前工作管道)头上时处理
			request := <-in

			result, err := e.RequestProcessor(request)
			if err != nil {
				continue
			}

			//把处理结果传到输出管道
			out <- result
		}
	}()
}
