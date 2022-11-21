package defined

type JoinChatRoomMes struct {
	UserId string `json:"user_id"`
}

type ChatType string

const (
	ChatText  ChatType = "text"  //普通文本
	ChatHello ChatType = "hello" //进入聊天室时的招呼语
)

//聊天室聊天结构
type ChatMes struct {
	UserId   string   `json:"user_id"`
	UserName string   `json:"user_name"`
	Content  string   `json:"content"`
	Type     ChatType `json:"type"`
	Time     int64    `json:"time"`
}

//点对点聊天请求
type PtoPChatMes struct {
	ChatId       string   `json:"chat_id"`
	InfoId       int      `json:"info_id"`
	FromUserId   string   `json:"form_user_id"`
	FromUserName string   `json:"from_user_name"`
	ToUserId     string   `json:"to_user_id"`
	ToUserName   string   `json:"to_user_name"`
	Type         ChatType `json:"type"`
	Content      string   `json:"content"`
}

//更新信息状态请求
type UpdateInfoStatus struct {
	ChatId string `json:"chat_id"`
	InfoId string `json:"info_id"`
}


