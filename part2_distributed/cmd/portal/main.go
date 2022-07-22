package main

import (
	"context"
	"example/part2_distributed/log"
	"example/part2_distributed/portal"
	"example/part2_distributed/registry"
	"example/part2_distributed/service"
	"fmt"
	stlog "log"
)

func main() {
	err := portal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}
	host, port := "localhost", "5000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		ServiceName: registry.PortalService,
		ServiceUrl:  serviceAddress,
		RequiredServices: []registry.ServiceName{
			registry.LogService,
			registry.GradingService,
		},
		ServiceUpdateURL: serviceAddress + "/services",
		//HeartbeatURL: serviceAddress + "/heartbeat",
	}

	ctx, err := service.Start(context.Background(),
		r,
		host,
		port,
		portal.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err != nil {
		log.SetClientLogger(logProvider[0], r.ServiceName)
	}
	<- ctx.Done()
	fmt.Println("Shutting down portal.")
}
