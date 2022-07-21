package main

import (
	"context"
	"example/part2_distributed/registry"
	"fmt"
	"log"
	"net/http"
)

// 服务中心 server（注册服务，取消注册服务）
func main() {
	http.Handle("/services", &registry.RegistryService{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = registry.ServerPort

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		log.Println("Registry service started. Press any Key to stop.")
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()
	log.Println("Shutting down registry service")
}
