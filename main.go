package main

import "log"

func main() {
	config, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("无法加载配置文件 'config.yaml': %v", err)
	}
	log.Printf("配置加载成功: %d 个线程, %d 次测试, %d 个端口, %d 个CIDR地址段; 保存结果文件: %s\n",
		config.Threads, config.PingCount, len(config.Ports), len(config.IPv4CIDRs), config.SaveFileName)
	RunTasks(config)
}
