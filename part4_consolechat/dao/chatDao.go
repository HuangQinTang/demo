package dao

/*var (
	ChatDao     *chatDao
	chatDaoOnce sync.Once
)

type chatDao struct {
	pool *redis.Pool	//连接池
	tempConn redis.Conn	//临时连接，用来启用事务
}

func NewChatDao() *chatDao {
	chatDaoOnce.Do(func() {
		ChatDao = &chatDao{
			pool: redisPool.Pool,
		}
	})
	return ChatDao
}

//注入临时连接对象
func (d chatDao) UseTempConn(conn redis.Conn) *chatDao {
	d.tempConn = conn
	return &d
}

func(d chatDao) GetConn() redis.Conn {
	if d.tempConn.Err() == nil {	//临时连接对象未关闭
		return d.tempConn
	}
	return d.pool.Get()
}

//获取全局自增消息id
func (d chatDao) GetGobalChatId() (chatId int, err error) {
	conn := d.GetConn()
	defer conn.Close()

	chatId, err = redis.Int(conn.Do("incr", defined.Redis_Chat_Id))
	if err != nil {
		return
	}
	return chatId, nil
}*/


