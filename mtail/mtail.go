package mtail

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/Wongzeonwai/logagent/kafka"
	"github.com/Wongzeonwai/logagent/logger"
	"github.com/Wongzeonwai/logagent/module"
	"github.com/hpcloud/tail"
	"strings"
	"time"
)

type tailTask struct {
	path   string
	topic  string
	tObj   *tail.Tail
	ctx    context.Context
	cancel context.CancelFunc
}

var confChan chan []module.CollectConf

// 读取日志发往kafka
func (t *tailTask) start() {
	logger.Logg.Infof("tail %s start\n", t.path)
	for {
		select {
		case <-t.ctx.Done():
			logger.Logg.Warnf("tail %s stop\n", t.path)
			return
		case line, ok := <-t.tObj.Lines:
			if !ok {
				logger.Logg.Warnf("tail file %s is closed, try to reopen\n", t.path)
				time.Sleep(5 * time.Second)
				continue
			}
			if len(strings.Trim(line.Text, "\r")) == 0 {
				continue
			}
			msg := &sarama.ProducerMessage{}
			msg.Topic = t.topic
			msg.Value = sarama.StringEncoder(line.Text)
			kafka.SendMsgChan(msg)
		}
	}
}
