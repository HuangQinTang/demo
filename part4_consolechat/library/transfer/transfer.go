package transfer

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"chat/defined"
	"chat/utils"
	"io"
	"net"
)

//数据传输对象结构体定义
type Transfer struct {
	Conn net.Conn   //连接对象
	Buf  [8096]byte //传输时，使用缓存
}

func NewTransfer(conn net.Conn) *Transfer {
	return &Transfer{
		Conn: conn,
		Buf:  [8096]byte{},
	}
}

//读取连接对象发送的数据,注意,连接断开没有关闭连接对象，由调
func (this *Transfer) ReadPkg() (mes defined.Message, err error) {
	utils.DD("准备读取数据...")

	//我们和客户端约定先发送数据的长度，这里拿到的是将要发送数据的长度，把该长度放入Buf
	n, err := this.Conn.Read(this.Buf[:4]) //conn.Read，如果客户端没有断开连接，协程会阻塞在这
	if err != nil {
		if err == io.EOF {
			return defined.Message{}, defined.ERROR_CONN_LOST
		}
		utils.DD("conn.Read err = " + err.Error())
		return defined.Message{}, defined.ERROR_CONN_LOST
	}
	if n != 4 {
		return mes, defined.ERROR_DATA_LENGTH
	}

	var pkgLen uint32 //需要读取多少长度的数据
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4])
	utils.DD("读取的长度为...\n" + fmt.Sprintf("%v", pkgLen))

	//如果接收的数据大于8096字节，循环读取
	if pkgLen > 8096 {
		num := pkgLen/8096 + 1
		temp := make([]byte, 0, pkgLen)
		totalNum := 0
		for i := 0; i < int(num); i++ {
			read := 8096
			if totalNum + 8096 >= int(pkgLen) {
				read = int(pkgLen)-totalNum
			}
			if _, err = this.Conn.Read(this.Buf[:read]); err != nil {
				return defined.Message{}, defined.ERROR_CONN_LOST
			}
			temp = append(temp, this.Buf[:read]...)
			totalNum += 8096
		}
		if err = json.Unmarshal(temp, &mes); err != nil {
			return
		}
		utils.DD("读取到的数据...\n" + string(temp))
		return
	}

	//接收数据小于8096字节
	n, err = this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		return
	}
	if err = json.Unmarshal(this.Buf[:pkgLen], &mes); err != nil {
		return
	}
	utils.DD("读取到的数据...\n" + string(this.Buf[:pkgLen]))
	return
}

//向指定连接对象发送数据长度
func (this *Transfer) SendPkgLen(data []byte) (pkgLen uint32, err error) {
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:4], pkgLen)
	n, err := this.Conn.Write(buf[:4])
	if n != 4 {
		utils.DD("数据没有发送完整")
	}
	if err != nil {
		fmt.Println("conn.Write(bytes) fail ", err.Error())
		return 0, err
	}
	return
}

//发送数据给指定连接对象
func (this *Transfer) WritePkg(data []byte) (err error) {
	//先发送长度
	pkgLen, err := this.SendPkgLen(data)
	if err != nil {
		return
	}

	//发送的数据大于8096字节时，按最多8096字节作一批发送
	if pkgLen > 8096 {
		num := pkgLen/8096 + 1
		for i := 0; i < int(num); i++ {
			start := i * 8096
			end := start + 8096
			if end > int(pkgLen) {
				end = int(pkgLen)
			}
			if _, err = this.Conn.Write(data[start:end]); err != nil {
				utils.DD("数据发送失败")
				return errors.New("数据没有发送完整~")
			}
		}
		return
	}

	//发送数据本身
	n, err := this.Conn.Write(data)
	if int(pkgLen) != n {
		utils.DD("数据没有发送完整")
		return errors.New("数据没有发送完整~")
	}
	if err != nil {
		utils.DD("数据发送失败")
		return
	}
	return
}

func (this *Transfer) Close() {
	this.Conn.Close()
}

//检查错误，并判断传入的错误对象是否是ERROR_CONN_LOST，是则断开连接并返回true表示断开连接打印错误，不是只打印错误，只有客户端在用
func (this *Transfer) CheckError(err error) (result bool) {
	if err != nil { //报错
		utils.DD(err.Error())
		if err == defined.ERROR_CONN_LOST {
			this.Close()
			result = true
		} else { //普通异常，返回友好提示
			fmt.Println("服务异常~")
		}
	}
	return result
}
