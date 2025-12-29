package main

import (
	"log"

	"github.com/flm210/mihomo-in-k8s/utils"
)

func main() {

	if err := utils.GenerateConfigFromEnv(); err != nil {
		log.Fatalf("生成配置失败: %v", err)
	}
}
