package notify

import (
	"cmlruleprocess/config/reconf"
	"fmt"
	"sync/atomic"
)

type CmlRuleConfig struct {
	Hostname  string
	Port      int
	KafkaAddr string
	KafkaPort int
}

// reload()协程写 和 for循环的读，都是对Appconfig对象，因此有读写冲突
type CmlRuleConfigMgr struct {
	Config atomic.Value
}

// 初始化结构体
var CmlRuleConfigNotifyMgr = &CmlRuleConfigMgr{}

func (cml *CmlRuleConfigMgr) Callback(conf *reconf.Config) {
	cmlRuleConfig := &CmlRuleConfig{}
	hostname, err := conf.GetString("hostname")
	if err != nil {
		fmt.Printf("get hostname err: %v\n", err)
		return
	}
	cmlRuleConfig.Hostname = hostname

	kafkaPort, err := conf.GetInt("kafkaPort")
	if err != nil {
		fmt.Printf("get kafkaPort err: %v\n", err)
		return
	}
	cmlRuleConfig.KafkaPort = kafkaPort

	CmlRuleConfigNotifyMgr.Config.Store(cmlRuleConfig)

}
