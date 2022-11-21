package clientBoot

import (
	"chat/defined"
)

var (
	//菜单响应管道
	MenuJob chan defined.Message
)

//全局初始化
func init() {
	MenuJob = make(chan defined.Message)
}