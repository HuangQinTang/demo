package defined

import "errors"

var (
	ERROR_CONN_LOST       = errors.New("连接断开")
	ERROR_DATA_LENGTH     = errors.New("数据长度错误")
	ERROR_USER_NOT_ONLINE = errors.New("用户不在线")
)
