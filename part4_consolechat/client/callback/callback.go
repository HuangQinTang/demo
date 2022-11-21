package callback

import (
	"encoding/json"
	"fmt"
	"chat/client/service"
	"chat/defined"
	"chat/utils"
	"strconv"
	"time"
)

//	@Description 判断服务器的响应结果并返回菜单对象（包含下一步需要执行的对象）
var CallbackEven = callback{}

type callback struct{}

//处理菜单任务响应
func (this *callback) ProcessEenuHandle(response defined.Message, menu *service.Menu) (err error) {
	switch response.Type {
	case defined.LoginResMesType: //登录响应
		var res defined.LoginResMsg
		if err = json.Unmarshal([]byte(response.Data), &res); err != nil {
			return err
		}
		fmt.Println(res.Mes)
		if res.Code == defined.Success { //成功
			menu.Loop = true
			menu.EventKey = defined.LoginSuccessEvent //登录成功后的事件
			menu.Cache.UserId = res.Data.UserId
			menu.Cache.UserName = res.Data.UserName
			menu.FriendList = res.Data.FriendList //好友列表
			menu.InfoList = res.Data.Info         //聊天记录列表
			menu.ChatStatus = make(map[string]int, len(res.Data.Info))
			for chatId, infos := range res.Data.Info { //遍历会话统计未读消息
				for _, v := range infos { //遍历信息
					if v.ReadStatus == 0 { //未读信息
						menu.ChatStatus[chatId]++
					}
				}
			}
			menu.FriendChatId = make(map[string]string, len(res.Data.FriendList))
			for _, v := range res.Data.FriendList { //建立好友与会话id映射
				menu.FriendChatId[v.UserId] = v.ChatId
			}
			for _, v := range res.Data.ReceiveMes { //消息列表
				if v.ReadStatus == 0 {
					menu.Mes.UnRead++
				}
				menu.Mes.Content = append(menu.Mes.Content, defined.ReceiveMes{
					MesId:        v.MesId,
					Type:         v.Type,
					FromUserId:   v.FromUserId,
					FromUserName: v.FromUserName,
					ToUserName:   v.ToUserName,
					ToUserId:     v.ToUserId,
					ReadStatus:   v.ReadStatus,
					CreateTime:   v.CreateTime,
				})
			}
		} else if res.Code == defined.UnknowError { //登录失败
			fmt.Println("登录失败")
			menu.EventKey = defined.HomePageEvent //主菜单
		}

	case defined.RegisterResMesType: //注册响应
		var res defined.RegisterResMsg
		if err = json.Unmarshal([]byte(response.Data), &res); err != nil {
			return err
		}
		fmt.Println(res.Mes)

		if res.Code == defined.Success { //成功
			fmt.Println("您的账号密码依次为~")
			fmt.Println(res.Data.UserId)
			fmt.Println(res.Data.UserPwd)
			menu.EventKey = defined.HomePageEvent //回到首页
		} else if res.Code == defined.UnknowError {
			fmt.Println("请重新注册~")
			fmt.Println("返回首页~")
			menu.EventKey = defined.HomePageEvent //回到首页
		}

	case defined.GetAllOnlineUserNameResMesType: //查看所有在线用户昵称
		var res defined.GetAllOnlineUserNameResMsg
		if err = json.Unmarshal([]byte(response.Data), &res); err != nil {
			return err
		}
		if res.Code == defined.Success {
			fmt.Println("当前有", len(res.Data), "个用户在线, 如下:")
			for _, v := range res.Data {
				fmt.Println(v)
			}
			fmt.Println()
			fmt.Println("--回车返回上一步--")
			fmt.Scanln()
		} else if res.Code == defined.UnknowError {
			fmt.Println(res.Mes)
		}
		menu.EventKey = defined.MainPageEvent //回到主界面

	case defined.GetFriendapplyDetailResMesType: //好友申请详情展示
		var res defined.GetFriendApplyDetailResMes
		if err = json.Unmarshal([]byte(response.Data), &res); err != nil {
			return err
		}
		if res.Code == 1 {
			//修改阅读状态
			for k, v := range menu.Mes.Content {
				if v.MesId == res.Data.MesId {
					if v.ReadStatus == 0 {
						menu.Mes.UnRead--
					}
					menu.Mes.Content[k].ReadStatus = 1
				}
			}
			//打印消息内容
			fmt.Println("-----详情-----")
			fmt.Println("用户：", res.Data.FromUserName, " 请求添加好友")
			for _, v := range res.Data.Remark {
				fmt.Println(v.UserId + ":" + v.Content + "(" + time.Unix(int64(v.Time), 0).Format("2006-01-02 15:04:05") + ")")
			}
			//已是好友不需要走下面逻辑
			if res.Data.Status == 1 {
				fmt.Println("你们已成为好朋友啦！")
				fmt.Println("----回车返回上一步----")
				fmt.Scanln()
				menu.EventKey = defined.MainPageEvent //回到主界面
				return
			}

			//非好友进一步操作
			if res.Data.FromUserId == menu.Cache.UserId { //我发送的添加好友消息
				var remark string
				utils.DumpSAndScanVar("继续留言, q退出", &remark)
				if remark == "q" {
					menu.EventKey = defined.MainPageEvent
					return
				}
				menu.Request.Type = defined.UpdateFriendMesType
				menu.Request.Data = map[string]interface{}{
					"mes_id": res.Data.MesId,
					"status": 2,
					"remark": remark,
				}
				menu.EventKey = defined.SendCommonRequestEvent
			} else { //别人请求添加我为好友的消息
			GetFriendapplyDetailResMesType: //goto标志
				var status string
				var remark string
				utils.DumpSAndScanVar("1同意添加好友  0拒绝,  q退出", &status)
				if status == "q" {
					menu.EventKey = defined.MainPageEvent
					return
				}
				if status == "0" {
					utils.DumpSAndScanVar("留言, 回车跳过留言", &remark)
				}
				if status != "q" && status != "1" && status != "0" {
					fmt.Println("老铁非法操作丢！")
					goto GetFriendapplyDetailResMesType
				}
				menu.Request.Type = defined.UpdateFriendMesType
				statusInt, _ := strconv.Atoi(status)
				menu.Request.Data = map[string]interface{}{
					"mes_id": res.Data.MesId,
					"status": statusInt,
					"remark": remark,
				}
				menu.EventKey = defined.SendCommonRequestEvent
			}
		} else if res.Code == defined.UnknowError {
			fmt.Println(res.Mes)
			menu.EventKey = defined.MainPageEvent //回到主界面
		}

	case defined.AddFreindResMesType: //添加好友请求响应信息
		var res defined.FriendApplyResMes
		if err = json.Unmarshal([]byte(response.Data), &res); err != nil {
			return err
		}
		if res.Code == 1 {
			menu.Mes.Content = append(menu.Mes.Content, defined.ReceiveMes{
				MesId:        res.Data.MesId,
				Type:         res.Data.Type,
				FromUserId:   res.Data.FromUserId,
				FromUserName: res.Data.FromUserName,
				ToUserName:   res.Data.ToUserName,
				ToUserId:     res.Data.ToUserId,
				ReadStatus:   res.Data.ReadStatus,
				CreateTime:   res.Data.CreateTime,
			})
		}
		fmt.Println(res.Mes)
		fmt.Println("---回车返回上一步---")
		fmt.Scanln()
		menu.EventKey = defined.MainPageEvent //回到主界面

	case defined.UpdateFriendResMesType: //更新好友申请
		var res defined.Response
		if err = json.Unmarshal([]byte(response.Data), &res); err != nil {
			return err
		}
		if res.Code == 1 {
			fmt.Println(res.Mes)
		} else {
			fmt.Println("发送失败")
		}
		menu.EventKey = defined.MainPageEvent //回到主界面

	case defined.JoinChatRoomResMesType: //进入聊天室
		var res defined.Response
		if err = json.Unmarshal([]byte(response.Data), &res); err != nil {
			return err
		}
		if res.Code == 1 {
			fmt.Println(res.Mes)
			menu.EventKey = defined.ChatRoom
		} else {
			fmt.Println("发送失败")
		}

	default:
	}
	return nil
}

//处理nsq消息
func (this *callback) AsynchronNsqHandle(response defined.Message, menu *service.Menu) {
	var mes defined.ReceiveMes
	err := json.Unmarshal([]byte(response.Data), &mes)
	if err != nil {
		utils.CDD("消息解析失败：" + err.Error())
		return
	}
	switch mes.Type {
	case defined.NsqFriendapply: //好友申请消息
		menu.Mes.UnRead++ //消息未读数+1
		menu.Mes.Content = append(menu.Mes.Content, defined.ReceiveMes{
			MesId:        mes.MesId,
			Type:         mes.Type,
			FromUserId:   mes.FromUserId,
			FromUserName: mes.FromUserName,
			ReadStatus:   0, //新进来的消息，阅读状态默认为0
			CreateTime:   mes.CreateTime,
		})
	case defined.NsqUpdateFriendApply: //更新
		menu.Mes.UnRead++ //消息未读数+1
		for k, v := range menu.Mes.Content {
			if v.MesId == mes.MesId {
				menu.Mes.Content[k].ReadStatus = 0
				menu.Mes.Content[k].CreateTime = mes.CreateTime
			}
		}
	default:
		utils.CDD("未定义的消息类型")
	}
}

//处理聊天室消息
func (this *callback) AsynchronChatRoomHandle(response defined.Message, menu *service.Menu) {
	var res defined.ChatMes
	if err := json.Unmarshal([]byte(response.Data), &res); err != nil {
		utils.CDD(err.Error())
		return
	}
	switch res.Type {
	case defined.ChatText:
		if res.UserId != menu.Cache.UserId { //不是自己发的才打印
			fmt.Println("【", res.UserName, "】：", res.Content)
		}
	case defined.ChatHello:
		fmt.Println("【", res.UserName, "】：", "加入聊天室")
	}
}

//处理点对点聊天消息
func (this *callback) AsynchronPtoPChatHandle(response defined.Message, menu *service.Menu) {
	var res defined.ChatInfo
	if err := json.Unmarshal([]byte(response.Data), &res); err != nil {
		utils.CDD(err.Error())
		return
	}
	chatId := menu.FriendChatId[res.FromUserId]
	if menu.PtoPChatObj == res.FromUserId { //当前正在与这个对象聊天
		//更新聊天消息为已读
		if err := service.ChatService.UpdateInfoStatus(menu, chatId, fmt.Sprintf("%d", res.InfoId)); err != nil {
			utils.CDD("AsynchronPtoPChatHandle" + err.Error())
			return
		}
		//打印聊天消息
		fmt.Println("【", res.FromUserName, "】：", res.Content)
		res.ReadStatus = 1
	} else {
		res.ReadStatus = 0
		//未阅读数+1
		menu.ChatStatus[chatId]++
	}
	//当前聊天信息推入客户端缓存
	menu.InfoList[chatId] = append(menu.InfoList[chatId], res)
	return
}

//更新好友列表状态
func (this *callback) AsynchronUpdateFreindList(response defined.Message, menu *service.Menu) {
	var res defined.UpdateFriendListResMsg
	if err := json.Unmarshal([]byte(response.Data), &res); err != nil {
		utils.CDD(err.Error())
		return
	}
	for k, v := range menu.FriendList {
		if v.UserId == res.UserId {
			menu.FriendList[k].OnlineStatus = true
		}
	}
	return
}
