package mtail

import (
	"context"
	"github.com/Wongzeonwai/logagent/logger"
	"github.com/Wongzeonwai/logagent/module"
	"github.com/hpcloud/tail"
	"time"
)

// tails管理者
type TailMgr struct {
	tailMap         map[string]*tailTask
	collectConfList []module.CollectConf      // 所有配置项
	confChan        chan []module.CollectConf // 等待新配置项的管道
}

var tailMgr *TailMgr

func InitTail(collects []module.CollectConf) (err error) {
	tailMgr = &TailMgr{
		tailMap:         make(map[string]*tailTask, 20),
		collectConfList: collects,
		confChan:        make(chan []module.CollectConf),
	}
	cfg := tail.Config{
		Follow:    true,
		ReOpen:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	for _, collect := range collects {
		tt := &tailTask{
			path:  collect.Path,
			topic: collect.Topic,
		}
		tt.tObj, err = tail.TailFile(collect.Path, cfg)
		if err != nil {
			logger.Logg.Errorf("tail %s failed, err: %v\n", collect.Path, err)
			continue
		}
		tailMgr.tailMap[collect.Path] = tt // 把创建的tail登记
		go tt.start()
	}
	go tailMgr.watch()
	newConf := <-confChan
	logger.Logg.Infof("get new etcd conf: %v", newConf)
	// 管理之前启动的tail
	return
}

func (tm *TailMgr) watch() {
	for {
		newConf := <-tm.confChan
		logger.Logg.Infof("get new etcd conf: %v", newConf)
		// 管理之前启动的tail
		for _, conf := range newConf {
			// 1.原来有的任务不需要改变
			if _, ok := tm.tailMap[conf.Path]; ok {
				continue
			}
			// 2.创建新加入的任务
			cfg := tail.Config{
				Follow:    true,
				ReOpen:    true,
				Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
				MustExist: false,
				Poll:      true,
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			tt := &tailTask{
				path:   conf.Path,
				topic:  conf.Topic,
				ctx:    ctx,
				cancel: cancel,
			}
			var err error
			tt.tObj, err = tail.TailFile(conf.Path, cfg)
			if err != nil {
				logger.Logg.Errorf("tail %s failed, err: %v\n", conf.Path, err)
				continue
			}
			tailMgr.tailMap[conf.Path] = tt // 把创建的tail登记
			go tt.start()
		}
		// 3.停止被删除的任务，找到tailMap存在，newconf不存在的任务
		for k, task := range tailMgr.tailMap {
			var found bool
			for _, conf := range newConf {
				if k == conf.Path {
					found = true
					break
				}
			}
			if !found {
				delete(tailMgr.tailMap, k)
				task.cancel()
			}
		}
	}
}

func SendNewConf(newConf []module.CollectConf) {
	tailMgr.confChan <- newConf
}
