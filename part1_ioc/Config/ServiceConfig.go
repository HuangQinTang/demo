package Config

import (
	"example/part1_ioc/services"
)

type ServiceConfig struct {
}

func NewServiceConfig() *ServiceConfig {
	return &ServiceConfig{}
}

func (s *ServiceConfig) OrderService() *services.OrderService {
	//fmt.Println("初始化OrderService")
	return services.NewOrderService()
}

func (s *ServiceConfig) DBService() *services.DBService {
	//fmt.Println("初始化DBService")
	return services.NewDBService()
}
