package defined

type EventKey int

const (
	HomePageEvent          EventKey = iota + 1 //回到首页事件
	LoginEvent                                 //登录事件
	LoginSuccessEvent                          //登录成功后事件
	RegisterEvent                              //注册事件
	AllOnlineUser                              //查看当前在线用户
	MainPageEvent                              //主界面菜单
	AddFriendEvent                             //添加好友
	ShowAllMesEvent                            //查看当前所有消息
	SendCommonRequestEvent                     //发送Menu.Request存放的数据
	JoinChatRoom                               //进入聊天室
	ChatRoom                                   //聊天室
	FriendListEvent                            //好友列表
	PtoPChatEvent                              //点对点聊天
)
