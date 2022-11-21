package dao

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"chat/defined"
	"chat/library/redisPool"
	"chat/model"
	"chat/utils"
	"sync"
)

var (
	UserDao     *userDao
	userDaoOnce sync.Once
)

type userDao struct {
	pool *redis.Pool
}

func NewUserDao() *userDao {
	userDaoOnce.Do(func() {
		UserDao = &userDao{
			pool: redisPool.Pool,
		}
	})
	return UserDao
}

//根据id获取一条用户记录
func (this userDao) GetUserDetailById(userId string) (user model.User, err error) {
	//从连接池获取一根连接
	conn := this.pool.Get()
	//用完丢回连接池
	defer conn.Close()

	res, err := redis.StringMap(conn.Do("hgetall", defined.Redis_User_Info+userId))
	if err != nil {
		return model.User{}, err
	}
	user.UserId = res["user_id"]
	user.UserPwd = res["user_pwd"]
	user.UserName = res["user_name"]
	return user, nil
}

//创建一个用户
func (this userDao) CreateUser(userId, userPwd, userName string) bool {
	conn := this.pool.Get()
	defer conn.Close()

	//lua脚本，以下两个操作作一个原子操作
	//hmset user:userid:root user_id root user_pwd root user_name root
	//set username:root:userid root
	_, err := conn.Do("eval", `
redis.call('hmset', KEYS[1], KEYS[2], KEYS[3], KEYS[4], KEYS[5], KEYS[6], KEYS[7]);
redis.call('set', KEYS[8], KEYS[9]);
return "ok";
`, 9, defined.Redis_User_Info+userId,
		"user_id", userId,
		"user_pwd", userPwd,
		"user_name", userName,
		defined.Redis_UserName_Prefix+userName+defined.Redis_UserName_Postfix, userId)

	if err != nil {
		conn.Do("del", defined.Redis_User_Info+userId)
		conn.Do("del", defined.Redis_UserName_Prefix+userName+defined.Redis_UserName_Postfix)
		utils.SDD(err.Error())
		return false
	}
	return true
}

//判断用户id是否存在
func (this userDao) ExistsUser(userId string) bool {
	conn := this.pool.Get()
	defer conn.Close()

	res, _ := redis.Bool(conn.Do("exists", defined.Redis_User_Info+userId))
	return res
}

//判断用户昵称是否存在
func (this userDao) ExistUserName(userName string) bool {
	conn := this.pool.Get()
	defer conn.Close()

	res, _ := redis.Bool(conn.Do("exists", defined.Redis_UserName_Prefix+userName+defined.Redis_UserName_Postfix))
	return res
}

//userName换取userId,返回string，不存在返回""
func (this userDao) GetUserIdByUserName(userName string) (string, error) {
	conn := this.pool.Get()
	defer conn.Close()

	res, err := conn.Do("get", defined.Redis_UserName_Prefix+userName+defined.Redis_UserName_Postfix)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", err
	}
	return fmt.Sprintf("%s", res), nil
}

//根据userId换取userName，返回string
func (this userDao) GetUserNameByUserId(userId string) (string, error) {
	conn := this.pool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("hget", defined.Redis_User_Info+userId, "user_name"))
	if err != nil {
		return "", err
	}
	return res, nil
}

//根据多个userId换取多个userName，返回map [userId]userName
func (this userDao) GetUsersNameByUserId(userIds []string) (res map[string]string, err error) {
	conn := this.pool.Get()
	defer conn.Close()

	sendNum := 0
	for _, v := range userIds {
		if err = conn.Send("hget", defined.Redis_User_Info+v, "user_name"); err != nil {
			return
		}
		sendNum++
	}
	//pipeline发送，减少网络开销
	if err = conn.Flush(); err != nil {
		return
	}

	data := make([]string, 0, len(userIds))
	for i := 0; i < sendNum; i++ {
		if value, err := redis.String(conn.Receive()); err != nil {
			return res, err
		} else {
			data = append(data, value)
		}
	}

	res = make(map[string]string, len(userIds))
	for k, v := range data {
		res[userIds[k]] = v
	}
	return res, nil
}

//添加在线用户集合
func (this userDao) AddOnlineUser(userName string) error {
	conn := this.pool.Get()
	defer conn.Close()

	_, err := conn.Do("sadd", defined.Redis_Online_User, userName)
	if err != nil {
		utils.SDD(err.Error())
		return err
	}
	return nil
}

//移除在线用户集合
func (this userDao) RemoveOnlineUser(userName string) error {
	conn := this.pool.Get()
	defer conn.Close()

	_, err := redis.Bool(conn.Do("srem", defined.Redis_Online_User, userName))
	return err
}

//返回当前所有在线用户,[]string{用户昵称1, 用户昵称2}
func (this userDao) GetAllOnlineUserName() ([]string, error) {
	conn := this.pool.Get()
	defer conn.Close()

	res, err := redis.Strings(conn.Do("smembers", defined.Redis_Online_User))
	if err != nil {
		utils.SDD(err.Error())
		return []string{}, err
	}
	return res, nil
}

//查询当前用户总数
func (this userDao) GetAllUserNum() int {
	conn := this.pool.Get()
	defer conn.Close()

	curosr := "0"
	user := make([]string, 0, 100)
	for {
		//这里查询key用scan分步返回，用keys命令遍历所有key，造成阻塞
		res, _ := redis.Values(conn.Do("scan", curosr, "match", defined.Redis_User_Info+"*", "count", 100))
		curosr = fmt.Sprintf("%s", res[0])
		if userSlice, ok := res[1].([]interface{}); ok {
			for _, v := range userSlice {
				user = append(user, fmt.Sprintf("%s", v))
			}
		}
		if curosr == "0" { //游标为0表示遍历完成跳出
			break
		}
	}

	//scan命令可能会返回重复值，这里去重
	return len(utils.RemoveDuplicateElement(user))
}

//查询好友列表，返回好友的userId
func (this userDao) GetFriendList(userId string) ([]string, error) {
	conn := this.pool.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("SMEMBERS", defined.Redis_Friend+userId))
}

//查询全部好友的会话id map[userId]chatId
func (this userDao) GetFriendChatMapp(userId string) (map[string]string, error) {
	conn := this.pool.Get()
	defer conn.Close()

	//查询好友会话映射表
	friendChat, err := redis.ByteSlices(conn.Do("hgetall", defined.Redis_Friend_Chat+userId))
	if err != nil {
		return map[string]string{}, err
	}
	res := make(map[string]string, len(friendChat))
	for i := 0; i < len(friendChat); i += 2 {
		res[fmt.Sprintf("%s", friendChat[i])] = fmt.Sprintf("%s", friendChat[i+1])
	}
	return res, nil
}


