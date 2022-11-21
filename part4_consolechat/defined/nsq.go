package defined

//nsq消息类型
const (
	NsqFriendapply       = "friendapply"
	NsqUpdateFriendApply = "NsqUpdateFriendApply"
)

//发送消息结构体
type SendMes struct {
	FromUserId string `json:"from_user_id"`
	ToUserId   string `json:"to_user_id"`
	MesId      int    `json:"mes_id"`
	Type       string `json:"type"`
	CreateTime int64  `json:"create_time"`
}

//客户端消息结构
type ReceiveMes struct {
	MesId        int    `json:"mes_id"`
	Type         string `json:"type"`
	FromUserName string `json:"from_user_name"`
	FromUserId   string `json:"from_user_id"`
	ToUserId     string `json:"to_user_id"`
	ToUserName   string `json:"to_user_name"`
	ReadStatus   int    `json:"read_status"`
	CreateTime   int64  `json:"create_time"` //这里的这个当更新时间用
}
