package scheduler

import "crawler/engine"

type SimpleScheduler struct {
	workerChan chan engine.Request
}

// WorkerChan 返回工作管道
func (s *SimpleScheduler) WorkerChan() chan engine.Request {
	return s.workerChan
}

// WorkerReady 简单爬虫这是摆设
func (s *SimpleScheduler) WorkerReady(requests chan engine.Request) {}

// Run 创建工作管道
func (s *SimpleScheduler) Run() {
	s.workerChan = make(chan engine.Request)
}

// Submit 把任务提交到工作的管道
func (s *SimpleScheduler) Submit(r engine.Request) {
	//避免循环阻塞，这里用协程
	go func() { s.workerChan <- r }()
}
