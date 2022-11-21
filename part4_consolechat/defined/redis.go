package defined

const (
	Redis_User_Info         = "user:userid:" //用户表, value为用户信息,	hash
	Redis_UserName_Prefix   = "username:"    //冗余，方便用昵称换取id, value为用户id,	string
	Redis_UserName_Postfix  = ":userid"
	Redis_Online_User       = "online:username"               //当前在线用户, value为用户昵称,	sets
	Redis_Mes_Id            = "gobal:mesid"                   //全局自增消息id
	Redis_Mes               = "mes:mesid:"                    //消息 string 存json type,create_time,from_user_id,to_user_id,from_user_read,to_user_read
	Redis_Friend_Apply      = "friendapply:mesid:"            //好友申请表 hash from_user_id(申请人),to_user_id(被申请人),remark(附言json,userid用户id，content内容，time创建时间),status(申请状态)
	Redis_Receive_Mes       = "receive:mes:userid:"           //用户需要拉取的消息 sort set 存放消息id
	Redis_Friend            = "friend:userid:"                //好友表，set集合，存放好友id
	Redis_Chat_Id           = "gobal:chatid"                  //全局自增会话id
	Redis_Friend_Chat       = "friend:chat:userid:"           //好友会话映射表，hash，每个好友分配一个指定的会话id
	Redis_Chat_Queue        = "chat:queue:chatid:"            //聊天会话，list，存储聊天信息
	Redis_FriendChat_InfoId = "friend:chatInfoId:chatid:"     //当前会话下聊天记录id，string，自增
	Redis_Info_Status       = "friend:chatInfoStatus:chatid:" //会话下，信息的阅读状态，hash key为信息id，value值为1表示会话下这条消息已读
)

//redis消息结构(Redis_Mes)
type Mes struct {
	MesId              int    `json:"mes_id"`
	Type               string `json:"type"`
	CreateTime         int64  `json:"create_time"`
	UpdateTime         int64  `json:"update_time"`
	FromUserId         string `json:"from_user_id"`
	ToUserId           string `json:"to_user_id"`
	FromUserReadStatus int    `json:"from_user_read_status"` //阅读状态，0未读，1已读。
	ToUserReadStatus   int    `json:"to_user_read_status"`
}

//redis好友申请表结构体(Redis_Friend_Apply)
type FriendApply struct {
	FromUserId string              `json:"from_user_id" redis:"from_user_id"`
	ToUserId   string              `json:"to_user_id" redis:"to_user_id"`
	Status     int                 `json:"status" redis:"status"`
	Remark     []FriendApplyRemark `json:"remark" redis:"remark"`
}

type FriendApplyRemark struct {
	UserId  string `json:"user_id" redis:"user_id"`
	Content string `json:"content" redis:"content"`
	Time    int    `json:"time" redis:"time"`
}

//聊天信息结构体
type PtoPChat struct {
	InfoId       int      `json:"info_id"`
	FromUserId   string   `json:"form_user_id"`
	FromUserName string   `json:"from_user_name"`
	ToUserId     string   `json:"to_user_id"`
	ToUserName   string   `json:"to_user_name"`
	Type         ChatType `json:"type"`
	Time         int64    `json:"time"`
	Content      string   `json:"content"`
}


