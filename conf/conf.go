package conf

import (
	"github.com/Wongzeonwai/logagent/logger"
	"github.com/Wongzeonwai/logagent/module"
	"github.com/go-ini/ini"
)

func InitConfig() *module.Config {
	config := new(module.Config)
	err := ini.MapTo(config, "./conf/conf.ini")
	if err != nil {
		logger.Logg.Error("Open config file failed, err:", err)
	}
	return config
}
