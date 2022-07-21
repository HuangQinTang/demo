package service

import (
	"context"
	"example/part2_distributed/registry"
	"fmt"
	"log"
	"net/http"
)

// Start 启动服务
func Start(ctx context.Context, reg registry.Registration, host, port string, registerHandlersFun func()) (context.Context, error) {
	registerHandlersFun()
	ctx = startService(ctx, reg.ServiceName, host, port)
	err := registry.RegisterService(reg) //服务注册
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = host + ":" + port
	srv.Handler = nil	//DefaultServeMux(相当于路由器)

	go func() {
		//服务挂了或者关闭服务器时 -> 取消服务注册，执行cancel()
		log.Printf("%v %s", srv.ListenAndServe(), serviceName)
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		log.Println("cancel service registry success.")
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
