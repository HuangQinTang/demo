package main

import (
	"example/part1_ioc/Config"
	. "example/part1_ioc/Injector"
	"example/part1_ioc/services"
	"fmt"
)

func main() {
	//放容器存值，再通过类型从容器取值
	//BeanFactory.Set(services.NewOrderService())
	//order := BeanFactory.Get((*services.OrderService)(nil))
	//fmt.Printf("%T", order)


	//非表达式注入获取对象(inject:"-")
	//BeanFactory.Set(services.NewOrderService())
	//userService := services.NewUserService()
	//BeanFactory.Apply(userService)
	//fmt.Println(userService.Order2)


	//使用表达式的方式注入依赖(inject:"ServiceConfig.OrderService()")，需要在ExprMap维护表达式与对应的对象
	serviceConfig := Config.NewServiceConfig()
	BeanFactory.ExprMap = map[string]interface{}{
		"ServiceConfig":serviceConfig,
	}
	userService := services.NewUserService()
	BeanFactory.Apply(userService)
	fmt.Println(userService.Order)
}
