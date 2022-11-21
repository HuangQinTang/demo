package process

import (
	"encoding/json"
	"fmt"
	"chat/dao"
	"chat/defined"
	"chat/library/transfer"
	"chat/utils"
	"sync"
	"time"
)

var ChatRoom = chatRoom{
	Conns: make(map[string]*transfer.Transfer, 20),
}

var lock sync.RWMutex

type chatRoom struct {
	Conns map[string]*transfer.Transfer //聊天室所有用户都在这哟，map[userId]连接
}

//进入聊天室哟
func (p *chatRoom) Join(userId string, conn *transfer.Transfer) {
	lock.Lock()
	p.Conns[userId] = conn
	fmt.Println("加入聊天数组喽")
	lock.Unlock()
}

//进入聊天室时的提示，xxx进入了聊天室
func (p *chatRoom) Hello(userId string) {
	var mes defined.ChatMes
	mes.Type = defined.ChatHello
	userName, _ := dao.NewUserDao().GetUserNameByUserId(userId)
	mes.UserName = userName
	p.SendChat(mes)
}

//聊天室聊天消息
func (p *chatRoom) SendChat(mes defined.ChatMes) {
	if len(p.Conns) <= 0 {
		return
	}
	mes.Time = time.Now().Unix()
	var res defined.Message
	res.Type = defined.ChatRoomMesResMesType
	mesByte, err := json.Marshal(mes)
	if err != nil {
		utils.SDD("SendChat:" + err.Error())
	}
	res.Data = string(mesByte)
	resByte, err := json.Marshal(res)
	if err != nil {
		utils.SDD("SendChat" + err.Error())
	}
	for _, v := range p.Conns {
		if err = v.WritePkg(resByte); err != nil {
			utils.SDD(err.Error())
		}
	}
}
