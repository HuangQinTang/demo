package log

import (
	"io/ioutil"
	stlog "log"
	"net/http"
	"os"
)

var log *stlog.Logger

type fileLog string

// Write 将data写入到fileLog路径
func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

// Run 创建log日志对象
func Run(destination string) {
	log = stlog.New(fileLog(destination), "[go] - ", stlog.LstdFlags)
}

func RegisterHandlers() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := ioutil.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

// write 这里Printf还会执行log日志对象实现的Write方法
func write(message string) {
	log.Printf("%v\n", message)
}
