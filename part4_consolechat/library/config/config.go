package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	cfg * tomlConfig
	once sync.Once
)

func GetConfig() *tomlConfig {
	once.Do(func() {
		//当前文件路径
		_, filename, _, _ := runtime.Caller(0)			//D:/chat/library/config/config.go
		//配置文件路径
		configPath := path.Dir(path.Dir(path.Dir(filename))) + "/config/config.toml"
		//系统路径符处理
		configPath = filepath.FromSlash(configPath)		//window 反斜杠（\）  linux 正斜杠（\）

		//判断文件路径是否存在
		_, err := os.Lstat(configPath)
		if notExist := os.IsNotExist(err); notExist {
			log.Println("路径不存在")
		}
		//解析配置文件到全局变量中
		if _ , err = toml.DecodeFile(configPath, &cfg); err != nil {
			log.Fatalln(err.Error())
		}
	})
	return cfg
}

