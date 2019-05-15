package service

import (
	"cmlruleprocess/config/notify"
	"cmlruleprocess/config/reconf"
	"fmt"
)

func InitConfig(file string) {
	// [1] 打开配置文件
	conf, err := reconf.NewConfig(file)
	if err != nil {
		fmt.Printf("read config file err: %v\n", err)
		return
	}

	// 添加观察者
	conf.AddObserver(notify.CmlRuleConfigNotifyMgr)

	// [2]第一次读取配置文件
	var cmlRuleConfig notify.CmlRuleConfig
	cmlRuleConfig.Hostname, err = conf.GetString("hostname")
	if err != nil {
		fmt.Printf("get hostname err: %v\n", err)
		return
	}
	fmt.Println("Hostname:", cmlRuleConfig.Hostname)

	cmlRuleConfig.KafkaPort, err = conf.GetInt("kafkaPort")
	if err != nil {
		fmt.Printf("get kafkaPort err: %v\n", err)
		return
	}
	fmt.Println("kafkaPort:", cmlRuleConfig.KafkaPort)

	// [3] 把读取到的配置文件数据存储到atomic.Value
	notify.CmlRuleConfigNotifyMgr.Config.Store(&cmlRuleConfig)
	fmt.Println("first load sucess.")

}
