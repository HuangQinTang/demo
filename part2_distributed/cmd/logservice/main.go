package main

import (
	"context"
	"example/part2_distributed/log"
	"example/part2_distributed/registry"
	"example/part2_distributed/service"
	"fmt"
	stlog "log"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	reg := registry.Registration{
		ServiceName: "Log Service",
		ServiceUrl:  serviceAddress,
	}
	ctx, err := service.Start(context.Background(), reg, host, port, log.RegisterHandlers)
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down log service.")
}
