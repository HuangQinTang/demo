package process

import (
	"encoding/json"
	"chat/dao"
	"chat/defined"
	"chat/library/usermgr"
	"time"
)

type ChatProcess struct {
	Processor *Processor
}

//发送点对点聊天消息
func (c *ChatProcess) SendPtoPMes(mes defined.PtoPChatMes) error {
	//1.聊天信息写入redis
	var info defined.PtoPChat
	info.Time = time.Now().Unix()
	info.Type = mes.Type
	info.Content = mes.Content
	info.ToUserId = mes.ToUserId
	info.ToUserName = mes.ToUserName
	info.FromUserId = mes.FromUserId
	info.FromUserName = mes.FromUserName

	//创建信息id
	infoId, err := dao.NewMesDao().GetChatInfoId(mes.ChatId)
	if err != nil {
		return err
	}
	info.InfoId = infoId
	infoByte, err := json.Marshal(info)
	if err != nil {
		return err
	}
	if err = dao.NewMesDao().CreateChatQueue(mes.ChatId, string(infoByte)); err != nil {
		return err
	}

	//2. 在线推送过去
	tf, _ := usermgr.UserMgr.GetConnByUserId(info.ToUserId)
	if tf == nil {
		return nil
	}
	var res defined.Message
	res.Type = defined.PtoPChatMesResMesType
	res.Data = string(infoByte)
	resByte, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return tf.WritePkg(resByte)
}

//更新信息阅读状态
func (c *ChatProcess) UpdateInfoStatus(chatId, InfoId string) error {
	return dao.NewMesDao().UpdateInfoStatus(chatId, InfoId)
}
