package process

import (
	"chat/dao"
	"chat/defined"
	"chat/library/transfer"
	"chat/library/usermgr"
	"chat/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type Processor struct {
	Client *transfer.Transfer
	Ctx    *context.Context
}

func (this *Processor) ProcessConn() {
	fmt.Println("用户 " + this.Client.Conn.RemoteAddr().String() + " 进来了...")
	defer this.Client.Close()

	//心跳结构体
	heartbeat, _ := json.Marshal(defined.Message{Type: defined.ServerBeatMesType})

	for {
		//读取客户端数据
		mes, err := this.Client.ReadPkg()
		if this.checkError(err) { //判断是否读取成功
			return
		}

		//心跳回应
		if mes.Type == defined.ClientBeatMesType {
			this.Client.WritePkg(heartbeat)
			continue
		}

		//不同客户端数据执行不同逻辑
		err = this.serverProcessMes(&mes)
		if this.checkError(err) { //判断是否发送成功
			return
		}
	}
}

//判断消息类型，走不同逻辑
func (this *Processor) serverProcessMes(mes *defined.Message) (err error) {
	switch mes.Type {
	case defined.LoginMesType: //登陆
		//解析客户端信息
		var loginMes defined.LoginMes
		if err = json.Unmarshal([]byte(mes.Data), &loginMes); err != nil {
			return
		}
		processObj := &UserProcess{Processor: this}
		data := processObj.LoginProcess(this.Ctx, loginMes.UserId, loginMes.UserPwd)
		return this.send(defined.LoginResMesType, data)

	case defined.RegisterMesType: //注册
		var registerMes defined.RegisterMes
		if err = json.Unmarshal([]byte(mes.Data), &registerMes); err != nil {
			return err
		}
		processObj := &UserProcess{Processor: this}
		data := processObj.RegisterProcess(registerMes.UserId, registerMes.UserPwd, registerMes.UserName)
		return this.send(defined.RegisterResMesType, data)

	case defined.GetAllOnlineUserNameMesType: //获取所有在线用户昵称
		var getOnlineUserNameMes defined.GetAllOnlineUserNameMes
		if err = json.Unmarshal([]byte(mes.Data), &getOnlineUserNameMes); err != nil {
			return err
		}
		processObj := &UserProcess{Processor: this}
		data := processObj.GetAllOnlineUserName(getOnlineUserNameMes.UserId)
		return this.send(defined.GetAllOnlineUserNameResMesType, data)

	case defined.GetFriendapplyDetailMesType: //获取好友申请消息详情
		var getFriendApplyDetailMes defined.GetFriendApplyDetailMes
		if err = json.Unmarshal([]byte(mes.Data), &getFriendApplyDetailMes); err != nil {
			return err
		}
		processObj := &MesProcess{Processor: this}
		//查询数据
		data := processObj.GetFriendApplyDetail(getFriendApplyDetailMes.MesId)
		return this.send(defined.GetFriendapplyDetailResMesType, data)

	case defined.AddFriendMesType: //添加好友申请
		var addFreindMes defined.AddFriendMes
		if err = json.Unmarshal([]byte(mes.Data), &addFreindMes); err != nil {
			return err
		}
		processObj := &MesProcess{Processor: this}
		data := processObj.AddFriendApply(addFreindMes.UserName, addFreindMes.Remark)
		return this.send(defined.AddFreindResMesType, data)

	case defined.UpdateFriendMesType: //更新好友申请
		var updateFriendMes defined.UpdateFriendMes
		if err = json.Unmarshal([]byte(mes.Data), &updateFriendMes); err != nil {
			return err
		}
		processObj := &MesProcess{Processor: this}
		if err = processObj.UpdateFriend(updateFriendMes.MesId, updateFriendMes.Status, updateFriendMes.Remark); err != nil {
			this.send(defined.UpdateFriendResMesType, defined.FailResponse)
			return err
		}
		return this.send(defined.UpdateFriendResMesType, defined.OkResponse)

	case defined.JoinChatRoomMesType: //进入聊天室
		var joinChatMes defined.JoinChatRoomMes
		if err = json.Unmarshal([]byte(mes.Data), &joinChatMes); err != nil {
			return err
		}
		ctx := *this.Ctx
		myUserId, _ := ctx.Value("userId").(string)
		ChatRoom.Hello(myUserId)
		ChatRoom.Join(joinChatMes.UserId, this.Client)
		return this.send(defined.JoinChatRoomResMesType, defined.OkResponse)

	case defined.ChatRoomMesType: //聊天室发送消息
		var chatMes defined.ChatMes
		if err = json.Unmarshal([]byte(mes.Data), &chatMes); err != nil {
			return err
		}
		ChatRoom.SendChat(chatMes)
		return nil

	case defined.PtoPChatMesType:	//点对点聊天
		var chatMes defined.PtoPChatMes
		if err = json.Unmarshal([]byte(mes.Data), &chatMes); err != nil {
			return err
		}
		processObj := &ChatProcess{Processor: this}
		if err = processObj.SendPtoPMes(chatMes); err != nil {
			return err
		}
		return nil

	case defined.UpdatePtoPStatusMesType: //更新信息阅读状态
		var updateData defined.UpdateInfoStatus
		if err = json.Unmarshal([]byte(mes.Data), &updateData); err != nil {
			return err
		}
		processObj := &ChatProcess{Processor: this}
		if err = processObj.UpdateInfoStatus(updateData.ChatId, updateData.InfoId); err != nil {
			return err
		}
		return nil

	default:
		return errors.New("类型定义错误！")
	}
}

//检查连接状态，如果传入的错误对象是ERROR_CONN_LOST，表示断开连接并返回true
func (this *Processor) checkError(err error) (result bool) {
	if err != nil { //报错
		utils.SDD(err.Error())
		if err == defined.ERROR_CONN_LOST {
			fmt.Println("用户 " + this.Client.Conn.RemoteAddr().String() + " 断开连接...")
			//关闭连接对象
			this.Client.Close()
			//关闭连接管理对象(客户端登录成功后，会把userId存储在上下文中并把连接放入连接管理对象，连接断开时，如果存在userId，删除该userId在管理对象的映射连接)
			ctx := *this.Ctx
			userId, ok := ctx.Value("userId").(string)
			if ok {
				usermgr.UserMgr.DeleteOnlineUser(userId)
			}
			//移除在线用户集合
			userName, _ := dao.NewUserDao().GetUserNameByUserId(userId)
			dao.NewUserDao().RemoveOnlineUser(userName)
			result = true
		} else { //普通异常，返回友好提示
			fmt.Println("服务异常~")
		}
	}
	return result
}

//发送响应数据到客户端
func (this *Processor) send(mesType string, data interface{}) error {
	var resMes defined.Message
	resMes.Type = mesType
	dataJson, err := json.Marshal(data)
	if err != nil {
		utils.SDD("出大问题了，客户端响应send函数直接转空接口失败！" + err.Error())
		return err
	}
	resMes.Data = string(dataJson)
	//Marshal报错的情况 -> 不支持的标准类型有 Complex64 ，Complex128 ，Chan ，Func ，UnsafePointer ，这种情况下会返回 UnsupportedTypeError 。对于不支持的数据类型，需要实现 MarshalJSON 或者 encoding.TextMarshaler 接口。对于不支持的值，会返回 UnsupportedValueError 错误，如浮点数的无穷大，无穷小，NaN 和出现循环引用的 map、slice和pointer。
	dataJson, _ = json.Marshal(resMes)
	return this.Client.WritePkg(dataJson)
}
