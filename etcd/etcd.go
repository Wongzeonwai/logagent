package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Wongzeonwai/logagent/logger"
	"github.com/Wongzeonwai/logagent/module"
	"github.com/Wongzeonwai/logagent/mtail"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

var client *clientv3.Client

func InitEtcd(addr module.EtcdConfig) (err error) {
	addrs := strings.Split(addr.Address, ",")
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logger.Logg.Errorf("connect etcd failed,err:%v\n", err)
		return err
	}
	return
}

func GetCollectConf(key string) (conf []module.CollectConf, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := client.Get(ctx, key)
	if err != nil {
		logger.Logg.Errorf("get collect conf from etcd failed,err:%v\n", err)
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		logger.Logg.Error("collect conf len: 0")
		return nil, errors.New("collect conf len: 0")
	}
	ret := resp.Kvs[0].Value
	err = json.Unmarshal(ret, &conf)
	if err != nil {
		logger.Logg.Errorf("unmarshal conf from etcd failed,err:%v\n", err)
		return nil, err
	}
	return conf, nil
}

// 监控etcd配置变化
func WatchConf(key string) {
	for {
		wCh := client.Watch(context.Background(), key)
		for w := range wCh {
			for _, ev := range w.Events {
				var newConf []module.CollectConf
				fmt.Printf("Type: %v, Key: %v, Value: %v\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				if ev.Type == clientv3.EventTypeDelete {
					// 如果是删除操作
					logger.Logg.Warnf("etcd delete key: %v\n", ev.Kv.Key)
					mtail.SendNewConf(newConf)
					continue
				}
				if err := json.Unmarshal(ev.Kv.Value, &newConf); err != nil {
					logger.Logg.Errorf("unmarshal conf from etcd failed,err:%v\n", err)
					mtail.SendNewConf(newConf)
					continue
				}
			}
		}
	}
}
