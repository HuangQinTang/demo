package service

import (
	"encoding/json"
	"fmt"
	"chat/defined"
	"chat/library/config"
	"chat/library/transfer"
	"chat/utils"
	"net"
	"os"
)

var ApiService = apiService{}

type apiService struct{}

//向服务端发送一条数据
func (this apiService) Send(mesType string, data map[string]interface{}, menu *Menu) (err error) {
	//1.连接服务器
	if menu.Transfer == nil {	//如果为空则建立连接，不为空复用连接
		conn, err := net.Dial("tcp", config.GetConfig().Server.Address)
		if err != nil {
			utils.CDD("net.Dial err = " + err.Error())
			fmt.Println("连接服务器失败！")
			os.Exit(500)
		}
		menu.Transfer = transfer.NewTransfer(conn)
	}

	//2.准备发送消息给服务
	var mes defined.Message
	mes.Type = mesType
	djson, err := json.Marshal(data)
	if err != nil {
		utils.CDD("json.Marshal err = " + err.Error())
		return
	}
	mes.Data = string(djson)
	mString, err := json.Marshal(mes)
	if err != nil {
		utils.CDD("json.Marshal err = " + err.Error())
		return
	}

	//3.发送数据
	if err = menu.Transfer.WritePkg(mString); err != nil {
		return
	}

	utils.CDD("发送数据成功...")
	utils.CDD(fmt.Sprintf("数据大小...\n%d\n", len(data)))
	utils.CDD(fmt.Sprintf("数据\n%s\n", data))
	return nil
}

//向服务端发送一条数据，阻塞等待服务端响应，返回响应结果
func (this apiService) SendAndRead(mesType string, data map[string]interface{}, menu *Menu) (mes defined.Message, err error) {
	//发送数据
	if err = this.Send(mesType, data, menu); err != nil {
		return
	}

	//等待服务端响应
	mes, err = menu.Transfer.ReadPkg()
	return mes, err
}