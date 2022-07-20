package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// Start 启动服务
func Start(ctx context.Context, serviceName, host, port string, registerHandlersFun func()) (context.Context, error) {
	registerHandlersFun()
	ctx = startService(ctx, serviceName, host, port)
	return ctx, nil
}

func startService(ctx context.Context, serviceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = host + ":" + port

	go func() {
		//服务挂了 cancel()
		log.Printf("%v %s", srv.ListenAndServe(), serviceName)
		cancel()
	}()

	go func() {
		//控制台输入任意键cancel()
		log.Printf("%v started. Press any key to stop. \n", serviceName)
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()
	return ctx
}
