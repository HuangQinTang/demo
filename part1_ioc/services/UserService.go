package services

import "fmt"

type UserService struct {
	order *OrderService
}

func NewUserService(order *OrderService) *UserService {
	return &UserService{order: order}
}

func (u *UserService) GetUserInfo(uid int) {
	fmt.Println("获取用户ID=", uid, "的用户信息")
	u.order.GetOrderInfo(uid)
}
