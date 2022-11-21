package process

import (
	"context"
	"encoding/json"
	"fmt"
	"chat/dao"
	"chat/defined"
	"chat/library/logicerror"
	"chat/library/usermgr"
	"chat/utils"
)

//	@Description 处理用户相关逻辑,登录、注册等
type UserProcess struct {
	Processor *Processor
}

//	@Description 登录逻辑
func (this *UserProcess) LoginProcess(ctx *context.Context, userId, userPwd string) defined.LoginResMsg {
	var loginResMes defined.LoginResMsg
	loginResMes.Code = defined.UnknowError
	loginResMes.Mes = defined.ERROR_REGISTER.Error()
	//判断登录
	userInfo, err := dao.NewUserDao().GetUserDetailById(userId)
	if err != nil {
		utils.SDD("LoginProcess redis查询报错" + err.Error())
		return defined.LoginResMsg{}
	}

	if userInfo.UserId == "" { //用户不存在
		loginResMes.Code = defined.UnknowError
		loginResMes.Mes = defined.ERROR_USER_NOTEXISTS.Error()
	}
	if userPwd == userInfo.UserPwd { //登录成功
		loginResMes.Code = defined.Success
		loginResMes.Mes = defined.SUCCESS_LOGIN
		receiveMes, err := this.ClientMesReceive(userInfo.UserId)
		if err != nil { //查询消息报错
			utils.SDD("ClientMesReceive 查询消息报错" + err.Error())
			return defined.LoginResMsg{}
		}
		friendList, err := this.GetFriendList(userId) //好友列表
		if err != nil {
			utils.SDD("ClientMesReceive 查询消息报错" + err.Error())
			return defined.LoginResMsg{}
		}
		infoList, err := this.GetInfo(userId) //信息列表
		if err != nil {
			utils.SDD("ClientMesReceive 查询消息报错" + err.Error())
			return defined.LoginResMsg{}
		}
		loginResMes.Data = defined.LoginResData{
			UserName:   userInfo.UserName,
			UserId:     userInfo.UserId,
			ReceiveMes: receiveMes,
			FriendList: friendList,
			Info:       infoList,
		}
		this.loginSuccessHook(userId, friendList) //登录成功后的逻辑
	} else {
		loginResMes.Code = defined.UnknowError
		loginResMes.Mes = defined.ERROR_USER_PWD.Error()
	}
	return loginResMes
}

//	@Description 注册逻辑
func (this *UserProcess) RegisterProcess(userId, userPwd, userName string) defined.RegisterResMsg {
	//定义响应信息
	var registerResMes defined.RegisterResMsg

	//用户id和用户昵称唯一
	if dao.NewUserDao().ExistsUser(userId) { //判断用户id是否存在
		registerResMes.Code = defined.UnknowError
		registerResMes.Mes = defined.ERROR_USERID_EXISTS.Error()
	} else if dao.NewUserDao().ExistUserName(userName) { //判断昵称是否存在
		registerResMes.Code = defined.UnknowError
		registerResMes.Mes = defined.ERROR_USERNAME_EXISTS.Error()
	} else { //都不存在创建用户
		if dao.NewUserDao().CreateUser(userId, userPwd, userName) {
			registerResMes.Code = defined.Success
			registerResMes.Mes = defined.SUCCESS_REGISTER
			registerResMes.Data = defined.RegisterMes{userId, userPwd, userName}
		} else { //创建失败返回提示
			registerResMes.Code = defined.UnknowError
			registerResMes.Mes = defined.ERROR_REGISTER.Error()
		}
	}
	return registerResMes
}

//登录成功后处理逻辑
func (this *UserProcess) loginSuccessHook(userId string, friendList []defined.FriendList) {
	//添加在线集合
	userName, err := dao.NewUserDao().GetUserNameByUserId(userId)
	if err != nil {
		utils.SDD("loginSuccessHook 用户id" + userId + "添加在线用户集合失败" + err.Error())
	} else {
		_ = dao.NewUserDao().AddOnlineUser(userName)
	}

	//userid存入当前process上下文中
	*this.Processor.Ctx = context.WithValue(*this.Processor.Ctx, "userId", userId)
	//*this.Processor.Ctx = context.WithValue(*this.Processor.Ctx, "userName", userName)

	//把当前连接放入连接管理对象
	usermgr.UserMgr.AddOnlineUser(this.Processor.Client, userId)

	//通知好友我已上线
	for _, v := range friendList {
		tf, _ := usermgr.UserMgr.GetConnByUserId(v.UserId)
		if tf != nil {
			var resMes defined.Message
			resMes.Type = defined.FriendOnlineMesType
			resMes.Data = `{"user_id":"` + userId + `"}`
			resMesByte, err := json.Marshal(resMes)
			if err != nil {
				utils.SDD("loginSuccessHook err:" + err.Error())
			}
			if err = tf.WritePkg(resMesByte); err != nil {
				fmt.Println("发送失败")
				utils.SDD("loginSuccessHook err:" + err.Error())
			}
		}
	}
}

//获取所有在线用户昵称
func (this *UserProcess) GetAllOnlineUserName(userId string) defined.GetAllOnlineUserNameResMsg {
	//获取在线用户昵称
	onlineUser, err := dao.NewUserDao().GetAllOnlineUserName()
	if err != nil {
		utils.SDD("GetAllOnlineUserName redis查询报错：" + err.Error())
		return defined.GetAllOnlineUserNameResMsg{
			Code: defined.UnknowError,
			Mes:  "查询失败，请稍后再试!",
		}
	}

	//排除自己
	self, err := dao.NewUserDao().GetUserNameByUserId(userId)
	if err != nil {
		utils.SDD("GetAllOnlineUserName redis查询报错：" + err.Error())
		return defined.GetAllOnlineUserNameResMsg{
			Code: defined.UnknowError,
			Mes:  "查询失败，请稍后再试!",
		}
	}
	onlineUser = utils.StrSliceDelete(self, onlineUser)

	//拼接消息返回
	var msg defined.GetAllOnlineUserNameResMsg
	msg.Code = defined.Success
	msg.Data = onlineUser
	msg.Mes = "查询成功！"
	return msg
}

//获取最近的20条消息，分页什么的控制台做起来太麻烦了，不折腾了
func (this *UserProcess) ClientMesReceive(userId string) ([]defined.ReceiveMes, error) {
	//查询要拉取的消息id
	mesIds, err := dao.NewMesDao().GetReceiveMes(userId)
	if err != nil {
		return []defined.ReceiveMes{}, err
	}

	//根绝消息id获取消息
	mes, err := dao.NewMesDao().GetMesByMesIds(mesIds)
	if err != nil {
		return []defined.ReceiveMes{}, err
	}

	//拼装业务数据
	keyAndReadStuas := make(map[int]int, len(mes))
	userIds := []string{}
	for k, v := range mes {
		//循环获取消息中的userId，用于批量换取userName
		userIds = append(userIds, v.ToUserId)
		userIds = append(userIds, v.FromUserId)
		//建立消息id与阅读状态的映射关系
		if v.FromUserId == userId {
			keyAndReadStuas[k] = v.FromUserReadStatus
		} else {
			keyAndReadStuas[k] = v.ToUserReadStatus
		}
	}
	userIds = utils.RemoveDuplicateElement(userIds)                     //去重
	idAndNameMap, err := dao.NewUserDao().GetUsersNameByUserId(userIds) //换取userName
	if err != nil {
		return []defined.ReceiveMes{}, err
	}
	res := make([]defined.ReceiveMes, 0, len(mes))
	for _, v := range mes {
		res = append(res, defined.ReceiveMes{
			MesId:        v.MesId,
			Type:         v.Type,
			FromUserName: idAndNameMap[v.FromUserId],
			FromUserId:   v.FromUserId,
			ToUserId:     v.ToUserId,
			ToUserName:   idAndNameMap[v.ToUserId],
			ReadStatus:   keyAndReadStuas[v.MesId],
			CreateTime:   v.UpdateTime, //这里用更新时间
		})
	}
	return res, nil
}

//查询好友列表
func (this *UserProcess) GetFriendList(userId string) ([]defined.FriendList, error) {
	//查询好友id
	friendIds, err := dao.NewUserDao().GetFriendList(userId)
	if logicerror.PrintError(err) != nil {
		return []defined.FriendList{}, err
	}

	//没有好友
	if len(friendIds) == 0 {
		return []defined.FriendList{}, nil
	}

	//用id查询昵称
	friendMap, err := dao.NewUserDao().GetUsersNameByUserId(friendIds)
	if logicerror.PrintError(err) != nil {
		return []defined.FriendList{}, err
	}

	//查询好友会话id
	friendChatMap, err := dao.NewUserDao().GetFriendChatMapp(userId)
	if logicerror.PrintError(err) != nil {
		return []defined.FriendList{}, err
	}
	//拼接业务数据
	res := make([]defined.FriendList, 0, len(friendMap))
	for fUserId, fUserName := range friendMap {
		res = append(res, defined.FriendList{
			UserId:       fUserId,
			UserName:     fUserName,
			OnlineStatus: usermgr.UserMgr.IsOnline(fUserId),
			ChatId:       friendChatMap[fUserId],
		})
	}
	return res, nil
}

//查询全部消息
func (this *UserProcess) GetInfo(myUserId string) (map[string][]defined.ChatInfo, error) {
	//查询与所有好友的会话id
	friendChatMap, err := dao.NewUserDao().GetFriendChatMapp(myUserId)
	if logicerror.PrintError(err) != nil {
		return map[string][]defined.ChatInfo{}, err
	}

	res := make(map[string][]defined.ChatInfo, len(friendChatMap))
	for _, charId := range friendChatMap {
		//查询聊天信息
		info, err := dao.MesDao.GetLastInfo(charId)
		if err != nil {
			return map[string][]defined.ChatInfo{}, err
		}
		res[charId] = info
	}
	return res, nil
}
