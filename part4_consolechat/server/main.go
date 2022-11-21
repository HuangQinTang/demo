package main

import (
	_ "chat/boot/serverBoot"
	"chat/library/config"
	"chat/library/transfer"
	"chat/server/process"
	"context"
	"fmt"
	"net"
)

func main() {
	fmt.Println("服务启动...")
	listen, err := net.Listen("tcp", config.GetConfig().Server.ListenPort)
	if err != nil {
		fmt.Println("net.Listen err = ", err.Error())
		return
	}
	defer listen.Close()
	fmt.Println("监听6666端口...")

	//监听成功，等待客户端链接
	for {
		fmt.Println("等待新的客户端...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err = ", err.Error())
			return
		}

		//连接成功，启动协程保持连接
		ctx := context.Background()
		p:= process.Processor{Client: &transfer.Transfer{Conn: conn}, Ctx: &ctx}
		go p.ProcessConn()
	}
}
