package service

import ("chat/defined"
)

var UserService = &userService{}

type userService struct{}

//登录
func (this *userService) Login(userId, userPwd string, menu *Menu) (defined.Message, error) {
	return ApiService.SendAndRead(defined.LoginMesType, map[string]interface{}{
		"user_id":  userId,
		"user_pwd": userPwd,
	}, menu)
}

//注册
func (this *userService) Register(userId, userPwd, userName string, menu *Menu) (defined.Message, error) {
	return ApiService.SendAndRead(defined.RegisterMesType, map[string]interface{}{
		"user_id":   userId,
		"user_pwd":  userPwd,
		"user_name": userName,
	}, menu)
}

//获取当前所有在线用户
func (this *userService) GetAllOnlineUser(menu *Menu) error {
	return ApiService.Send(defined.GetAllOnlineUserNameMesType, map[string]interface{}{"user_id": menu.Cache.UserId}, menu)
}

//查询好友列表
func (this *userService) GetFriendList(menu *Menu) error {
	return ApiService.Send(defined.GetFriendListMesType, map[string]interface{}{"user_id": menu.Cache.UserId}, menu)
}