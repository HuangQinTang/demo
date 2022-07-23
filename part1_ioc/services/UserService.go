package services

import "fmt"

type UserService struct {
	Order *OrderService `inject:"-"`
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) GetUserInfo(uid int) {
	fmt.Println("获取用户ID=", uid, "的用户信息")
	u.Order.GetOrderInfo(uid)
}
