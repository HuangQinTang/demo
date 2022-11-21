package usermgr

import (
	"chat/defined"
	"chat/library/transfer"
	"sync"
)

var (
	UserMgr *userMgr //管理所有服务器连接
	lock    sync.RWMutex
)

type userMgr struct {
	onlineUsers map[string]*transfer.Transfer //当前在线用户连接，key为用户id，value存放用户的连接
}

func init() {
	UserMgr = &userMgr{
		onlineUsers: make(map[string]*transfer.Transfer, 1024),
	}
}

//添加在线数据
func (this *userMgr) AddOnlineUser(conn *transfer.Transfer, userId string) {
	lock.Lock()
	this.onlineUsers[userId] = conn
	lock.Unlock()
}

//删除，只删除没有关闭连接
func (this *userMgr) DeleteOnlineUser(userId string) {
	lock.Lock()
	delete(this.onlineUsers, userId)
	lock.Unlock()
}

//返回当前所有在线用户连接
func (this *userMgr) GetAllOnlineUser() map[string]*transfer.Transfer {
	return this.onlineUsers
}

//根据id返回对应连接
func (this *userMgr) GetConnByUserId(userId string) (*transfer.Transfer, error) {
	tf, ok := this.onlineUsers[userId]
	if !ok {
		return nil, defined.ERROR_USER_NOT_ONLINE
	}
	return tf, nil
}

//判断当前用户是否在线
func (this *userMgr) IsOnline(userId string) bool {
	_, ok := this.onlineUsers[userId]
	if !ok {
		return false
	}
	return true
}

//这里到时候调用要解耦，不要用dao
////返回当前在线的所有用户id以及昵称
//func (this *userMgr) GetAllOnlineUserInfo() map[string]string {
//	userId := make([]string, 0, len(this.onlineUsers))
//	for k, _ := range this.onlineUsers {
//		userId = append(userId, k)
//	}
//	res, err := dao.NewUserDao().GetUsersNameByUserId(userId)
//	if err != nil {
//		return map[string]string{}
//	}
//	return res
//}
