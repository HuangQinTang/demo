package service

import (
	"chat/dao"
	"chat/defined"
	"chat/library/config"
	mynsq "chat/library/nsq"
	"chat/library/usermgr"
	"chat/utils"
	"encoding/json"
	"errors"
	"github.com/nsqio/go-nsq"
	"os"
	"os/signal"
	"syscall"
)

var NsqService = nsqServerice{}

type nsqServerice struct{}

func (s nsqServerice) InitMesConsumerNsq() {
		conf := mynsq.NsqConsumerConfig{
		Topic:   config.GetConfig().Nsq.MesTopic,
		Channel: config.GetConfig().Nsq.MesChannel,
	}

	consumer, err := conf.NewConsumer(nsq.HandlerFunc(s.mesConsumerHandle))
	if err != nil {
		utils.SDD("InitFriendApplyConsumerMq 创建消费者失败" + err.Error())
		return
	}

	//集群连接
	err = consumer.ConnLookupd()
	if err != nil {
		utils.SDD("InitFriendApplyConsumerMq 连接nsqlookupd失败" + err.Error())
		return
	}

	//等待系统信号
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt
	consumer.Stop() //停止连接
}

//消费者回调处理函数，消息内容会传入message.Body，返回结果为error消息将会重新入队列，返回结果Nil则标记为消费成功
func (s nsqServerice) mesConsumerHandle(message *nsq.Message) (err error) {
	var mes defined.SendMes
	err = json.Unmarshal(message.Body, &mes)
	if err != nil {
		utils.SDD("mesConsumerHandle json解析失败失败" + err.Error())
		return err
	}

	switch mes.Type {
	//好友申请消息，可以不引入nsq，算了，当作练习nsq
	case defined.NsqFriendapply:	//添加好友申请
		transfer, _ := usermgr.UserMgr.GetConnByUserId(mes.ToUserId)
		if transfer == nil { //用户不在线消息放入redis等待用户主动拉取
			return dao.NewMesDao().AddReceiveMes(mes.ToUserId, mes.MesId, mes.CreateTime)
		} else { //用户在线，发送给用户
			userName, err := dao.NewUserDao().GetUserNameByUserId(mes.FromUserId)
			if err != nil {
				utils.SDD("mesConsumerHandle defined.NsqFriendapply类型 redis查询失败" + err.Error())
			}
			toUserName, err := dao.NewUserDao().GetUserNameByUserId(mes.ToUserId)
			if err != nil {
				utils.SDD("mesConsumerHandle defined.NsqFriendapply类型 redis查询失败" + err.Error())
			}
			clientReceiveMes := defined.ReceiveMes{
				MesId: mes.MesId,
				Type: mes.Type,
				FromUserId: mes.FromUserId,
				FromUserName: userName,
				ToUserId: mes.ToUserId,
				ToUserName: toUserName,
				ReadStatus: 0,
				CreateTime: mes.CreateTime,
			}
			clientReceiveMesJson, err := json.Marshal(clientReceiveMes)
			if err != nil {
				utils.SDD("mesConsumerHandle json解析失败失败" + err.Error())
				return err
			}
			var resMes defined.Message
			resMes.Type = defined.NsqMesType
			resMes.Data = string(clientReceiveMesJson)
			resMesByet, _ := json.Marshal(resMes)
			err = transfer.WritePkg(resMesByet)		//如果用户突然离线了，发送失败会重新投递，积累的白白跑生产者cpu?
			return err
		}

	case defined.NsqUpdateFriendApply:	//更新好友申请消息
		transfer, _ := usermgr.UserMgr.GetConnByUserId(mes.ToUserId)
		if transfer == nil {	//接收方不在线,不作任何操作
			return err
		}
		//在线推过去
		var resMes defined.Message
		resMes.Type = defined.NsqMesType
		clientReceiveMes := defined.ReceiveMes{
			MesId: mes.MesId,
			Type: mes.Type,
			CreateTime: mes.CreateTime,
		}
		clientReceiveMesJson, err := json.Marshal(clientReceiveMes)
		if err != nil {
			utils.SDD("mesConsumerHandle json解析失败失败" + err.Error())
			return err
		}
		resMes.Data = string(clientReceiveMesJson)
		resMesByet, _ := json.Marshal(resMes)
		return transfer.WritePkg(resMesByet)

	default:
		utils.SDD("mesConsumerHandle 消息类型未定义" + mes.Type)
		return errors.New("消息类型未定义")
	}
}
