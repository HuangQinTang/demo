package Config

import "example/part1_ioc/services"

type ServiceConfig struct {
}

func NewServiceConfig() *ServiceConfig {
	return &ServiceConfig{}
}

func (s *ServiceConfig) OrderService() *services.OrderService {
	return services.NewOrderService()
}
