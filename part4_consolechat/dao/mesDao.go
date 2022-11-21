package dao

import (
	"chat/defined"
	"chat/library/redisPool"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"sync"
)

var (
	MesDao     *mesDao
	mesDaoOnce sync.Once
)

type mesDao struct {
	pool *redis.Pool
}

func NewMesDao() *mesDao {
	mesDaoOnce.Do(func() {
		MesDao = &mesDao{
			pool: redisPool.Pool,
		}
	})
	return MesDao
}

//添加 待拉取的消息
func (d mesDao) AddReceiveMes(userId string, mesId int, createTime int64) error {
	conn := d.pool.Get()
	defer conn.Close()

	_, err := conn.Do("zadd", defined.Redis_Receive_Mes+userId, createTime, mesId)
	return err
}

//查询消息详情——好友申请
func (d mesDao) GetFriendApplyDetail(mesId string) (defined.FriendApply, error) {
	conn := d.pool.Get()
	defer conn.Close()

	var friendApply defined.FriendApply
	res, err := redis.StringMap(conn.Do("hgetall", defined.Redis_Friend_Apply+mesId))
	if err != nil {
		return defined.FriendApply{}, err
	}
	status, _ := strconv.Atoi(res["status"])
	friendApply.Status = status
	friendApply.FromUserId = res["from_user_id"]
	friendApply.ToUserId = res["to_user_id"]

	var remark []defined.FriendApplyRemark
	if res["remark"] != "" {
		err = json.Unmarshal([]byte(res["remark"]), &remark)
		if err != nil {
			return defined.FriendApply{}, err
		}
	}
	friendApply.Remark = remark
	return friendApply, nil
}

//判断是否已经添加好友
func (d mesDao) IsFriend(myUserId, friendId string) (bool, error) {
	conn := d.pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("SISMEMBER", defined.Redis_Friend+myUserId, friendId))
}

//多条好友申请详情
func (d mesDao) GetFriendApplyByMesIds(mesIds []string) (res []defined.FriendApply, err error) {
	conn := d.pool.Get()
	defer conn.Close()

	for _, v := range mesIds {
		if err = conn.Send("hgetall", defined.Redis_Friend_Apply+v); err != nil {
			return res, err
		}
	}
	if err = conn.Flush(); err != nil {
		return res, err

	}
	for i := 1; i <= len(mesIds); i++ {
		friendApplyData, err := redis.StringMap(conn.Receive())
		if err != nil {
			return res, err
		}
		status, _ := strconv.Atoi(friendApplyData["status"])
		var remark []defined.FriendApplyRemark
		if friendApplyData["remark"] != "" {
			err = json.Unmarshal([]byte(friendApplyData["remark"]), &remark)
		}
		res = append(res, defined.FriendApply{
			FromUserId: friendApplyData["from_user_id"],
			ToUserId:   friendApplyData["to_user_id"],
			Status:     status,
			Remark:     remark,
		})
	}
	return res, nil
}

//获取全局自增消息id
func (d mesDao) GetGobalMesId() (mesId int, err error) {
	conn := d.pool.Get()
	defer conn.Close()

	mesId, err = redis.Int(conn.Do("incr", defined.Redis_Mes_Id))
	if err != nil {
		return
	}
	return mesId, nil
}

//添加好友请求（会写入消息表，好友申请表）
func (d mesDao) AddFriendMes(mesId, mesStr string, friendApplyRedis defined.FriendApplyRedis) (err error) {
	conn := d.pool.Get()
	defer conn.Close()

	//插入消息表，好友申请表,multi只能批处理redis命令，但无法保证两条命令作为一个原子操作，如果命令存在语法问题，事务提交失败，如果存在运行时错误，忽略错误命令
	if err = conn.Send("multi"); err != nil {
		return err
	}
	if err = conn.Send("set", defined.Redis_Mes+mesId, mesStr); err != nil {
		return err
	}
	if err = conn.Send("hmset", redis.Args{defined.Redis_Friend_Apply + mesId}.AddFlat(friendApplyRedis)...); err != nil {
		return err
	}
	_, err = conn.Do("exec")
	if err != nil {
		return err
	}
	return nil
}

//更新好友申请表留言
func (d mesDao) UpdateFriendRemark(mesId int, remark []defined.FriendApplyRemark) error {
	conn := d.pool.Get()
	defer conn.Close()

	remarkByte, _ := json.Marshal(remark)
	_, err := conn.Do("hset", defined.Redis_Friend_Apply+fmt.Sprintf("%d", mesId), "remark", string(remarkByte))
	return err
}

//更新消息时间
func (d mesDao) UpdateMesTime(mesId int, time int64) error {
	conn := d.pool.Get()
	defer conn.Close()

	mesStr, err := redis.String(conn.Do("get", defined.Redis_Mes+fmt.Sprintf("%d", mesId)))
	var mes defined.Mes
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(mesStr), &mes); err != nil {
		return err
	}
	mes.UpdateTime = time
	newMes, _ := json.Marshal(mes)
	_, err = conn.Do("set", defined.Redis_Mes+fmt.Sprintf("%d", mesId), string(newMes))
	return err
}

//更新好友申请表状态
func (d mesDao) UpdateFriendStatus(mesId int, status int, time int64) (err error) {
	conn := d.pool.Get()
	defer conn.Close()

	if err = d.UpdateMesTime(mesId, time); err != nil {
		return err
	}
	if _, err = conn.Do("hset", defined.Redis_Friend_Apply+fmt.Sprintf("%d", mesId), "status", status); err != nil {
		return err
	}
	return nil
}

//同意添加好友
func (d mesDao) AddFriend(mesId int, userId string, time int64) (err error) {
	conn := d.pool.Get()
	defer conn.Close()

	friendId, err := redis.String(conn.Do("hget", defined.Redis_Friend_Apply+fmt.Sprintf("%d", mesId), "from_user_id"))
	if err != nil {
		return
	}

	//创建会话id，我不知道为什么这一步放在事务里会出问题，感觉是redisgo包的bug
	chatId, err := redis.Int(conn.Do("incr", defined.Redis_Chat_Id))
	if err != nil {
		return
	}
	if err = conn.Send("multi"); err != nil {
		return err
	}
	if err = d.UpdateMesTime(mesId, time); err != nil {
		conn.Do("DISCARD") //取消事务
		return err
	}
	//更改好友申请表状态
	if _, err = conn.Do("hset", defined.Redis_Friend_Apply+fmt.Sprintf("%d", mesId), "status", 1); err != nil {
		conn.Do("DISCARD") //取消事务
		return err
	}
	//添加好友到集合
	if _, err = conn.Do("sadd", defined.Redis_Friend+userId, friendId); err != nil {
		conn.Do("DISCARD")
		return err
	}
	if _, err = conn.Do("sadd", defined.Redis_Friend+friendId, userId); err != nil {
		conn.Do("DISCARD")
		return err
	}
	//添加会话映射
	if _, err = conn.Do("hset", defined.Redis_Friend_Chat+userId, friendId, chatId); err != nil {
		fmt.Println(3)
		fmt.Println(err.Error())
		conn.Do("DISCARD")
		return err
	}
	if _, err = conn.Do("hset", defined.Redis_Friend_Chat+friendId, userId, chatId); err != nil {
		fmt.Println(4)
		fmt.Println(err.Error())
		conn.Do("DISCARD")
		return err
	}
	if _, err = conn.Do("exec"); err != nil {
		conn.Do("DISCARD")
		return err
	}
	return nil
}

//查询最近20条消息的消息id
func (d mesDao) GetReceiveMes(userId string) ([]string, error) {
	conn := d.pool.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("ZREVRANGEBYSCORE", defined.Redis_Receive_Mes+userId, "+inf", "-inf", "limit", 0, 20))
}

//根据多个消息id换取多个消息
func (d mesDao) GetMesByMesIds(mesIds []string) ([]defined.Mes, error) {
	conn := d.pool.Get()
	defer conn.Close()

	for _, mesId := range mesIds {
		if err := conn.Send("get", defined.Redis_Mes+mesId); err != nil {
			return []defined.Mes{}, err
		}
	}
	if err := conn.Flush(); err != nil {
		return []defined.Mes{}, err
	}
	var res []defined.Mes
	for i := 0; i < len(mesIds); i++ {
		item, err := redis.String(conn.Receive())
		if err != nil {
			return []defined.Mes{}, err
		}
		var mes defined.Mes
		if err = json.Unmarshal([]byte(item), &mes); err != nil {
			return []defined.Mes{}, err
		}
		res = append(res, mes)
	}
	return res, nil
}

//创建会话下的信息id
func (d mesDao) GetChatInfoId(chatId string) (int, error) {
	conn := d.pool.Get()
	defer conn.Close()

	infoId, err := redis.Int(conn.Do("incr", defined.Redis_FriendChat_InfoId+chatId))
	if err != nil {
		return 0, err
	}
	return infoId, nil
}

//创建会话队列
func (d mesDao) CreateChatQueue(chatId string, info string) error {
	conn := d.pool.Get()
	defer conn.Close()

	_, err := conn.Do("rpush", defined.Redis_Chat_Queue+chatId, info)
	return err
}

//更新信息状态
func (d mesDao) UpdateInfoStatus(chatId, InfoId string) error {
	conn := d.pool.Get()
	defer conn.Close()

	_, err := conn.Do("hset", defined.Redis_Info_Status+chatId, InfoId, 1)
	return err
}

//查询会话下最新50条消息
func (this mesDao) GetLastInfo(chatId string) ([]defined.ChatInfo, error) {
	conn := this.pool.Get()
	defer conn.Close()

	//查询最近50条信息
	res, err := redis.Strings(conn.Do("lrange", defined.Redis_Chat_Queue+chatId, 0, 49))
	if err != nil {
		return []defined.ChatInfo{}, err
	}
	if len(res) == 0 {
		return []defined.ChatInfo{}, err
	}
	info := make([]defined.ChatInfo, 0, len(res))
	infoIds := redis.Args{}.Add(defined.Redis_Info_Status + chatId)
	for _, v := range res {
		var temp defined.ChatInfo
		if err = json.Unmarshal([]byte(v), &temp); err != nil {
			return []defined.ChatInfo{}, err
		}
		info = append(info, temp)
		infoIds = infoIds.Add(temp.InfoId) //推入聊天信息id，用于批量查询信息是否阅读
	}
	//查询聊天信息id是否已读
	readStatus, err := redis.Strings(conn.Do("hmget", infoIds...))
	if err != nil {
		return []defined.ChatInfo{}, err
	}
	//瓶装信息阅读状态
	for k, v := range readStatus {
		status, _ := strconv.Atoi(v)
		info[k].ReadStatus = status
	}
	return info, nil
}
