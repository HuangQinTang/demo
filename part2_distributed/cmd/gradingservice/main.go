package main

import (
	"context"
	"example/part2_distributed/grades"
	"example/part2_distributed/log"
	"example/part2_distributed/registry"
	"example/part2_distributed/service"
	"fmt"
	stlog "log"
)

//学生信息 server
func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	r := registry.Registration{
		ServiceName:      registry.GradingService,
		ServiceUrl:       serviceAddress,
		RequiredServices: []registry.ServiceName{registry.LogService},
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}
	ctx, err := service.Start(context.Background(),
		r,
		host,
		port,
		grades.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}

	//从服务提供者获取依赖服务信息
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("Logging service found at: %s\n", logProvider)
		log.SetClientLogger(logProvider[0], r.ServiceName)
	}

	<-ctx.Done()
	fmt.Printf("Shutting down %s.\n", registry.GradingService)
}
