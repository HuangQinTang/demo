package main

import (
	"example/part1_ioc/Config"
	. "example/part1_ioc/Injector"
	"example/part1_ioc/services"
	"fmt"
)

func main() {
	serviceConfig := Config.NewServiceConfig() //获取配置对象，配置对象的各方法名为要存入容器的对象名，方法返回值为要存入容器的对象（约定一个方法返回一个值）
	BeanFactory.Config(serviceConfig)          //循环执行配置对象各方法，将方法名，方法返回的对象映射到容器中，用以后续依赖注入时，直接从容器中获取
	{
		//tag inject:"-"，单例直接从容器中取
		userService := services.NewUserService()
		BeanFactory.Apply(userService) //注入依赖
		fmt.Println(userService.Order)
	}
	{
		//tag带表达式,多例,重新示例化依赖对象放入容器
		adminService := services.NewAdminService()
		BeanFactory.Apply(adminService)
		fmt.Println(adminService.Order)
	}
}
