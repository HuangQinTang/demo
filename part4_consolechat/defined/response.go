package defined

type (
	//状态码
	HttpCode uint

	//消息结构体
	Message struct {
		Type string `json:"type"` //消息类型
		Data string `json:"data"` //消息内容
	}
)

//状态码值
const (
	ParamsError HttpCode = iota + 10000
	ValueNotExist
	UnknowError HttpCode = 0
	Success     HttpCode = 1
)

//通用响应
type Response = struct {
	Code HttpCode               `json:"code"`
	Data map[string]interface{} `json:"data"`
	Mes  string                 `json:"mes"`
}

//通用成功响应
var OkResponse = Response{
	Code: Success,
	Data: map[string]interface{}{},
	Mes:  "操作成功",
}

//通用失败响应
var FailResponse = Response{
	Code: UnknowError,
	Data: map[string]interface{}{},
	Mes:  "操作失败",
}

const (
	ClientBeatMesType        = "ping"           //客户端心跳
	ServerBeatMesType        = "pong"           //服务端心跳
	MenuMesType       string = "MenuResMesType" //菜单操作

	//请求类型
	LoginMesType    string = "LoginMes"        //登录
	LoginResMesType string = "LoginResMesType" //登录响应

	RegisterMesType    string = "RegisterMesType"    //注册
	RegisterResMesType string = "RegisterResMesType" //注册响应

	GetAllOnlineUserNameMesType    string = "GetAllOnlineUserNameMesType" //获取所有在线用户昵称
	GetAllOnlineUserNameResMesType string = "GetAllOnlineUserNameResMesType"

	NsqMesType                     = "NsqMesType"              //nsq消息
	UpdatePtoPStatusMesType string = "UpdatePtoPStatusMesType" //更新消息阅读状态
	FriendOnlineMesType            = "FriendOnlineMesType"     //好友上线通知

	GetFriendapplyDetailMesType    string = "GetFriendapplyDetailMesType"    //获取好友申请消息详情
	GetFriendapplyDetailResMesType string = "GetFriendapplyDetailResMesType" //获取好友消息详情响应

	AddFriendMesType    string = "AddFriendMesType"    //添加好友请求
	AddFreindResMesType string = "AddFreindResMesType" //添加好友请求响应

	UpdateFriendMesType    string = "UpdateFriendMesType"    //更新添加好友消息请求
	UpdateFriendResMesType string = "UpdateFriendResMesType" //更新添加好友消息响应

	JoinChatRoomMesType    string = "JoinChatRoomMesType"    //进入聊天室请求
	JoinChatRoomResMesType string = "JoinChatRoomResMesType" //进入聊天室响应

	ChatRoomMesType       string = "ChatRoomMesType"       //发送聊天室消息
	ChatRoomMesResMesType string = "ChatRoomMesResMesType" //聊天室消息响应

	GetFriendListMesType    string = "GetFriendListMesType"    //查询好友列表请求
	GetFriendListResMesType string = "GetFriendListResMesType" //查询好友列表响应

	PtoPChatMesType       string = "PtoPChatMesType"       //点对点聊天
	PtoPChatMesResMesType string = "PtoPChatMesResMesType" //点对点聊天响应
)
