package service

import (
	"chat/defined"
)

var ChatService = &chatService{}

type chatService struct{}

//加入聊天室请求
func (this chatService) JoinChatRoom(menu *Menu, userId string) error {
	return ApiService.Send(defined.JoinChatRoomMesType, map[string]interface{}{"user_id": userId}, menu)
}

//聊天室发送信息
func (this chatService) ChatRoomMes(menu *Menu, userId, userName, content string) error {
	return ApiService.Send(defined.ChatRoomMesType, map[string]interface{}{
		"user_id":   userId,
		"user_name": userName,
		"content":   content,
		"type":      "text",
	}, menu)
}

//点对点聊天
func (this chatService) PtoPMes(menu *Menu, chatMes defined.PtoPChatMes) error {
	return ApiService.Send(defined.PtoPChatMesType, map[string]interface{}{
		"chat_id":        chatMes.ChatId,
		"form_user_id":   chatMes.FromUserId,
		"from_user_name": chatMes.FromUserName,
		"to_user_id":     chatMes.ToUserId,
		"to_user_name":   chatMes.ToUserName,
		"type":           chatMes.Type,
		"content":        chatMes.Content,
	}, menu)
}

//更新阅读信息阅读状态
func (this chatService) UpdateInfoStatus(menu *Menu, chatId, InfoId string) error {
	return ApiService.Send(defined.UpdatePtoPStatusMesType, map[string]interface{}{
		"chat_id": chatId,
		"info_id": InfoId,
	}, menu)
}
