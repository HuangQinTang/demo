package event

import (
	"bufio"
	"fmt"
	"chat/boot/clientBoot"
	"chat/client/callback"
	"chat/client/service"
	"chat/defined"
	"chat/utils"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//	@Description 事件处理中心
var BaseEvent = baseEvent{}

type baseEvent struct{}

func (this *baseEvent) ProcessHandle(menu *service.Menu) {
RESELECTEVENT:
	switch menu.EventKey {
	case defined.HomePageEvent: //首页菜单
		menu.Loop = true
		menu.IndexMenu()
		goto RESELECTEVENT //菜单操作完后(会把后续后续执行的事件记录在menu.Eventkey)重新选择事件

	case defined.MainPageEvent: //主菜单
		menu.Loop = true
		menu.MainMenu()
		goto RESELECTEVENT

	case defined.LoginEvent: //登录事件
		var userId string
		var userPwd string
		utils.DumpSAndScanVar("请输入用户id, q返回首页", &userId)
		if userId == "q" {
			menu.EventKey = defined.HomePageEvent
			goto RESELECTEVENT
		}
		utils.DumpSAndScanVar("请输入用户密码, q返回首页", &userPwd)
		if userPwd == "q" {
			menu.EventKey = defined.HomePageEvent
			goto RESELECTEVENT
		}
		msg, err := service.UserService.Login(userId, userPwd, menu)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("登录失败，请联系管理员~")
			menu.EventKey = defined.HomePageEvent
			goto RESELECTEVENT
		}
		clientBoot.MenuJob <- msg //把消息往任务管道发送

	case defined.LoginSuccessEvent: //登录成功后事件
		go this.SendBeat(menu) //心跳
		go this.Receive(menu)  //接收服务端响应
		menu.MainMenu()
		goto RESELECTEVENT //菜单操作完后，重新选择事件

	case defined.RegisterEvent: //注册事件
		var userId string
		var userPwd string
		var userName string
		utils.DumpSAndScanVar("请输入用户id", &userId)
		utils.DumpSAndScanVar("请输入用户密码", &userPwd)
		utils.DumpSAndScanVar("请输入用户昵称（将作为聊天标识）", &userName)
		msg, err := service.UserService.Register(userId, userPwd, userName, menu)
		if err != nil {
			fmt.Println("注册失败，请联系管理员~")
			menu.EventKey = defined.HomePageEvent
			goto RESELECTEVENT
		}
		clientBoot.MenuJob <- msg

	case defined.AllOnlineUser: //查看当前所有在线用户
		err := service.UserService.GetAllOnlineUser(menu)
		if err != nil {
			fmt.Println("查看失败，请联系管理员~")
			menu.MainMenu()
			goto RESELECTEVENT //菜单操作完后，重新选择事件
		}

	case defined.ShowAllMesEvent: //查看当前消息
		menu.ShowMes() //展示信息
		var value string
		utils.DumpSAndScanVar("请输入编号查看详情, (q)退出返回上一步~", &value)
		if value == "q" {
			menu.EventKey = defined.MainPageEvent //回到主界面
			goto RESELECTEVENT                    //重新判断所选事件
		}
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("老铁，请输入整数哟！")
			goto RESELECTEVENT
		}
		if valueInt > len(menu.Mes.Content) || valueInt <= 0 {
			fmt.Println("老铁，请输入对应编号哟！")
			goto RESELECTEVENT
		}

		switch menu.Mes.Content[valueInt-1].Type {
		case defined.NsqFriendapply: //好友申请消息
			//请求消息详情
			err = service.MesService.GetFriendapplyDetail(menu, menu.Mes.Content[valueInt-1].MesId)
			if err != nil {
				fmt.Println("我丢，查看失败！")
				utils.CDD("查看消息详情失败" + err.Error())
				menu.EventKey = defined.MainPageEvent //回到主界面
				goto RESELECTEVENT                    //重新判断所选事件
			}
		default:
		}

	case defined.AddFriendEvent: //添加好友
		var userName string
		utils.DumpSAndScanVar("请输入用户昵称(q返回上一步)", &userName)
		if userName == "q" {
			menu.Loop = true
			menu.MainMenu()
			goto RESELECTEVENT
		}
		var remark string
		utils.DumpSAndScanVar("请输入留言(q返回上一步)", &remark)
		if userName == "q" {
			menu.Loop = true
			menu.MainMenu()
			goto RESELECTEVENT
		}
		if err := service.MesService.SendFriendApplyMes(menu, userName, remark); err != nil {
			utils.CDD("defined.AddFriendEvent 发送好友申请失败" + err.Error())
			menu.EventKey = defined.MainPageEvent //回到主界面
			goto RESELECTEVENT                    //重新判断所选事件
		}

	case defined.SendCommonRequestEvent: //发送公共请求到服务器
		if err := this.CommonSend(menu); err != nil {
			fmt.Println("操作失败！")
			utils.CDD("defined.AddFrieSendCommonRequestndEvent 错误" + err.Error())
		}
		menu.Request.Type = "" //发送完请求把类型清空表示当前没有需要发送的请求

	case defined.JoinChatRoom: //请求进入聊天室
		if err := service.ChatService.JoinChatRoom(menu, menu.Cache.UserId); err != nil {
			utils.CDD(err.Error())
			fmt.Println("进入聊天室失败")
			menu.EventKey = defined.MainPageEvent //回到主界面
			goto RESELECTEVENT                    //重新判断所选事件
		}

	case defined.ChatRoom: //聊天室
		fmt.Println("---欢迎进入聊天室，quit退出聊天室---")
		for {
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadBytes('\n')
			if err != nil || err == io.EOF {
				fmt.Println("发送失败")
				continue
			}
			content := strings.Replace(string(input), "\n", "", -1)
			content = strings.Replace(content, "\r", "", -1)
			if content == "quit" {
				break
			}
			if err = service.ChatService.ChatRoomMes(menu, menu.Cache.UserId, menu.Cache.UserName, content); err != nil {
				fmt.Println("发送失败")
				continue
			}
		}
		menu.EventKey = defined.MainPageEvent //回到主界面
		goto RESELECTEVENT                    //重新判断所选事件

	case defined.FriendListEvent: //查看好友列表
		menu.ShowFriendList() //展示好友列表
	GetFriendListResMesType:
		var value string
		utils.DumpSAndScanVar("请输入编号即可发送消息, (q)退出返回上一步~", &value)
		if value == "q" {
			menu.EventKey = defined.MainPageEvent
			goto RESELECTEVENT //重新判断所选事件
		}
		vInt, _ := strconv.Atoi(value)
		if vInt-1 < 0 || vInt-1 >= len(menu.FriendList) {
			fmt.Println("输入有误~")
			goto GetFriendListResMesType
		}
		menu.ShowInfo(menu.FriendList[vInt-1].ChatId) //展示聊天记录
		var chatMes defined.PtoPChatMes
		friendId := menu.FriendList[vInt-1].UserId //好友的id
		chatMes.ToUserId = friendId
		chatMes.ToUserName = menu.FriendList[vInt-1].UserName
		chatMes.Type = defined.ChatText
		chatMes.FromUserName = menu.Cache.UserName
		chatMes.FromUserId = menu.Cache.UserId
		chatMes.ChatId = menu.FriendChatId[friendId]

		//发送点对点聊天信息
		fmt.Println("---quit退出聊天---")
		menu.PtoPChatObj = chatMes.ToUserId
		for {
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadBytes('\n')
			if err != nil || err == io.EOF {
				fmt.Println("发送失败")
				continue
			}
			content := strings.Replace(string(input), "\n", "", -1)
			content = strings.Replace(content, "\r", "", -1)
			if content == "quit" {
				menu.PtoPChatObj = ""
				break
			}
			chatMes.Content = content
			if err = service.ChatService.PtoPMes(menu, chatMes); err != nil {
				fmt.Println("发送失败")
				continue
			}
		}
		menu.EventKey = defined.MainPageEvent //回到主界面
		goto RESELECTEVENT                    //重新判断所选事件

	default:
		log.Panicln("输入错误！")
	}
}

//发送心跳，保持连接
func (this *baseEvent) SendBeat(menu *service.Menu) {
	timeTicker := time.NewTicker(30 * time.Second)
	for {
		<-timeTicker.C
		if err := service.ApiService.Send(defined.ClientBeatMesType, map[string]interface{}{}, menu); err != nil {
			utils.CDD("心跳发送失败:" + err.Error())
			return
		}
	}
}

//接收服务端响应
func (this *baseEvent) Receive(menu *service.Menu) {
	for {
		mes, err := menu.Transfer.ReadPkg()
		if menu.Transfer.CheckError(err) { //如果连接断开,程序结束
			os.Exit(500)
		}
		switch mes.Type {
		case defined.ServerBeatMesType: //服务端心跳，不作处理
		case defined.NsqMesType: //处理nsq消息
			callback.CallbackEven.AsynchronNsqHandle(mes, menu)
		case defined.ChatRoomMesResMesType: //聊天室消息
			callback.CallbackEven.AsynchronChatRoomHandle(mes, menu)
		case defined.PtoPChatMesResMesType: //点对点聊天消息
			callback.CallbackEven.AsynchronPtoPChatHandle(mes, menu)
		case defined.FriendOnlineMesType: //好友上线通知
			callback.CallbackEven.AsynchronUpdateFreindList(mes, menu)
		default: //默认是菜单操作响应消息，推入菜单响应管道
			clientBoot.MenuJob <- mes
		}
	}
}

//发送公共请求数据到服务端
func (this *baseEvent) CommonSend(menu *service.Menu) error {
	if menu.Request.Type != "" {
		return service.ApiService.Send(menu.Request.Type, menu.Request.Data, menu)
	}
	return nil
}
