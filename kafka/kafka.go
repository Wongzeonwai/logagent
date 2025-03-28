package kafka

import (
	"github.com/IBM/sarama"
	"github.com/Wongzeonwai/logagent/logger"
	"github.com/Wongzeonwai/logagent/module"
	"strings"
)

var (
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

func InitKafka(kafkaConf module.KafkaConfig) (err error) {
	config := sarama.NewConfig()
	// 发送数据后需要主从生产者都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 新选出一个partition
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 成功交付的消息将在channel中返回
	config.Producer.Return.Successes = true

	// 连接Kafka
	addrs := strings.Split(kafkaConf.Address, ",")
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		logger.Logg.Error("Connect Kafka failed, err:", err)
		return
	}
	msgChan = make(chan *sarama.ProducerMessage, kafkaConf.ChanSize)
	go sendMsg()
	return
}

// 从通道中读取消息发送给Kafka
func sendMsg() {
	for {
		select {
		case msg := <-msgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logger.Logg.Error("send message failed, err:", err)
				return
			}
			logger.Logg.Infof("send message success, pid:%v, offset:%v\n", pid, offset)
		}
	}
}

func SendMsgChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
