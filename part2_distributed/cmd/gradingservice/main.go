package main

import (
	"context"
	"example/part2_distributed/grades"
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
		ServiceName: registry.GradingService,
		ServiceUrl:  serviceAddress,
	}
	ctx, err := service.Start(context.Background(),
		r,
		host,
		port,
		grades.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}

	<-ctx.Done()
	fmt.Printf("Shutting down %s.\n", registry.GradingService)
}
