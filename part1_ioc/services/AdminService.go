package services

type AdminService struct {
	Order *OrderService `inject:"ServiceConfig.OrderService()"`
	Name  string
}

func NewAdminService() *AdminService {
	return &AdminService{}
}
