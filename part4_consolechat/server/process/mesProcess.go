package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"chat/dao"
	"chat/defined"
	"chat/library/config"
	"chat/library/logicerror"
	"chat/library/nsq"
	"chat/server/service"
	"chat/utils"
	"strconv"
	"time"
)

type MesProcess struct {
	Processor *Processor
}

//获取好友申请详情
func (s MesProcess) GetFriendApplyDetail(mesId string) defined.GetFriendApplyDetailResMes {
	ctx := *s.Processor.Ctx
	myUserId, _ := ctx.Value("userId").(string)
	friendApplyData, err := dao.NewMesDao().GetFriendApplyDetail(mesId)
	if err != nil {
		utils.SDD("GetFriendApplyDetail 请求失败，错误信息：" + err.Error())
		return defined.GetFriendApplyDetailResMes{
			Code: defined.UnknowError,
			Mes:  defined.ERROR_MES_QUERY.Error(),
			Data: defined.FriendApplyDetailRes{},
		}
	}
	//userId替换为昵称，方便客户端遍历展示
	friendName := ""
	for k, v := range friendApplyData.Remark {
		if v.UserId != myUserId && friendName != "" {
			friendApplyData.Remark[k].UserId = friendName
		} else if friendName == "" && v.UserId != myUserId {
			friendName, _ = dao.NewUserDao().GetUserNameByUserId(v.UserId)
			friendApplyData.Remark[k].UserId = friendName
		} else {
			friendApplyData.Remark[k].UserId = "我"
		}
	}
	if err != nil {
		utils.SDD("GetFriendApplyDetail 请求失败，json转换失败：" + err.Error())
		return defined.GetFriendApplyDetailResMes{
			Code: defined.UnknowError,
			Mes:  defined.ERROR_MES_QUERY.Error(),
			Data: defined.FriendApplyDetailRes{},
		}
	}
	mesIdInt, _ := strconv.Atoi(mesId)
	return defined.GetFriendApplyDetailResMes{
		Code: defined.Success,
		Mes:  defined.SUCCESS_MES_QUERY,
		Data: defined.FriendApplyDetailRes{
			MesId:        mesIdInt,
			FromUserId:   friendApplyData.FromUserId,
			FromUserName: friendName,
			ToUserId:     friendApplyData.ToUserId,
			Status:       friendApplyData.Status,
			Remark:       friendApplyData.Remark,
		},
	}
}

//添加好友申请
func (s MesProcess) AddFriendApply(userName string, remark string) (res defined.FriendApplyResMes) {
	toUserId := service.VerifyService.IsExistUserName(userName)
	if toUserId == "" {	//用户不存在
		return defined.FriendApplyResMes{Code: defined.UnknowError, Mes: defined.ERROR_USER_NOTEXISTS.Error(), Data: defined.ReceiveMes{}}
	}
	ctx := *s.Processor.Ctx
	fromUserId, _ := ctx.Value("userId").(string)
	isFriend, _ := dao.NewMesDao().IsFriend(fromUserId, toUserId)
	if isFriend {	//已经是好友了
		return defined.FriendApplyResMes{Code: defined.UnknowError, Mes: defined.ERROR_REPEAT_FRIEND.Error(), Data: defined.ReceiveMes{}}
	}
	//构造消息表数据(Redis_Mes)
	mesId, err := dao.NewMesDao().GetGobalMesId()
	createTime := time.Now().Unix()

	if logicerror.PrintError(err) != nil {
		return defined.FriendApplyResMes{Code: defined.UnknowError, Mes: defined.ERROR_MES_SEND.Error(), Data: defined.ReceiveMes{}}
	}

	mes := defined.Mes{
		MesId:              mesId,
		Type:               defined.NsqFriendapply,
		FromUserId:         fromUserId,
		ToUserId:           toUserId,
		FromUserReadStatus: 1,
		ToUserReadStatus:   0,
		CreateTime:         createTime,
		UpdateTime:         createTime,
	}
	mesJson, err := json.Marshal(mes)
	if logicerror.PrintError(err) != nil {
		return defined.FriendApplyResMes{Code: defined.UnknowError, Mes: defined.ERROR_MES_SEND.Error(), Data: defined.ReceiveMes{}}
	}

	//构造好友申请表数据(Redis_Friend_Apply)
	var remarkData = defined.FriendApplyRemark{}
	remarkTemp := []byte{}
	if remark != "" { //第一条留言为发送者
		remarkData = defined.FriendApplyRemark{
			UserId:  fromUserId,
			Content: remark,
			Time:    int(createTime),
		}
		remarkArr := append(make([]defined.FriendApplyRemark, 0, 1), remarkData)
		remarkTemp, err = json.Marshal(remarkArr)
	}
	if logicerror.PrintError(err) != nil {
		return defined.FriendApplyResMes{Code: defined.UnknowError, Mes: defined.ERROR_MES_SEND.Error(), Data: defined.ReceiveMes{}}
	}
	friendApplyRedis := defined.FriendApplyRedis{
		FromUserId: fromUserId,
		ToUserId:   toUserId,
		Status:     0,
		Remark:     string(remarkTemp),
	}

	//写入redis
	if err = dao.NewMesDao().AddFriendMes(fmt.Sprintf("%d", mesId), string(mesJson), friendApplyRedis); logicerror.PrintError(err) != nil {
		return defined.FriendApplyResMes{Code: defined.UnknowError, Mes: defined.ERROR_MES_SEND.Error(), Data: defined.ReceiveMes{}}
	}

	//发送nsq
	nsqJson, _ := json.Marshal(defined.SendMes{
		FromUserId: fromUserId,
		ToUserId:   toUserId,
		MesId:      mesId,
		Type:       defined.NsqFriendapply,
		CreateTime: createTime,
	})
	if !nsq.HttpPush(config.GetConfig().Nsq.Nsqd1Http, config.GetConfig().Nsq.MesTopic, string(nsqJson)) {
		return defined.FriendApplyResMes{Code: defined.UnknowError, Mes: defined.ERROR_MES_SEND.Error(), Data: defined.ReceiveMes{}}
	}

	fromuserName, err := dao.UserDao.GetUserNameByUserId(fromUserId)
	if logicerror.PrintError(err) != nil {
		return defined.FriendApplyResMes{Code: defined.UnknowError, Mes: defined.ERROR_MES_SEND.Error(), Data: defined.ReceiveMes{}}
	}
	//把消息响应回给客户端
	return defined.FriendApplyResMes{Code: defined.Success, Mes: defined.SUCCESS_MES_SEND, Data: defined.ReceiveMes{
		MesId:        mesId,
		Type:         defined.NsqFriendapply,
		FromUserName: fromuserName,
		FromUserId:   fromUserId,
		ToUserId:     toUserId,
		ToUserName:   userName,
		ReadStatus:   1,
		CreateTime:   createTime,
	}}
}

//更新好友申请（即是否同意添加好友等操作）
func (s MesProcess) UpdateFriend(mesId, status int, remark string) (err error) {
	//获取当前连接userId
	ctx := *s.Processor.Ctx
	myUserId, _ := ctx.Value("userId").(string)

	//查询好友申请表
	applyInfo, err := dao.NewMesDao().GetFriendApplyDetail(fmt.Sprintf("%d", mesId))
	if err != nil {
		return err
	}
	time := time.Now().Unix()
	if remark != "" {
		applyInfo.Remark = append(applyInfo.Remark, defined.FriendApplyRemark{
			UserId:  myUserId,
			Content: remark,
			Time:    int(time),
		})
	}
	switch status {
	case 0: //拒绝
		if err = dao.NewMesDao().UpdateFriendStatus(mesId, status, time); logicerror.PrintError(err) != nil { //更新状态
			return
		}
		if remark != "" {
			if err = dao.NewMesDao().UpdateFriendRemark(mesId, applyInfo.Remark); logicerror.PrintError(err) != nil { //存在留言更新留言
				return
			}
		}
	case 1: //同意添加好友
		if err = dao.NewMesDao().AddFriend(mesId, myUserId, time); logicerror.PrintError(err) != nil {
			return
		}
	case 2: //继续留言
		if err = dao.NewMesDao().UpdateFriendRemark(mesId, applyInfo.Remark); logicerror.PrintError(err) != nil {
			return
		}
	}

	//发送nsq信息
	sendUserId := applyInfo.ToUserId //消息接收人
	if sendUserId == myUserId {
		sendUserId = applyInfo.FromUserId
	}
	nsqJson, _ := json.Marshal(defined.SendMes{
		ToUserId:   sendUserId,
		MesId:      mesId,
		Type:       defined.NsqUpdateFriendApply,
		CreateTime: time,
	})
	if !nsq.HttpPush(config.GetConfig().Nsq.Nsqd1Http, config.GetConfig().Nsq.MesTopic, string(nsqJson)) {
		return errors.New("更新好友申请发送nsq失败")
	}
	return nil
}
