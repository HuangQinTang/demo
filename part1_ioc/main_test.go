package main

import (
	"example/part1_ioc/Config"
	"example/part1_ioc/Injector"
	"example/part1_ioc/services"
	"fmt"
	"testing"
)

func TestBeanFactory(t *testing.T) {
	//往容器存值(只能支持指针对象）
	Injector.BeanFactory.Set(&services.AdminService{Name: "测试"})
	//从容器取值
	adminService := Injector.BeanFactory.Get((*services.AdminService)(nil))
	fmt.Println("adminService", adminService) //&{<nil> 测试}

	//我们一个约定的对象（要往容器中注入N个对象，则该对象就声明N个方法，方法名就是要注入的对象名，方法返回值是这个对象的指针）格式可参照/part1_ioc/Config/ServiceConfig.go
	Injector.BeanFactory.Config(Config.NewServiceConfig()) //Config方法会通过放射遍历配置对象的方法集，把方法集中返回的对象放入容器(Set方法)
	userService := services.NewUserService()               //我们约定该结构体tag标签为inject的字段是要注入的对象,inject:"-"单例,inject:"表达式"多例
	Injector.BeanFactory.Apply(userService)                //Apply方法会遍历userService结构体的字段，把字段中tag标签为inject的字段尝试从容器中取值并通过反射赋值
	fmt.Println("userService.Order", userService.Order)
	fmt.Println("userService.IOrder", userService.IOrder)
}
