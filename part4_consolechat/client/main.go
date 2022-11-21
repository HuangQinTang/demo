package main

import (
	"chat/boot/clientBoot"
	_ "chat/boot/clientBoot"
	"chat/client/callback"
	"chat/client/event"
	"chat/client/service"
	"chat/defined"
	"chat/utils"
	"fmt"
)

func main() {
	//启动菜单
	menu := service.MenuService
	menu.EventKey = defined.HomePageEvent
	menu.Response = make(chan defined.Message)
	menu.Mes = defined.MenuMes{UnRead: 0, Content: make([]defined.ReceiveMes, 0)}

	go work(menu.Response) //开启一个协程，如果菜单响应管道里有消息，转推到消息待处理管道

	for {
		//调取对应的事件（我把menu当作上下文来使用，贯穿整个生命周期）
		event.BaseEvent.ProcessHandle(menu)

		//阻塞等待响应
		mes := <-menu.Response

		//处理响应
		if err := callback.CallbackEven.ProcessEenuHandle(mes, menu); err != nil {
			utils.CDD("响应处理出错了~~" + err.Error())
			fmt.Println("程序异常~")
			break
		}
	}
}

func work(mes chan defined.Message) {
	for {
		mes <- <-clientBoot.MenuJob
	}
}
