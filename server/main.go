package main

import (
	"github.com/Sugarwithfish/Software-System-Design/server"
	"github.com/Sugarwithfish/Software-System-Design/utils"
	"log"

	_ "github.com/Sugarwithfish/Software-System-Design/config"
	configs "github.com/Sugarwithfish/Software-System-Design/config"
	"github.com/Sugarwithfish/Software-System-Design/db"
	"github.com/Sugarwithfish/Software-System-Design/routes"

)

func main() {
	// 加载配置
	config := configs.GetConfig()

	//初始化大模型配置
	server.InitLLM(config.OpenAIAPIKey,config.OpenAIBaseURL,config.OpenAIModel)
	// 初始化数据库
	db.InitDB(config.DBPath)
	//初始化雪花算法节点
	utils.InitSnowflakeNode(1000)


	r := routes.SetupRouter()
	log.Printf("服务器启动于 %s 端口", config.ServerAddr)

	err := r.Run(config.ServerAddr)
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}

	r.Run(config.ServerAddr)
}
