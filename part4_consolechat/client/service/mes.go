package service

import (
	"fmt"
	"chat/defined"
)

var MesService = &mesService{}

type mesService struct{}

//查看好友申请详情
func (this mesService) GetFriendapplyDetail(menu *Menu, mesId int) error {
	return ApiService.Send(defined.GetFriendapplyDetailMesType, map[string]interface{}{"mes_id": fmt.Sprintf("%d", mesId)}, menu)
}

//请求添加好友
func (this mesService) SendFriendApplyMes(menu *Menu, userName, remark string) error {
	return ApiService.Send(defined.AddFriendMesType, map[string]interface{}{"user_name": userName, "remark": remark}, menu)
}
