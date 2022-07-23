package main

import (
	. "example/part1_ioc/Injector"
	"example/part1_ioc/services"
	"fmt"
)

func main() {
	//uid := 123
	//userService := services.NewUserService(services.NewOrderService())
	//userService.GetUserInfo(uid)

	BeanFactory.Set(services.NewOrderService())
	order := BeanFactory.Get((*services.OrderService)(nil))
	fmt.Printf("%T", order)
}
