package main

import (
	"log"

	_ "github.com/Sugarwithfish/Software-System-Design/config"
	configs "github.com/Sugarwithfish/Software-System-Design/config"
	"github.com/Sugarwithfish/Software-System-Design/db"
	"github.com/Sugarwithfish/Software-System-Design/routes"
)

func main() {
	// 加载配置
	config := configs.GetConfig()

	// 初始化数据库
	db.InitDB(config.DBPath)

	r := routes.SetupRouter()
	log.Printf("服务器启动于 %s 端口", config.ServerAddr)

	err := r.Run(config.ServerAddr)
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}

	r.Run(config.ServerAddr)
}
