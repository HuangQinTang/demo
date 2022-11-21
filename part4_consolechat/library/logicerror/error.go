package logicerror

import (
	"chat/defined"
)

// Lerror Lerror
type Lerror struct {
	code    defined.HttpCode
	message string
}

func (e Lerror) Error() string {
	return e.message
}

// Code Code
func (e Lerror) Code() defined.HttpCode {
	return e.code
}

// New New
func New(code defined.HttpCode, message string) *Lerror {

	return &Lerror{code, message}

}

//错误处理函数
//func CheckErr(err error, extra string) bool {
//	if err != nil {
//		formatStr := " Err : %s\n"
//		if extra != "" {
//			formatStr = extra + formatStr
//		}
//		logs.Log("rpc", "err", "formatStr %s  error %v", formatStr, err)
//		//g.Log().Error("formatStr %s  error %v", formatStr, err)
//		return true
//	}
//	return false
//}
