package defined

import "errors"

const (
	SUCCESS_MES_QUERY = "查询成功！"
	SUCCESS_MES_SEND  = "发送成功！"
)

//错误变量
var (
	ERROR_MES_QUERY = errors.New("查询失败，请稍候再试~")
	ERROR_MES_SEND  = errors.New("发送失败")
)

//查看好友申请详情请求
type GetFriendApplyDetailMes struct {
	MesId string `json:"mes_id"`
}

//查看好友申请详情响应
type GetFriendApplyDetailResMes struct {
	Code HttpCode             `json:"code"`
	Mes  string               `json:"mes"`
	Data FriendApplyDetailRes `json:"data"`
}

type FriendApplyDetailRes struct {
	MesId        int                 `json:"mes_id"`
	FromUserId   string              `json:"from_user_id"`
	FromUserName string              `json:"from_user_name"`
	ToUserId     string              `json:"to_user_id"`
	Status       int                 `json:"status"`
	Remark       []FriendApplyRemark `json:"remark"`
}

//添加好友申请消息写入redis结构体
type FriendApplyRedis struct {
	FromUserId string `redis:"from_user_id"`
	ToUserId   string `redis:"to_user_id"`
	Status     int    `redis:"status"`
	Remark     string `redis:"remark"`
}

//添加好友请求
type AddFriendMes struct {
	UserName string `json:"user_name"`
	Remark   string `json:"remark"`
}

//添加好友响应结构体
type FriendApplyResMes struct {
	Code HttpCode   `json:"code"`
	Mes  string     `json:"mes"`
	Data ReceiveMes `json:"data"`
}

type UpdateFriendMes struct {
	MesId  int    `json:"mes_id"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}
