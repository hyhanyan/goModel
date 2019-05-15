package service

import (
	"cmlruleprocess/config/notify"
	"fmt"
	"time"
)

func Run() {
	for {
		cmlRuleConfig := notify.CmlRuleConfigNotifyMgr.Config.Load().(*notify.CmlRuleConfig)

		fmt.Println("Hostname:", cmlRuleConfig.Hostname)
		fmt.Println("kafkaPort:", cmlRuleConfig.KafkaPort)
		fmt.Printf("%v\n", "--------------------")
		time.Sleep(5 * time.Second)
	}
}
