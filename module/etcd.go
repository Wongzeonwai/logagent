package module

// ETCD收集的日志配置项
type CollectConf struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}
