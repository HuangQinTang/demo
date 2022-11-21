package tailservice

import (
	"github.com/hpcloud/tail"
)

// InitTail 使用tail读取日志，返回tail对象
func InitTail(path string) (*tail.Tail, error) {
	cfg := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}

	tailObj, err := tail.TailFile(path, cfg)
	if err != nil {
		return nil, err
	}
	return tailObj, nil
}
