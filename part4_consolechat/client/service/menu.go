package service

import (
	"fmt"
	"chat/defined"
	"chat/library/transfer"
	"chat/utils"
	"os"
	"time"
)

//显示菜单的对象，也是贯穿整个客户端的上下文
var MenuService = new(Menu)

type Menu struct {
	Key          int                           //接收用户的选择
	Loop         bool                          //判断是否还继续显示菜单
	EventKey     defined.EventKey              //用户想要执行的事件id
	Transfer     *transfer.Transfer            //连接对象
	Response     chan defined.Message          //响应信息管道
	Cache        defined.MenuCache             //缓存...userId之类
	Mes          defined.MenuMes               //客户端接收的消息...好友申请之类的
	Request      defined.Request               //存放将要发送给服务器的数据
	FriendList   []defined.FriendList          //好友列表
	FriendChatId map[string]string             //好友与会话id映射，map[userId]chatId
	ChatStatus   map[string]int                //会话未读数, map[chatId]未读数
	PtoPChatObj  string                        //点对点聊天时，当前的聊天对象，如果没有为空窜
	InfoList     map[string][]defined.ChatInfo //存储聊天信息的链表,[chatId]聊天信息数组
}

// @Description 获取输入值赋值到this.key	go官方不推荐接收者命名为this,self,me这类，这里我懒得改了，实际上我这里用的是指针，this就是Menu对象的地址
func (this *Menu) input() {
	_, err := fmt.Scanf("%d\r\n", &this.Key)
	if err != nil {
		fmt.Println("?????")
		this.Loop = true
	}
}

// @Description 首页界面
func (this *Menu) IndexMenu() {
	for this.Loop {
		fmt.Println("------------------棠大聊天室----------------------")
		fmt.Println("                1 登录聊天室")
		fmt.Println("                2 注册用户")
		fmt.Println("                3 退出系统")
		fmt.Println("                请选择（1-3）")

		this.input()
		switch this.Key {
		case 1:
			this.EventKey = defined.LoginEvent
			this.Loop = false
		case 2:
			this.EventKey = defined.RegisterEvent
			this.Loop = false
		case 3:
			this.Transfer.Close()
			fmt.Println("滴滴~成功退出")
			os.Exit(886)
		default:
			fmt.Println("兄弟，输错啦！")
			this.Loop = true
		}
	}
}

//	@Description 主界面
func (this *Menu) MainMenu() {
	for this.Loop {
		fmt.Println("------------------欢迎", this.Cache.UserName, "热烈欢迎------------------")
		fmt.Println("                1 显示在线用户列表")
		fmt.Println("                2 聊天室")
		fmt.Println("                3 好友列表【", this.CalcInfoUnReadNum(), "】")
		fmt.Println("                4 消息" + "【" + fmt.Sprintf("%d", this.Mes.UnRead) + "信息】")
		fmt.Println("                5 添加好友")
		fmt.Println("                6 刷新当前界面信息")
		fmt.Println("                7 退出系统")
		fmt.Println("                请选择（1-7）")

		this.input()
		switch this.Key {
		case 1:
			this.EventKey = defined.AllOnlineUser //显示在线用户
			this.Loop = false
		case 2:
			this.EventKey = defined.JoinChatRoom
			this.Loop = false
		case 3:
			this.EventKey = defined.FriendListEvent
			this.Loop = false
		case 4:
			this.EventKey = defined.ShowAllMesEvent //查看消息列表
			this.Loop = false
		case 5:
			this.EventKey = defined.AddFriendEvent //添加好友
			this.Loop = false
		case 6:
			continue //刷新展示菜单
		case 7:
			this.Transfer.Close()
			fmt.Println("滴滴~成功退出")
			os.Exit(886)
		default:
			fmt.Println("兄弟，输错啦！")
			this.Loop = true
		}
	}
}

//计算当前所有会话未读消息数
func (this *Menu) CalcInfoUnReadNum() string {
	res := 0
	for _, v := range this.ChatStatus {
		res += v
	}
	return fmt.Sprintf("%d", res)
}

//展示消息
func (this *Menu) ShowMes() {
	fmt.Println("                当前共有【", this.Mes.UnRead, "】条新消息！")
	temp := this.bubbleSort(this.Mes.Content) //按创建时间排序
	for k, v := range temp {
		switch v.Type {
		case defined.NsqFriendapply: //好友申请消息
			if v.ReadStatus == 0 {
				fmt.Print("(未读)")
			}
			if v.FromUserId == this.Cache.UserId { //自己发出去的信息
				fmt.Println("【"+fmt.Sprintf("%d", k+1)+"】", "我请求添加"+v.ToUserName+"好友！", "\t\t", time.Unix(v.CreateTime, 0).Format("2006-01-02 15:04:05"))
			} else {
				fmt.Println("【"+fmt.Sprintf("%d", k+1)+"】", v.FromUserName, "请求添加好友！", "\t\t", time.Unix(v.CreateTime, 0).Format("2006-01-02 15:04:05"))
			}
		default:
			utils.CDD("showMes信息展示失败，类型为：" + v.Type)
		}
	}
}

//按CreatTime 冒泡排序 menu.Mes.Content
func (this *Menu) bubbleSort(slice []defined.ReceiveMes) []defined.ReceiveMes {
	for n := 0; n <= len(slice); n++ {
		for i := 1; i < len(slice)-n; i++ {
			if slice[i].CreateTime > slice[i-1].CreateTime {
				slice[i], slice[i-1] = slice[i-1], slice[i]
			}
		}
	}
	return slice
}

//展示好友列表
func (this *Menu) ShowFriendList() {
	for k, v := range this.FriendList {
		online := "离线"
		if v.OnlineStatus {
			online = "在线"
		}
		fmt.Println("【", k+1, "】", v.UserName, "--", online, "【", this.ChatStatus[v.ChatId], "消息】")
	}
}

//展示聊天记录
func (this *Menu) ShowInfo(chatId string) {
	if _, ok := this.InfoList[chatId]; !ok { //不存在
		return
	}
	for k, v := range this.InfoList[chatId] {
		switch v.Type {
		case defined.ChatText: //普通文本
			if v.FromUserId != this.Cache.UserId {
				fmt.Println("【", v.FromUserId, "】：", v.Content)
			} else {
				fmt.Println("【我】：", v.Content)
			}
		}
		//更改阅读状态
		if v.ReadStatus != 1 {
			if err := ChatService.UpdateInfoStatus(this, chatId, fmt.Sprintf("%d", v.InfoId)); err != nil {
				utils.CDD("ShowInfo" + err.Error())
			}
			this.InfoList[chatId][k].ReadStatus = 1
			this.ChatStatus[chatId]--
		}
	}
	return
}
