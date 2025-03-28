# 日志收集项目

### 使用技术
|   | tech      | version | desc    | addr                         |
|---|-----------|---------|---------|------------------------------|
| 1 | Go        | 1.23.5  | 开发语言    | https://go.dev/              |
| 2 | Kafka     | 3.9.0   | 消息中间件   | https://kafka.apache.org     |
| 3 | Zookeeper | 3.8.4   | 分布式协调服务 | https://zookeeper.apache.org |
| 4 | Etcd      | 3.5.20  | 服务注册发现  | https://etcd.io/             |

### 使用的Go包
|    | pkg      | version | desc        | addr                               |
|----|----------|---------|-------------|------------------------------------|
| 1  | sarama   | 1.45.1  | Kafka的Go客户端 | https://github.com/IBM/sarama      |
| 2  | go-ini   | 1.67.0  | Go语言INI文件操作 | https://github.com/go-ini/ini      |
| 3  | tail     | 1.0.0   | 监控文件        | https://github.com/hpcloud/tail    |
| 4  | logrus   | 1.9.3   | 打印日志        | https://github.com/sirupsen/logrus |
| 5  | gopsutil | 3.21.11 | 采集系统信息      | https://github.com/shirou/gopsutil |

## 用ETCD存储日志配置项
key: collect
```json
[
  {
    "path": "D:\\xxx.log",
    "topic": "xxx"
  }
]
```