package conf

import (
	"gopkg.in/ini.v1"
	"log"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	cfg      = new(Config)
	onceSync sync.Once
)

// GetCfg 单例获取配置文件对象
func GetCfg() *Config {
	onceSync.Do(func() {
		//1.载入配置
		_, fileName, _, _ := runtime.Caller(0) //当前文件路径
		confPath := path.Dir(fileName) + "/config.ini"
		confPath = filepath.FromSlash(confPath) //系统路径符转换
		conf, err := ini.Load(confPath)
		if err != nil {
			log.Fatalf("load conf failed, err: %s", err.Error())
		}

		//2.解析配置
		err = conf.MapTo(cfg)
		if err != nil {
			log.Fatalf("load conf failed, err: %s", err.Error())
		}
	})
	return cfg
}
