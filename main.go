package main

import (
	"fmt"
	"github.com/Wongzeonwai/logagent/conf"
	"github.com/Wongzeonwai/logagent/etcd"
	"github.com/Wongzeonwai/logagent/kafka"
	"github.com/Wongzeonwai/logagent/logger"
	"github.com/Wongzeonwai/logagent/module"
	"github.com/Wongzeonwai/logagent/mtail"
	util "github.com/Wongzeonwai/logagent/utils"
)

var (
	config *module.Config
)

// 日志收集客户端
func main() {
	logger.NewLogger()
	// 获取本机IP地址
	ip, err := util.GetIP()
	if err != nil {
		logger.Logg.Error("Get IP fail, err:", err)
		return
	}
	// 1、读取配置文件
	config = conf.InitConfig()
	// 2、初始化，连接Kafka
	err = kafka.InitKafka(config.KafkaConfig)
	if err != nil {
		logger.Logg.Error("Init Kafka failed, err:", err)
		return
	}
	// 3、初始化etcd
	err = etcd.InitEtcd(config.EtcdConfig)
	if err != nil {
		logger.Logg.Error("Init etcd failed, err:", err)
		return
	}
	collectKey := fmt.Sprintf(config.CollectKey, ip)
	collectLogConf, err := etcd.GetCollectConf(collectKey)
	if err != nil {
		logger.Logg.Error("get collect conf from etcd failed, err:", err)
		return
	}
	go etcd.WatchConf(collectKey)
	// 4、根据配置文件中的路径使用tail收集日志
	err = mtail.InitTail(collectLogConf)
	if err != nil {
		logger.Logg.Error("Init Tail failed, err:", err)
		return
	}
	select {}
}
