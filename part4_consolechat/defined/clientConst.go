package defined

//客户端菜单缓存
type MenuCache struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
}

//客户端菜单消息
type MenuMes struct {
	UnRead  int
	Content []ReceiveMes
}

type Request struct {
	Type string
	Data map[string]interface{}
}

//客户端聊天信息
type ChatInfo struct {
	InfoId       int      `json:"info_id"`
	FromUserId   string   `json:"form_user_id"`
	FromUserName string   `json:"from_user_name"`
	ToUserId     string   `json:"to_user_id"`
	ToUserName   string   `json:"to_user_name"`
	Type         ChatType `json:"type"`
	Time         int64    `json:"time"`
	Content      string   `json:"content"`
	ReadStatus   int      `json:"read_status"`
}
