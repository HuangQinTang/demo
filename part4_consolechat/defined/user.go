package defined

import (
	"errors"
)

//成功响应
const (
	SUCCESS_LOGIN    = "登录成功"
	SUCCESS_REGISTER = "欢迎欢迎，注册成功！"
	SUCCESS_QUERY    = "查询成功"
)

//错误变量
var (
	ERROR_USER_NOTEXISTS  = errors.New("用户不存在")
	ERROR_USERID_EXISTS   = errors.New("用户id已经存在")
	ERROR_USERNAME_EXISTS = errors.New("用户昵称已经存在")
	ERROR_USER_PWD        = errors.New("密码不正确")
	ERROR_REGISTER        = errors.New("注册失败，请稍后再试")
	ERROR_REPEAT_FRIEND   = errors.New("你们已经是好友了")
	ERROR_USER_QUERY_FAIL = errors.New("查询失败")
)

//登录请求
type LoginMes struct {
	UserId   string `json:"user_id"`   //用户id
	UserPwd  string `json:"user_pwd"`  //用户密码
	UserName string `json:"user_name"` //用户名
}

//登录响应(LoginResMesType)
type LoginResMsg struct {
	Code HttpCode     `json:"code"` //返回状态码 1成功 0失败
	Mes  string       `json:"mes"`  //错误信息
	Data LoginResData `json:"data"` //响应数据
}

type LoginResData struct {
	UserName   string                `json:"user_name"`
	UserId     string                `json:"user_id"`
	ReceiveMes []ReceiveMes          `json:"receive_mes"` //带接收消息
	FriendList []FriendList          `json:"friend_list"` //好友列表
	Info       map[string][]ChatInfo `json:"info"`        //聊天消息，map[会话]聊天消息数组
}

//注册请求
type RegisterMes struct {
	UserId   string `json:"user_id"`   //用户id
	UserPwd  string `json:"user_pwd"`  //用户密码
	UserName string `json:"user_name"` //昵称，默认同userid
}

//注册响应(RegisterResMesType)
type RegisterResMsg struct {
	Code HttpCode    `json:"code"` //1成功 0失败
	Mes  string      `json:"mes"`
	Data RegisterMes `json:"data"`
}

//获取在线用户昵称请求
type GetAllOnlineUserNameMes struct {
	UserId string `json:"user_id"`
}

//获取在线用户昵称响应(GetAllOnlineUserNameResMesType)
type GetAllOnlineUserNameResMsg struct {
	Code HttpCode `json:"code"`
	Mes  string   `json:"mes"`
	Data []string `json:"data"`
}

//查询好友列表请求
type GetFriendListMes struct {
	UserId string `json:"user_id"`
}

//好友列表
type FriendList struct {
	UserId       string `json:"user_id"`
	UserName     string `json:"user_name"`
	OnlineStatus bool   `json:"online_status"` //是否在线
	ChatId       string `json:"chat_id"`       //会话id
}

//更新好友列表状态
type UpdateFriendListResMsg struct {
	UserId string `json:"user_id"`
}
