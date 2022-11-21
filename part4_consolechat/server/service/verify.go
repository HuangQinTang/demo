package service

import "chat/dao"

var VerifyService = verifyService{}

type verifyService struct {
}

// @Description 检查用户是否存在， 返回userId,或空窜（不存在）
func (s verifyService) IsExistUserName(userName string) string {
	res, _ := dao.NewUserDao().GetUserIdByUserName(userName)
	return res
}