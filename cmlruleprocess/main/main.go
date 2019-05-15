package main

import (
	"cmlruleprocess/service"
)

func main() {
	confFile := "../config/test.cfg"
	service.InitConfig(confFile)
	// 应用程序 很多配置已经不是存在文件中而是etcd
	service.Run()
}
